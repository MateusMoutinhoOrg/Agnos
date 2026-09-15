package serializables

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// yamlIndent is how many columns one block level is indented by, and
// yamlSimpleKeyLimit how long a key may be before it is written in the
// explicit "? key" form. Both are the defaults of gopkg.in/yaml.v3's Marshal,
// which this encoder reproduces byte for byte.
const (
	yamlIndent         = 4
	yamlSimpleKeyLimit = 128
)

// The scalar styles this encoder writes, in the order yaml.v3 falls back
// through them: a plain scalar that cannot be written plainly becomes single
// quoted, and one that cannot be single quoted becomes double quoted.
const (
	stylePlain = iota
	styleSingleQuoted
	styleDoubleQuoted
	styleLiteral
)

// yamlEncoder is a block-style YAML emitter: a port of the libyaml emitter as
// gopkg.in/yaml.v3 drives it, keeping the same indentation rule, the same
// scalar-style choice and the same escaping. Line breaking is left out on
// purpose — Marshal never sets a width, and an unset width is an unlimited
// one, so no scalar this encoder writes is ever folded.
//
// column, whitespace and indention are the emitter state the indentation
// rules read: whether anything was written on this line, whether the last
// byte was a space, and whether nothing but indentation was.
type yamlEncoder struct {
	out        []byte
	column     int
	indent     int
	indents    []int
	whitespace bool
	indention  bool
}

// encodeYaml renders one value as a YAML document, byte for byte the way
// yaml.Marshal renders it. The accepted values are the ones ParseYaml
// produces and the ones the Create* constructors build: nil, bool, the
// integer and floating-point kinds, string, []any and map[string]any.
func encodeYaml(value any) (string, error) {
	encoder := &yamlEncoder{indent: -1, whitespace: true, indention: true}

	if err := encoder.emitNode(value, false); err != nil {
		return "", err
	}

	// The document end writes one final indent, which on a line that already
	// holds something is just the closing line break.
	encoder.indent = -1
	encoder.writeIndent()

	return string(encoder.out), nil
}

// put appends one byte, counting it as one column.
func (encoder *yamlEncoder) put(char byte) {
	encoder.out = append(encoder.out, char)
	encoder.column++
}

// writeText appends text, counting columns in runes rather than bytes.
func (encoder *yamlEncoder) writeText(text string) {
	encoder.out = append(encoder.out, text...)
	encoder.column += len([]rune(text))
}

// putBreak ends the current line.
func (encoder *yamlEncoder) putBreak() {
	encoder.out = append(encoder.out, '\n')
	encoder.column = 0
	encoder.indention = true
}

// writeIndent opens the line the next token belongs on: it breaks the current
// line unless nothing but indentation is on it, then pads to the current
// indent.
func (encoder *yamlEncoder) writeIndent() {
	indent := encoder.indent
	if indent < 0 {
		indent = 0
	}

	if !encoder.indention || encoder.column > indent || (encoder.column == indent && !encoder.whitespace) {
		encoder.putBreak()
	}

	for encoder.column < indent {
		encoder.put(' ')
	}

	encoder.whitespace = true
}

// writeIndicator writes one piece of YAML punctuation — "-", ":", "?", a
// quote, a block-scalar header. need_whitespace asks for a separating space,
// is_whitespace and is_indention say what the indicator leaves behind.
func (encoder *yamlEncoder) writeIndicator(indicator string, need_whitespace bool, is_whitespace bool, is_indention bool) {
	if need_whitespace && !encoder.whitespace {
		encoder.put(' ')
	}
	encoder.writeText(indicator)
	encoder.whitespace = is_whitespace
	encoder.indention = encoder.indention && is_indention
}

// increaseIndent opens one block level. The first level inside a sequence
// item only skips the "- " indicator, which is what lines a mapping up with
// the key that shares the dash's line; every other level aligns to the next
// multiple of yamlIndent.
func (encoder *yamlEncoder) increaseIndent(seq_item bool) {
	encoder.indents = append(encoder.indents, encoder.indent)

	if encoder.indent < 0 {
		encoder.indent = 0
		return
	}
	if seq_item {
		encoder.indent += 2
		return
	}
	encoder.indent = yamlIndent * ((encoder.indent + yamlIndent) / yamlIndent)
}

// decreaseIndent closes the level increaseIndent opened.
func (encoder *yamlEncoder) decreaseIndent() {
	encoder.indent = encoder.indents[len(encoder.indents)-1]
	encoder.indents = encoder.indents[:len(encoder.indents)-1]
}

