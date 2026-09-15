package serializables

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// The tags one plain scalar can resolve to. They are the subset of the YAML
// core schema this codec reads, spelled the way yaml.v3 spells them.
const (
	yamlStrTag       = "!!str"
	yamlNullTag      = "!!null"
	yamlBoolTag      = "!!bool"
	yamlIntTag       = "!!int"
	yamlFloatTag     = "!!float"
	yamlTimestampTag = "!!timestamp"
	yamlMergeTag     = "!!merge"
)

// yamlFloatPattern is the shape a plain scalar has to have before it is read
// as a float, and yamlBase60Pattern the sexagesimal notation YAML 1.1 had:
// this codec never reads one, and the encoder quotes anything shaped like one
// so an older parser does not either.
var (
	yamlFloatPattern  = regexp.MustCompile(`^[-+]?(\.[0-9]+|[0-9]+(\.[0-9]*)?)([eE][-+]?[0-9]+)?$`)
	yamlBase60Pattern = regexp.MustCompile(`^[-+]?[0-9][0-9_]*(?::[0-5]?[0-9])+(?:\.[0-9_]*)?$`)
)

// yamlTimestampFormats are the timestamp spellings a plain scalar is checked
// against. A match is kept as the text it was written as: the contract this
// adapter fills has no date node, so a timestamp is a string here.
var yamlTimestampFormats = []string{
	"2006-1-2T15:4:5.999999999Z07:00",
	"2006-1-2t15:4:5.999999999Z07:00",
	"2006-1-2 15:4:5.999999999",
	"2006-1-2",
}