// emitNode writes one value. seq_item reports whether it is the node of a
// sequence item, which is the one case indentation is counted differently in.
func (encoder *yamlEncoder) emitNode(value any, seq_item bool) error {
	switch typed := value.(type) {
	case map[string]any:
		return encoder.emitMapping(typed, seq_item)
	case []any:
		return encoder.emitSequence(typed, seq_item)
	default:
		return encoder.emitScalar(value, false, seq_item)
	}
}

// emitMapping writes a block mapping, its keys in the order yaml.v3 sorts
// them. An empty mapping has no block form, so it is written in flow style.
func (encoder *yamlEncoder) emitMapping(mapping map[string]any, seq_item bool) error {
	if len(mapping) == 0 {
		encoder.writeIndicator("{", true, true, false)
		encoder.writeIndicator("}", false, false, false)
		return nil
	}

	keys := make([]string, 0, len(mapping))
	for key := range mapping {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i int, j int) bool { return lessYamlKey(keys[i], keys[j]) })

	encoder.increaseIndent(seq_item)
	for _, key := range keys {
		encoder.writeIndent()

		simple := !analyzeYamlScalar(key).multiline && len(key) <= yamlSimpleKeyLimit
		if !simple {
			encoder.writeIndicator("?", true, false, true)
		}
		if err := encoder.emitScalar(key, simple, false); err != nil {
			return err
		}

		if simple {
			encoder.writeIndicator(":", false, false, false)
		} else {
			encoder.writeIndent()
			encoder.writeIndicator(":", true, false, true)
		}

		if err := encoder.emitNode(mapping[key], false); err != nil {
			return err
		}
	}
	encoder.decreaseIndent()

	return nil
}

// emitSequence writes a block sequence. An empty one, like an empty mapping,
// has only a flow form.
func (encoder *yamlEncoder) emitSequence(sequence []any, seq_item bool) error {
	if len(sequence) == 0 {
		encoder.writeIndicator("[", true, true, false)
		encoder.writeIndicator("]", false, false, false)
		return nil
	}

	encoder.increaseIndent(seq_item)
	for _, item := range sequence {
		encoder.writeIndent()
		encoder.writeIndicator("-", true, false, true)
		if err := encoder.emitNode(item, true); err != nil {
			return err
		}
	}
	encoder.decreaseIndent()

	return nil
}

// emitScalar writes one scalar in the style yamlScalarStyle picked for it. A
// scalar opens an indentation level of its own, which is what a literal block
// is indented by.
func (encoder *yamlEncoder) emitScalar(value any, key bool, seq_item bool) error {
	text, style, tag, err := yamlScalarStyle(value, key)
	if err != nil {
		return err
	}

	if tag != "" {
		encoder.writeIndicator(tag, true, false, false)
	}

	encoder.increaseIndent(seq_item)
	switch style {
	case stylePlain:
		encoder.writePlain(text)
	case styleSingleQuoted:
		encoder.writeSingleQuoted(text)
	case styleDoubleQuoted:
		encoder.writeDoubleQuoted(text)
	case styleLiteral:
		encoder.writeLiteral(text)
	}
	encoder.decreaseIndent()

	return nil
}

// writePlain writes an unquoted scalar.
func (encoder *yamlEncoder) writePlain(text string) {
	if len(text) > 0 && !encoder.whitespace {
		encoder.put(' ')
	}
	encoder.writeText(text)
	if len(text) > 0 {
		encoder.whitespace = false
	}
	encoder.indention = false
}

// writeSingleQuoted writes a single-quoted scalar, doubling every quote it
// holds. The style is never picked for a value with a line break in it, so
// nothing here has to fold one.
func (encoder *yamlEncoder) writeSingleQuoted(text string) {
	encoder.writeIndicator("'", true, false, false)
	encoder.writeText(strings.ReplaceAll(text, "'", "''"))
	encoder.writeIndicator("'", false, false, false)
	encoder.whitespace = false
	encoder.indention = false
}

// writeDoubleQuoted writes a double-quoted scalar, escaping exactly what
// yaml.v3 escapes: the quote, the backslash, every line break and everything
// unprintable.
func (encoder *yamlEncoder) writeDoubleQuoted(text string) {
	encoder.writeIndicator("\"", true, false, false)
	for _, char := range text {
		encoder.writeDoubleQuotedRune(char)
	}
	encoder.writeIndicator("\"", false, false, false)
	encoder.whitespace = false
	encoder.indention = false
}

// writeDoubleQuotedRune writes one rune of a double-quoted scalar, escaped
// when it has to be.
func (encoder *yamlEncoder) writeDoubleQuotedRune(char rune) {
	if isYamlPrintable(char) && !isYamlBreak(char) && char != '"' && char != '\\' && char != 0xFEFF {
		encoder.writeText(string(char))
		return
	}

	encoder.put('\\')
	switch char {
	case 0x00:
		encoder.put('0')
	case 0x07:
		encoder.put('a')
	case 0x08:
		encoder.put('b')
	case 0x09:
		encoder.put('t')
	case 0x0A:
		encoder.put('n')
	case 0x0B:
		encoder.put('v')
	case 0x0C:
		encoder.put('f')
	case 0x0D:
		encoder.put('r')
	case 0x1B:
		encoder.put('e')
	case 0x22:
		encoder.put('"')
	case 0x5C:
		encoder.put('\\')
	case 0x85:
		encoder.put('N')
	case 0xA0:
		encoder.put('_')
	case 0x2028:
		encoder.put('L')
	case 0x2029:
		encoder.put('P')
	default:
		digits := 8
		switch {
		case char <= 0xFF:
			encoder.put('x')
			digits = 2
		case char <= 0xFFFF:
			encoder.put('u')
			digits = 4
		default:
			encoder.put('U')
		}
		for shift := (digits - 1) * 4; shift >= 0; shift -= 4 {
			digit := byte((char >> uint(shift)) & 0x0F)
			if digit < 10 {
				encoder.put(digit + '0')
			} else {
				encoder.put(digit + 'A' - 10)
			}
		}
	}
}

// writeLiteral writes a literal block scalar: the "|" header, the indentation
// hint a leading blank needs and the chomping hint the trailing breaks ask
// for, then every line at the scalar's own indent.
func (encoder *yamlEncoder) writeLiteral(text string) {
	encoder.writeIndicator("|", true, false, false)

	if strings.HasPrefix(text, " ") || strings.HasPrefix(text, "\n") {
		encoder.writeIndicator(strconv.Itoa(yamlIndent), false, false, false)
	}
	if hint := yamlChompHint(text); hint != "" {
		encoder.writeIndicator(hint, false, false, false)
	}

	encoder.whitespace = true
	breaks := true
	for _, char := range text {
		if char == '\n' {
			encoder.putBreak()
			breaks = true
			continue
		}
		if breaks {
			encoder.writeIndent()
		}
		encoder.writeText(string(char))
		encoder.indention = false
		breaks = false
	}
}

// yamlChompHint returns the chomping indicator a literal block needs: "-"
// when the text does not end in a break, "+" when it ends in more than one,
// and nothing when the single trailing break the default keeps is right.
func yamlChompHint(text string) string {
	if !strings.HasSuffix(text, "\n") {
		return "-"
	}
	if strings.HasSuffix(text, "\n\n") || text == "\n" {
		return "+"
	}
	return ""
}

// yamlScalarStyle returns the text of one scalar, the style it is written in
// and the tag it carries. It is the union of yaml.v3's stringv, which picks a
// style from what the text would resolve back to, and libyaml's scalar
// analysis, which vetoes a style the text cannot be written in.
func yamlScalarStyle(value any, key bool) (string, int, string, error) {
	text, plain_typed, err := yamlScalarText(value)
	if err != nil {
		return "", 0, "", err
	}

	if plain_typed {
		return text, stylePlain, "", nil
	}

	// Text that is not valid UTF-8 cannot be written as itself: it becomes
	// the base64 of a !!binary scalar, and the tag is what says so, which
	// leaves nothing for the style to resolve back to.
	tag := ""
	resolves := true
	if !utf8.ValidString(text) {
		text = encodeYamlBase64(text)
		tag = "!!binary"
		resolves = false
	}

	analysis := analyzeYamlScalar(text)

	style := stylePlain
	switch {
	case strings.Contains(text, "\n"):
		style = styleLiteral
	case resolves && !yamlPlainSafe(text):
		style = styleDoubleQuoted
	}

	if key && analysis.multiline {
		style = styleDoubleQuoted
	}
	if style == stylePlain {
		if !analysis.plain_allowed {
			style = styleSingleQuoted
		}
		if len(text) == 0 && key {
			style = styleSingleQuoted
		}
	}
	if style == styleSingleQuoted && !analysis.single_quoted_allowed {
		style = styleDoubleQuoted
	}
	if style == styleLiteral && (!analysis.block_allowed || key) {
		style = styleDoubleQuoted
	}

	return text, style, tag, nil
}