// decodeYaml parses one YAML document into the same Go values
// yaml.Unmarshal into an interface{} produces: map[string]any, []any,
// string, int, int64, float64, bool and nil.
//
// It reads the block subset this generator writes and projects hand-write:
// block mappings and sequences, flow mappings and sequences, plain, single
// quoted and double quoted scalars, literal and folded blocks, comments and
// an optional document marker. Anchors, aliases, merge keys and explicit
// tags are reported as errors rather than quietly dropped.
func decodeYaml(data string) (any, error) {
	lines := strings.Split(strings.ReplaceAll(data, "\r\n", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	parser := &yamlParser{lines: lines}
	return parser.parseDocument()
}

// yamlParser walks the document one line at a time. The lines are held as
// written, because a block scalar takes its content raw, and are stripped of
// comments only where a comment can appear.
type yamlParser struct {
	lines []string
	index int
}

// parseDocument parses the whole document: the optional markers around it,
// then the one node it holds.
func (parser *yamlParser) parseDocument() (any, error) {
	parser.skipBlanks()

	if parser.index < len(parser.lines) && strings.TrimRight(parser.lines[parser.index], " ") == "---" {
		parser.index++
		parser.skipBlanks()
	}

	if parser.index >= len(parser.lines) {
		return nil, nil
	}

	value, err := parser.parseBlock(yamlLineIndent(parser.lines[parser.index]))
	if err != nil {
		return nil, err
	}

	parser.skipBlanks()
	if parser.index < len(parser.lines) {
		line := strings.TrimRight(parser.lines[parser.index], " ")
		if line != "..." && line != "---" {
			return nil, fmt.Errorf("yaml: line %d: %q is not part of the document", parser.index+1, trimYamlBlanks(line))
		}
	}

	return value, nil
}

// parseBlock parses the block whose content starts at column indent.
func (parser *yamlParser) parseBlock(indent int) (any, error) {
	parser.skipBlanks()
	if parser.index >= len(parser.lines) {
		return nil, nil
	}

	content := yamlStripComment(parser.lines[parser.index])
	if content == "-" || strings.HasPrefix(content, "- ") {
		return parser.parseSequence(indent)
	}
	if content == "?" || strings.HasPrefix(content, "? ") {
		return parser.parseMapping(indent)
	}
	if _, _, ok := yamlSplitEntry(content); ok {
		return parser.parseMapping(indent)
	}

	return parser.parseScalarBlock(indent)
}

// parseMapping parses every key at column indent into one map.
func (parser *yamlParser) parseMapping(indent int) (any, error) {
	mapping := map[string]any{}

	for {
		parser.skipBlanks()
		if parser.index >= len(parser.lines) {
			break
		}

		line := parser.lines[parser.index]
		if yamlLineIndent(line) != indent {
			if yamlLineIndent(line) < indent {
				break
			}
			return nil, fmt.Errorf("yaml: line %d: unexpected indentation", parser.index+1)
		}

		content := yamlStripComment(line)
		if content == "---" || content == "..." {
			break
		}

		if content == "?" || strings.HasPrefix(content, "? ") {
			key, value, err := parser.parseExplicitEntry(indent, content)
			if err != nil {
				return nil, err
			}
			mapping[fmt.Sprint(key)] = value
			continue
		}

		key_text, rest, ok := yamlSplitEntry(content)
		if !ok {
			break
		}

		key, err := parser.scalarValue(key_text)
		if err != nil {
			return nil, err
		}
		if key == yamlMergeTag {
			return nil, fmt.Errorf("yaml: line %d: merge keys are not supported", parser.index+1)
		}
		parser.index++

		value, err := parser.parseEntryValue(rest, indent)
		if err != nil {
			return nil, err
		}
		mapping[fmt.Sprint(key)] = value
	}

	return mapping, nil
}

// parseSequence parses every item whose dash is at column indent.
func (parser *yamlParser) parseSequence(indent int) (any, error) {
	sequence := []any{}

	for {
		parser.skipBlanks()
		if parser.index >= len(parser.lines) {
			break
		}

		line := parser.lines[parser.index]
		if yamlLineIndent(line) != indent {
			if yamlLineIndent(line) < indent {
				break
			}
			return nil, fmt.Errorf("yaml: line %d: unexpected indentation", parser.index+1)
		}

		content := yamlStripComment(line)
		if content != "-" && !strings.HasPrefix(content, "- ") {
			break
		}

		value, err := parser.parseIndicatorNode(indent, content)
		if err != nil {
			return nil, err
		}
		sequence = append(sequence, value)
	}

	return sequence, nil
}

// parseIndicatorNode parses the node one indicator opens: the "-" of a
// sequence item, or the "?" and ":" of an explicit mapping entry. The node
// may sit on the indicator's own line or in the block below it.
func (parser *yamlParser) parseIndicatorNode(indent int, content string) (any, error) {
	rest := strings.TrimLeft(content[1:], " ")
	if rest == "" {
		parser.index++
		return parser.parseChildBlock(indent)
	}

	// The node's own content starts after the indicator and the blanks
	// behind it. Blanking the indicator out leaves a line that reads like
	// any other at that column, so a mapping or a sequence opened here is
	// parsed by the same two functions.
	node_indent := indent + len(content) - len(rest)
	if yamlOpensBlock(rest) {
		parser.lines[parser.index] = strings.Repeat(" ", node_indent) + rest
		return parser.parseBlock(node_indent)
	}

	parser.index++
	return parser.parseEntryValue(rest, node_indent-2)
}

// yamlOpensBlock reports whether the text after an indicator opens a
// collection of its own rather than being a scalar: another indicator, or a
// mapping entry.
func yamlOpensBlock(rest string) bool {
	if rest == "-" || strings.HasPrefix(rest, "- ") || rest == "?" || strings.HasPrefix(rest, "? ") {
		return true
	}
	_, _, ok := yamlSplitEntry(rest)
	return ok
}

// parseExplicitEntry parses one entry written in the explicit form, the
// "? key" and ": value" pair a mapping falls back to for a key that is too
// long or holds a line break.
func (parser *yamlParser) parseExplicitEntry(indent int, content string) (any, any, error) {
	key, err := parser.parseIndicatorNode(indent, content)
	if err != nil {
		return nil, nil, err
	}

	parser.skipBlanks()
	if parser.index >= len(parser.lines) {
		return key, nil, nil
	}

	line := parser.lines[parser.index]
	value_content := yamlStripComment(line)
	if yamlLineIndent(line) != indent || (value_content != ":" && !strings.HasPrefix(value_content, ": ")) {
		return key, nil, nil
	}

	value, err := parser.parseIndicatorNode(indent, value_content)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

// parseEntryValue parses what follows a mapping key or a sequence dash: a
// block scalar header, an inline value, or nothing, which means the value is
// the block below or a null.
func (parser *yamlParser) parseEntryValue(rest string, indent int) (any, error) {
	if rest == "" {
		return parser.parseChildBlock(indent)
	}
	if strings.HasPrefix(rest, "|") || strings.HasPrefix(rest, ">") {
		return parser.parseBlockScalar(rest, indent)
	}
	return parser.parseInline(rest)
}

// parseChildBlock parses the block a key or a dash left open. It is indented
// deeper than the line that opened it, except for a sequence, which YAML lets
// share its parent's column.
func (parser *yamlParser) parseChildBlock(indent int) (any, error) {
	parser.skipBlanks()
	if parser.index >= len(parser.lines) {
		return nil, nil
	}

	child_indent := yamlLineIndent(parser.lines[parser.index])
	if child_indent > indent {
		return parser.parseBlock(child_indent)
	}

	content := yamlStripComment(parser.lines[parser.index])
	if child_indent == indent && (content == "-" || strings.HasPrefix(content, "- ")) {
		return parser.parseSequence(indent)
	}

	return nil, nil
}

// parseScalarBlock parses a scalar standing where a block was expected,
// including a plain one carried over more than one line.
func (parser *yamlParser) parseScalarBlock(indent int) (any, error) {
	content := yamlStripComment(parser.lines[parser.index])
	parser.index++

	if strings.HasPrefix(content, "|") || strings.HasPrefix(content, ">") {
		return parser.parseBlockScalar(content, indent-1)
	}

	folded := parser.foldContinuation(content, indent)
	return parser.parseInline(folded)
}

// parseInline parses one value written on the line of its key: a flow
// collection or a scalar, quoted or not, folded over as many lines as its
// quotes take.
func (parser *yamlParser) parseInline(text string) (any, error) {
	if strings.HasPrefix(text, "[") || strings.HasPrefix(text, "{") {
		flow, err := parser.completeFlow(text)
		if err != nil {
			return nil, err
		}
		value, rest, err := parseYamlFlow(flow)
		if err != nil {
			return nil, err
		}
		if trimYamlBlanks(rest) != "" {
			return nil, fmt.Errorf("yaml: %q has trailing content after the flow collection", text)
		}
		return value, nil
	}

	if strings.HasPrefix(text, "\"") || strings.HasPrefix(text, "'") {
		quoted, err := parser.completeQuoted(text)
		if err != nil {
			return nil, err
		}
		value, rest, err := parseYamlQuoted(quoted)
		if err != nil {
			return nil, err
		}
		if trimYamlBlanks(rest) != "" {
			return nil, fmt.Errorf("yaml: %q has trailing content after the quoted scalar", text)
		}
		return value, nil
	}

	return parser.scalarValue(text)
}

// scalarValue resolves one plain scalar, or unquotes a quoted one that is
// already complete.
func (parser *yamlParser) scalarValue(text string) (any, error) {
	text = trimYamlBlanks(text)

	if strings.HasPrefix(text, "\"") || strings.HasPrefix(text, "'") {
		value, rest, err := parseYamlQuoted(text)
		if err != nil {
			return nil, err
		}
		if trimYamlBlanks(rest) != "" {
			return nil, fmt.Errorf("yaml: %q has trailing content after the quoted scalar", text)
		}
		return value, nil
	}

	if strings.HasPrefix(text, "&") || strings.HasPrefix(text, "*") {
		return nil, fmt.Errorf("yaml: %q: anchors and aliases are not supported", text)
	}
	if strings.HasPrefix(text, "!") {
		return nil, fmt.Errorf("yaml: %q: explicit tags are not supported", text)
	}

	_, value := resolveYamlScalar(text)
	return value, nil
}

// parseBlockScalar reads a literal or folded block: its header decides how
// the lines are joined and how the trailing breaks are kept, and every line
// deeper than indent belongs to it.
func (parser *yamlParser) parseBlockScalar(header string, indent int) (any, error) {
	folded := strings.HasPrefix(header, ">")

	header = yamlStripComment(header)
	chomp := ""
	explicit_indent := 0
	for _, char := range header[1:] {
		switch {
		case char == '-' || char == '+':
			chomp = string(char)
		case char >= '1' && char <= '9':
			explicit_indent = indent + int(char-'0')
		case char == ' ':
		default:
			return nil, fmt.Errorf("yaml: %q is not a block scalar header", header)
		}
	}

	var raw []string
	content_indent := explicit_indent
	for parser.index < len(parser.lines) {
		line := parser.lines[parser.index]
		if strings.Trim(line, " ") == "" {
			raw = append(raw, "")
			parser.index++
			continue
		}
		if yamlLineIndent(line) <= indent {
			break
		}
		if content_indent == 0 {
			content_indent = yamlLineIndent(line)
		}
		if yamlLineIndent(line) < content_indent {
			break
		}
		raw = append(raw, line[content_indent:])
		parser.index++
	}

	// Blank lines after the last content line belong to the document, not to
	// the scalar, unless the block keeps its trailing breaks.
	last := len(raw)
	for last > 0 && raw[last-1] == "" {
		last--
	}
	trailing := len(raw) - last
	raw = raw[:last]

	text := ""
	if folded {
		text = yamlFoldLines(raw)
	} else if len(raw) > 0 {
		text = strings.Join(raw, "\n")
	}

	switch chomp {
	case "-":
		// Strip: the trailing breaks are dropped, the text ends at its
		// last character.
	case "+":
		// Keep: every trailing break is content, and a block that is
		// nothing but blank lines is those breaks.
		if len(raw) > 0 {
			text += strings.Repeat("\n", trailing+1)
		} else {
			text += strings.Repeat("\n", trailing)
		}
	default:
		// Clip: one trailing break is kept.
		if len(raw) > 0 {
			text += "\n"
		}
	}

	return text, nil
}

// foldContinuation joins the lines a plain scalar runs over, which YAML folds
// into single spaces.
func (parser *yamlParser) foldContinuation(first string, indent int) string {
	folded := first
	for parser.index < len(parser.lines) {
		line := parser.lines[parser.index]
		if trimYamlBlanks(line) == "" || yamlLineIndent(line) < indent {
			break
		}
		content := yamlStripComment(line)
		if content == "" {
			break
		}
		if _, _, ok := yamlSplitEntry(content); ok {
			break
		}
		if content == "-" || strings.HasPrefix(content, "- ") {
			break
		}
		folded += " " + content
		parser.index++
	}
	return folded
}

// completeQuoted returns the whole of a quoted scalar, reading on while its
// closing quote is still missing.
func (parser *yamlParser) completeQuoted(first string) (string, error) {
	text := first
	for {
		if _, _, err := parseYamlQuoted(text); err == nil {
			return text, nil
		}
		if parser.index >= len(parser.lines) {
			return "", fmt.Errorf("yaml: %q is missing its closing quote", first)
		}
		text += "\n" + strings.TrimLeft(parser.lines[parser.index], " ")
		parser.index++
	}
}

// completeFlow returns the whole of a flow collection, reading on while its
// brackets are still open.
func (parser *yamlParser) completeFlow(first string) (string, error) {
	text := first
	for {
		if _, _, err := parseYamlFlow(text); err == nil {
			return text, nil
		}
		if parser.index >= len(parser.lines) {
			return "", fmt.Errorf("yaml: %q is missing its closing bracket", first)
		}
		text += " " + strings.TrimLeft(parser.lines[parser.index], " ")
		parser.index++
	}
}

// skipBlanks advances over empty lines and lines holding nothing but a
// comment.
func (parser *yamlParser) skipBlanks() {
	for parser.index < len(parser.lines) {
		line := parser.lines[parser.index]
		if trimYamlBlanks(line) != "" && yamlStripComment(line) != "" {
			return
		}
		parser.index++
	}
}

// trimYamlBlanks removes the blanks YAML knows from both ends of text. It is
// not strings.TrimSpace: that one also removes the non-breaking space, which
// a YAML document carries as content.
func trimYamlBlanks(text string) string {
	return strings.Trim(text, " \t")
}

// yamlLineIndent is the column one line's content starts at.
func yamlLineIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// yamlStripComment returns one line's content with its indentation and its
// trailing comment removed. A "#" only opens a comment at the start of the
// content or after a blank, and never inside quotes.
func yamlStripComment(line string) string {
	content := strings.TrimLeft(line, " ")

	quote := byte(0)
	for index := 0; index < len(content); index++ {
		char := content[index]
		switch {
		case quote == '\'':
			if char == '\'' {
				if index+1 < len(content) && content[index+1] == '\'' {
					index++
				} else {
					quote = 0
				}
			}
		case quote == '"':
			if char == '\\' {
				index++
			} else if char == '"' {
				quote = 0
			}
		case (char == '\'' || char == '"') && yamlOpensNode(content, index):
			quote = char
		case char == '#' && (index == 0 || content[index-1] == ' '):
			return strings.TrimRight(content[:index], " ")
		}
	}

	return strings.TrimRight(content, " ")
}

// yamlOpensNode reports whether the character at index stands where a node
// starts, which is the only place a quote opens a quoted scalar: anywhere
// else a quote is an ordinary character of a plain one. A block indicator is
// not in the set because one that opens a node is followed by a blank, and
// that blank is what is seen here.
func yamlOpensNode(content string, index int) bool {
	if index == 0 {
		return true
	}
	switch content[index-1] {
	case ' ', '\t', ',', '[', '{':
		return true
	}
	return false
}

// yamlSplitEntry splits one mapping entry into its key and what follows the
// colon, reporting whether the line is an entry at all: the colon has to be
// outside quotes and brackets, and followed by a blank or by nothing.
func yamlSplitEntry(content string) (string, string, bool) {
	quote := byte(0)
	depth := 0

	for index := 0; index < len(content); index++ {
		char := content[index]
		switch {
		case quote != 0:
			switch {
			case quote == '"' && char == '\\':
				index++
			case quote == '\'' && char == '\'' && index+1 < len(content) && content[index+1] == '\'':
				index++
			case char == quote:
				quote = 0
			}
		case (char == '\'' || char == '"') && (index == 0 || depth > 0):
			quote = char
		case (char == '[' || char == '{') && (index == 0 || depth > 0):
			depth++
		case (char == ']' || char == '}') && depth > 0:
			depth--
		case char == ':' && depth == 0:
			if index+1 == len(content) {
				return trimYamlBlanks(content[:index]), "", true
			}
			if content[index+1] == ' ' {
				return trimYamlBlanks(content[:index]), trimYamlBlanks(content[index+1:]), true
			}
		}
	}

	return "", "", false
}

// yamlFoldLines joins the lines of a folded block: one break between two
// non-empty lines becomes a space, an empty line becomes a break of its own,
// and a line indented deeper than the block keeps the break before it.
func yamlFoldLines(lines []string) string {
	text := ""
	for index, line := range lines {
		if index == 0 {
			text = line
			continue
		}
		switch {
		case line == "":
			text += "\n"
		case strings.HasPrefix(line, " ") || strings.HasPrefix(lines[index-1], " "):
			text += "\n" + line
		case strings.HasSuffix(text, "\n"):
			text += line
		default:
			text += " " + line
		}
	}
	return text
}

// parseYamlQuoted reads one quoted scalar off the front of text and returns
// it with whatever followed it. The line breaks the scalar was written over
// are folded before anything is unescaped, so a "\n" that was written as an
// escape stays the break it asked for.
func parseYamlQuoted(text string) (any, string, error) {
	if strings.HasPrefix(text, "'") {
		for index := 1; index < len(text); index++ {
			if text[index] != '\'' {
				continue
			}
			if index+1 < len(text) && text[index+1] == '\'' {
				index++
				continue
			}
			return strings.ReplaceAll(yamlFoldQuoted(text[1:index]), "''", "'"), text[index+1:], nil
		}
		return nil, "", fmt.Errorf("yaml: %q is missing its closing quote", text)
	}

	for index := 1; index < len(text); index++ {
		switch text[index] {
		case '\\':
			index++
		case '"':
			value, err := unescapeYamlText(yamlFoldQuoted(text[1:index]))
			if err != nil {
				return nil, "", err
			}
			return value, text[index+1:], nil
		}
	}

	return nil, "", fmt.Errorf("yaml: %q is missing its closing quote", text)
}

// unescapeYamlText resolves every escape sequence of one double-quoted
// scalar's text.
func unescapeYamlText(text string) (string, error) {
	var value strings.Builder

	for index := 0; index < len(text); index++ {
		if text[index] != '\\' {
			value.WriteByte(text[index])
			continue
		}
		if index+1 >= len(text) {
			return "", fmt.Errorf("yaml: %q ends in an escape", text)
		}
		consumed, unescaped, err := unescapeYaml(text[index+1:])
		if err != nil {
			return "", err
		}
		value.WriteString(unescaped)
		index += consumed
	}

	return value.String(), nil
}

// yamlFoldQuoted folds the line breaks a quoted scalar was written over: one
// break is a space, and each further one is kept as a break.
func yamlFoldQuoted(text string) string {
	if !strings.Contains(text, "\n") {
		return text
	}

	parts := strings.Split(text, "\n")
	folded := strings.TrimRight(parts[0], " \t")
	for index, part := range parts[1:] {
		part = strings.TrimLeft(part, " \t")
		if index < len(parts)-2 {
			part = strings.TrimRight(part, " \t")
		}
		if part == "" {
			folded += "\n"
			continue
		}
		if strings.HasSuffix(folded, "\n") {
			folded += part
			continue
		}
		folded += " " + part
	}

	return folded
}

// unescapeYaml reads one escape sequence, returning how many bytes after the
// backslash it took and what it stands for.
func unescapeYaml(text string) (int, string, error) {
	switch text[0] {
	case '0':
		return 1, "\x00", nil
	case 'a':
		return 1, "\a", nil
	case 'b':
		return 1, "\b", nil
	case 't':
		return 1, "\t", nil
	case 'n':
		return 1, "\n", nil
	case 'v':
		return 1, "\v", nil
	case 'f':
		return 1, "\f", nil
	case 'r':
		return 1, "\r", nil
	case 'e':
		return 1, "\x1b", nil
	case ' ':
		return 1, " ", nil
	case '"':
		return 1, "\"", nil
	case '/':
		return 1, "/", nil
	case '\\':
		return 1, "\\", nil
	case 'N':
		return 1, "", nil
	case '_':
		return 1, " ", nil
	case 'L':
		return 1, " ", nil
	case 'P':
		return 1, " ", nil
	case '\n':
		return 1, "", nil
	case 'x', 'u', 'U':
		digits := map[byte]int{'x': 2, 'u': 4, 'U': 8}[text[0]]
		if len(text) < digits+1 {
			return 0, "", fmt.Errorf("yaml: %q is not a complete escape", text)
		}
		code, err := strconv.ParseUint(text[1:digits+1], 16, 32)
		if err != nil {
			return 0, "", fmt.Errorf("yaml: %q is not a valid escape: %w", text[:digits+1], err)
		}
		return digits + 1, string(rune(code)), nil
	}

	return 0, "", fmt.Errorf("yaml: %q is not a valid escape", text[:1])
}

// parseYamlFlow reads one flow collection off the front of text and returns
// it with whatever followed it.
func parseYamlFlow(text string) (any, string, error) {
	text = strings.TrimLeft(text, " ")

	if strings.HasPrefix(text, "[") {
		items := []any{}
		rest := strings.TrimLeft(text[1:], " ")
		if strings.HasPrefix(rest, "]") {
			return items, rest[1:], nil
		}
		for {
			item, remainder, err := parseYamlFlowNode(rest)
			if err != nil {
				return nil, "", err
			}
			items = append(items, item)
			remainder = strings.TrimLeft(remainder, " ")
			if strings.HasPrefix(remainder, ",") {
				rest = strings.TrimLeft(remainder[1:], " ")
				continue
			}
			if strings.HasPrefix(remainder, "]") {
				return items, remainder[1:], nil
			}
			return nil, "", fmt.Errorf("yaml: %q is not a closed flow sequence", text)
		}
	}

	if !strings.HasPrefix(text, "{") {
		return nil, "", fmt.Errorf("yaml: %q is not a flow collection", text)
	}

	mapping := map[string]any{}
	rest := strings.TrimLeft(text[1:], " ")
	if strings.HasPrefix(rest, "}") {
		return mapping, rest[1:], nil
	}
	for {
		key, remainder, err := parseYamlFlowNode(rest)
		if err != nil {
			return nil, "", err
		}
		remainder = strings.TrimLeft(remainder, " ")
		if !strings.HasPrefix(remainder, ":") {
			return nil, "", fmt.Errorf("yaml: %q is missing a value in a flow mapping", text)
		}
		value, remainder, err := parseYamlFlowNode(strings.TrimLeft(remainder[1:], " "))
		if err != nil {
			return nil, "", err
		}
		mapping[fmt.Sprint(key)] = value

		remainder = strings.TrimLeft(remainder, " ")
		if strings.HasPrefix(remainder, ",") {
			rest = strings.TrimLeft(remainder[1:], " ")
			continue
		}
		if strings.HasPrefix(remainder, "}") {
			return mapping, remainder[1:], nil
		}
		return nil, "", fmt.Errorf("yaml: %q is not a closed flow mapping", text)
	}
}

// parseYamlFlowNode reads one node of a flow collection: a nested collection,
// a quoted scalar, or a plain one ending at the next separator.
func parseYamlFlowNode(text string) (any, string, error) {
	text = strings.TrimLeft(text, " ")

	if strings.HasPrefix(text, "[") || strings.HasPrefix(text, "{") {
		return parseYamlFlow(text)
	}
	if strings.HasPrefix(text, "\"") || strings.HasPrefix(text, "'") {
		return parseYamlQuoted(text)
	}

	end := len(text)
	for index := 0; index < len(text); index++ {
		if strings.ContainsRune(",]}:", rune(text[index])) {
			end = index
			break
		}
	}

	_, value := resolveYamlScalar(trimYamlBlanks(text[:end]))
	return value, text[end:], nil
}

// resolveYamlScalar reads one plain scalar the way yaml.v3 resolves it: the
// tag it carries and the Go value it stands for.
func resolveYamlScalar(text string) (string, any) {
	if constant, ok := yamlConstant(text); ok {
		return constant.tag, constant.value
	}

	if text == "" {
		return yamlNullTag, nil
	}

	switch yamlResolveHint(text[0]) {
	case 'M':
		// Every value that hint stands for is a constant, checked above.
	case '.':
		if number, err := strconv.ParseFloat(text, 64); err == nil {
			return yamlFloatTag, number
		}
	case 'D', 'S':
		if _, ok := parseYamlTimestamp(text); ok {
			return yamlTimestampTag, text
		}

		plain := strings.ReplaceAll(text, "_", "")
		if number, err := strconv.ParseInt(plain, 0, 64); err == nil {
			if number == int64(int(number)) {
				return yamlIntTag, int(number)
			}
			return yamlIntTag, number
		}
		if number, err := strconv.ParseUint(plain, 0, 64); err == nil {
			return yamlIntTag, number
		}
		if yamlFloatPattern.MatchString(plain) {
			if number, err := strconv.ParseFloat(plain, 64); err == nil {
				return yamlFloatTag, number
			}
		}
		if strings.HasPrefix(plain, "-0b") {
			if number, err := strconv.ParseInt("-"+plain[3:], 2, 64); err == nil {
				return yamlIntTag, int(number)
			}
		}
		if strings.HasPrefix(plain, "-0o") {
			if number, err := strconv.ParseInt("-"+plain[3:], 8, 64); err == nil {
				return yamlIntTag, int(number)
			}
		}
	}

	return yamlStrTag, text
}

// yamlConstantValue is one of the spellings that stand for a fixed value.
type yamlConstantValue struct {
	tag   string
	value any
}

// yamlConstant reports the value one of the fixed spellings stands for.
func yamlConstant(text string) (yamlConstantValue, bool) {
	switch text {
	case "true", "True", "TRUE":
		return yamlConstantValue{yamlBoolTag, true}, true
	case "false", "False", "FALSE":
		return yamlConstantValue{yamlBoolTag, false}, true
	case "", "~", "null", "Null", "NULL":
		return yamlConstantValue{yamlNullTag, nil}, true
	case ".nan", ".NaN", ".NAN":
		return yamlConstantValue{yamlFloatTag, yamlNaN()}, true
	case ".inf", ".Inf", ".INF", "+.inf", "+.Inf", "+.INF":
		return yamlConstantValue{yamlFloatTag, yamlInf(1)}, true
	case "-.inf", "-.Inf", "-.INF":
		return yamlConstantValue{yamlFloatTag, yamlInf(-1)}, true
	case "<<":
		return yamlConstantValue{yamlMergeTag, yamlMergeTag}, true
	}
	return yamlConstantValue{}, false
}

// yamlResolveHint classifies a plain scalar by its first byte, exactly as
// yaml.v3's resolve table does: a sign or a digit may be a number, the
// letters of the constants may be one of those, and a dot may be a float.
func yamlResolveHint(first byte) byte {
	switch {
	case first == '+' || first == '-':
		return 'S'
	case first >= '0' && first <= '9':
		return 'D'
	case strings.IndexByte("yYnNtTfFoO~", first) >= 0:
		return 'M'
	case first == '.':
		return '.'
	}
	return 0
}

// parseYamlTimestamp reports whether text is written as one of the timestamp
// spellings YAML resolves.
func parseYamlTimestamp(text string) (time.Time, bool) {
	digits := 0
	for digits < len(text) && text[digits] >= '0' && text[digits] <= '9' {
		digits++
	}
	if digits != 4 || digits == len(text) || text[digits] != '-' {
		return time.Time{}, false
	}

	for _, format := range yamlTimestampFormats {
		if parsed, err := time.Parse(format, text); err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

// isYamlBase60Float reports whether text is written in the sexagesimal
// notation YAML 1.1 resolved as a float.
func isYamlBase60Float(text string) bool {
	if text == "" {
		return false
	}
	first := text[0]
	if !(first == '+' || first == '-' || first >= '0' && first <= '9') || !strings.Contains(text, ":") {
		return false
	}
	return yamlBase60Pattern.MatchString(text)
}

// yamlNaN and yamlInf build the floats YAML spells .nan and .inf without
// naming the math package for two constants.
func yamlNaN() float64 {
	var zero float64
	return zero / zero
}

func yamlInf(sign int) float64 {
	var zero float64
	if sign < 0 {
		return -1 / zero
	}
	return 1 / zero
}