// yamlScalarText renders one scalar value as text, reporting whether it is a
// value that is always written plainly — everything but a string is.
func yamlScalarText(value any) (string, bool, error) {
	switch typed := value.(type) {
	case nil:
		return "null", true, nil
	case bool:
		if typed {
			return "true", true, nil
		}
		return "false", true, nil
	case int:
		return strconv.FormatInt(int64(typed), 10), true, nil
	case int8:
		return strconv.FormatInt(int64(typed), 10), true, nil
	case int16:
		return strconv.FormatInt(int64(typed), 10), true, nil
	case int32:
		return strconv.FormatInt(int64(typed), 10), true, nil
	case int64:
		return strconv.FormatInt(typed, 10), true, nil
	case uint:
		return strconv.FormatUint(uint64(typed), 10), true, nil
	case uint8:
		return strconv.FormatUint(uint64(typed), 10), true, nil
	case uint16:
		return strconv.FormatUint(uint64(typed), 10), true, nil
	case uint32:
		return strconv.FormatUint(uint64(typed), 10), true, nil
	case uint64:
		return strconv.FormatUint(typed, 10), true, nil
	case float32:
		return yamlFloatText(float64(typed), 32), true, nil
	case float64:
		return yamlFloatText(typed, 64), true, nil
	case string:
		return typed, false, nil
	default:
		return "", false, fmt.Errorf("yaml: cannot serialize a value of type %T", value)
	}
}

// yamlFloatText formats one float the way yaml.v3 does, with the infinities
// and the not-a-number spelled as YAML spells them.
func yamlFloatText(value float64, precision int) string {
	text := strconv.FormatFloat(value, 'g', -1, precision)
	switch text {
	case "+Inf":
		return ".inf"
	case "-Inf":
		return "-.inf"
	case "NaN":
		return ".nan"
	}
	return text
}

// yamlPlainSafe reports whether text can be written unquoted without reading
// back as something other than that same text: a plain "null", "3" or "-1"
// would return as a null, an int and an int, and the YAML 1.1 booleans and
// base-60 floats are quoted for the parsers that still resolve them.
func yamlPlainSafe(text string) bool {
	tag, _ := resolveYamlScalar(text)
	return tag == yamlStrTag && !isYamlBase60Float(text) && !isYamlOldBool(text)
}

// yamlScalarAnalysis is what libyaml's scalar analysis concludes about one
// piece of text: which styles can hold it as it is.
type yamlScalarAnalysis struct {
	multiline             bool
	plain_allowed         bool
	single_quoted_allowed bool
	block_allowed         bool
}

// analyzeYamlScalar decides which styles may hold text. A block scalar cannot
// keep a space at the end of a line or an unprintable character, a single
// quoted one cannot keep a tab or a break next to a space, and a plain one
// cannot start with an indicator or hold anything that would read as one.
func analyzeYamlScalar(text string) yamlScalarAnalysis {
	if len(text) == 0 {
		return yamlScalarAnalysis{plain_allowed: true, single_quoted_allowed: true}
	}

	analysis := yamlScalarAnalysis{plain_allowed: true, single_quoted_allowed: true, block_allowed: true}

	block_indicators := strings.HasPrefix(text, "---") || strings.HasPrefix(text, "...")
	line_breaks := false
	special_characters := false
	tab_characters := false
	leading_space := false
	leading_break := false
	trailing_space := false
	trailing_break := false
	break_space := false
	space_break := false
	previous_space := false
	previous_break := false
	preceded_by_whitespace := true

	runes := []rune(text)
	for index, char := range runes {
		followed_by_whitespace := index+1 >= len(runes) || runes[index+1] == ' ' || runes[index+1] == '\t'

		if index == 0 {
			switch char {
			case '#', ',', '[', ']', '{', '}', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
				block_indicators = true
			case '?', ':', '-':
				if followed_by_whitespace {
					block_indicators = true
				}
			}
		} else {
			switch char {
			case ':':
				if followed_by_whitespace {
					block_indicators = true
				}
			case '#':
				if preceded_by_whitespace {
					block_indicators = true
				}
			}
		}

		switch {
		case char == '\t':
			tab_characters = true
		case !isYamlPrintable(char):
			special_characters = true
		}

		switch {
		case char == ' ':
			if index == 0 {
				leading_space = true
			}
			if index == len(runes)-1 {
				trailing_space = true
			}
			if previous_break {
				break_space = true
			}
			previous_space, previous_break = true, false
		case isYamlBreak(char):
			line_breaks = true
			if index == 0 {
				leading_break = true
			}
			if index == len(runes)-1 {
				trailing_break = true
			}
			if previous_space {
				space_break = true
			}
			previous_space, previous_break = false, true
		default:
			previous_space, previous_break = false, false
		}

		preceded_by_whitespace = char == ' ' || char == '\t' || isYamlBreak(char)
	}

	analysis.multiline = line_breaks

	if leading_space || leading_break || trailing_space || trailing_break {
		analysis.plain_allowed = false
	}
	if trailing_space {
		analysis.block_allowed = false
	}
	if break_space {
		analysis.plain_allowed = false
		analysis.single_quoted_allowed = false
	}
	if space_break || tab_characters || special_characters {
		analysis.plain_allowed = false
		analysis.single_quoted_allowed = false
	}
	if space_break || special_characters {
		analysis.block_allowed = false
	}
	if line_breaks {
		analysis.plain_allowed = false
	}
	if block_indicators {
		analysis.plain_allowed = false
	}

	return analysis
}

// isYamlPrintable reports whether a rune may appear as itself in a YAML
// document. It is the set libyaml calls printable, which leaves out the tab,
// the C1 controls the next-line lives in, the surrogates, the byte-order mark
// and everything above the basic multilingual plane.
func isYamlPrintable(char rune) bool {
	switch {
	case char == 0x0A:
		return true
	case char >= 0x20 && char <= 0x7E:
		return true
	case char >= 0xA0 && char <= 0xD7FF:
		return true
	case char >= 0xE000 && char <= 0xFFFD && char != 0xFEFF:
		return true
	}
	return false
}

// isYamlBreak reports whether a rune is one of the line breaks YAML knows.
// A double-quoted scalar escapes every one of them rather than writing the
// break itself.
func isYamlBreak(char rune) bool {
	switch char {
	case 0x0A, 0x0D, 0x85, 0x2028, 0x2029:
		return true
	}
	return false
}

// isYamlOldBool reports whether text is one of the booleans YAML 1.1 had and
// YAML 1.2 dropped. They are quoted on the way out so an older parser reads
// them back as the strings they are.
func isYamlOldBool(text string) bool {
	switch text {
	case "y", "Y", "yes", "Yes", "YES", "on", "On", "ON",
		"n", "N", "no", "No", "NO", "off", "Off", "OFF":
		return true
	}
	return false
}

// encodeYamlBase64 renders text that is not valid UTF-8 as the base64 of a
// !!binary scalar, wrapped the way yaml.v3 wraps it.
func encodeYamlBase64(text string) string {
	const line_length = 70

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	if len(encoded) <= line_length {
		return encoded
	}

	var lines []string
	for start := 0; start < len(encoded); start += line_length {
		end := start + line_length
		if end > len(encoded) {
			end = len(encoded)
		}
		lines = append(lines, encoded[start:end])
	}
	return strings.Join(lines, "\n") + "\n"
}

// lessYamlKey orders two mapping keys the way yaml.v3 orders them: letters
// alphabetically, runs of digits by the number they spell, and a run of
// digits before a letter only when the two keys already share one.
func lessYamlKey(first string, second string) bool {
	first_runes, second_runes := []rune(first), []rune(second)

	digits := false
	for index := 0; index < len(first_runes) && index < len(second_runes); index++ {
		if first_runes[index] == second_runes[index] {
			digits = unicode.IsDigit(first_runes[index])
			continue
		}

		first_letter := unicode.IsLetter(first_runes[index])
		second_letter := unicode.IsLetter(second_runes[index])
		if first_letter && second_letter {
			return first_runes[index] < second_runes[index]
		}
		if first_letter || second_letter {
			if digits {
				return first_letter
			}
			return second_letter
		}

		var first_end, second_end int
		var first_number, second_number int64
		if first_runes[index] == '0' || second_runes[index] == '0' {
			for back := index - 1; back >= 0 && unicode.IsDigit(first_runes[back]); back-- {
				if first_runes[back] != '0' {
					first_number, second_number = 1, 1
					break
				}
			}
		}
		for first_end = index; first_end < len(first_runes) && unicode.IsDigit(first_runes[first_end]); first_end++ {
			first_number = first_number*10 + int64(first_runes[first_end]-'0')
		}
		for second_end = index; second_end < len(second_runes) && unicode.IsDigit(second_runes[second_end]); second_end++ {
			second_number = second_number*10 + int64(second_runes[second_end]-'0')
		}

		if first_number != second_number {
			return first_number < second_number
		}
		if first_end != second_end {
			return first_end < second_end
		}
		return first_runes[index] < second_runes[index]
	}

	return len(first_runes) < len(second_runes)
}
