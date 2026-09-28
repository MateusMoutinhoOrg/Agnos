package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RouteDocField is one value a caller sends, as docs/Routes prints it: the
// name it goes under, where it goes, what it holds in plain words, whether it
// may be left out, the value the page's requests send and the declared
// description.
type RouteDocField struct {
	Name        string
	Where       string
	Type        string
	Required    string
	Example     string
	Description string
	// From is, for a parameter a route in front of this one reads, that
	// route — "" for one of the route's own — and FromPage its page.
	From     string
	FromPage string
}

// RouteDocBodyField is one property of a json body's schema: its dotted path,
// what it holds, whether it may be left out and its bounds, in plain words.
type RouteDocBodyField struct {
	Name     string
	Type     string
	Required string
	Rules    string
}

// RouteDocBody is what a route's page says about the request body: one
// sentence on what to send, the schema's properties, and a sample of it.
type RouteDocBody struct {
	Summary    string
	Fields     []RouteDocBodyField
	Sample     string
	SampleLang string
}

// RouteDocRequest is one ready-to-run call of a route: the smallest request
// it answers, or one sending every value it reads.
type RouteDocRequest struct {
	Title   string
	Command string
}

// RouteDocStatus is one status a route answers with and what it means to the
// caller.
type RouteDocStatus struct {
	Code    string
	Meaning string
}

// RouteDocReach is one route that runs in front of another one: its page, and
// the condition it runs on, "" when it always does.
type RouteDocReach struct {
	Name      string
	Page      string
	Condition string
}

// RouteDoc is one route's page of docs/Routes, rendered from its route.yaml
// and the routes crossed with it: nothing here is written by hand on the
// page. The page is read by whoever calls the server, so every label is in
// plain words and every request is one that runs.
type RouteDoc struct {
	Name            string
	Method          string
	Pattern         string
	Help            string
	LongDescription string
	// Address is one row per part of the path the caller fills in or that
	// carries a rule; a literal part is spelled by Pattern alone.
	Address    []RouteDocField
	Parameters []RouteDocField
	Body       *RouteDocBody
	Requests   []RouteDocRequest
	// Examples are the declared ones a generated request does not repeat.
	Examples []string
	Statuses []RouteDocStatus
	// Middlewares are the routes on a lower rung that run in front of it.
	Middlewares []RouteDocReach
}

// RouteDocGroup is one category section of docs/Routes, holding the routes
// that declare that category.
type RouteDocGroup struct {
	Category string
	Routes   []RouteDoc
}

// routeDocOther is the category a route with no declared one falls into.
const routeDocOther = "Other"

// routeDocHost is where the generated requests are sent: the first port
// `start-server` tries by default.
const routeDocHost = "localhost:3000"

// routeDocMayRun is how a page words a route in front whose trigger cannot be
// crossed with this one without a request.
const routeDocMayRun = "depends on the address — `explain-route` gives the exact answer"

// The samples a generated request sends for a value of a type that has one.
const (
	routeDocSampleUuid     = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	routeDocSampleDateTime = "2026-01-02T03:04:05Z"
	routeDocSampleEmail    = "someone@example.com"
	routeDocSampleUri      = "https://example.com"
	routeDocSampleText     = "text"
	routeDocSampleFile     = "my-file"
)

// routeDocEntry is one declared route, read once for every crossing.
type routeDocEntry struct {
	name string
	conf *routeconf.RouteConf
}

// routeDocInherited is one parameter a route in front of another reads, and
// whether that route runs on every request of this one — only then is a
// value it requires one this route's requests have to send.
type routeDocInherited struct {
	parameter routeconf.Parameter
	always    bool
}

// CollectRouteDocs renders every sandbox/internal/routeslist/<name>/route.yaml into
// the pages docs/Routes prints, grouped by category in first-seen order — the
// server layer's CollectCommandDocs. Hidden routes are skipped. Every route is
// crossed with every route on a lower rung whose triggers may hold on it, and
// its page lists them; the parameters of the ones sure to run in front of it
// join its own, and its requests send the ones they require.
//
// The declaration is the only source: a route, a field or an example reaches
// the page by being declared with `add-route`, `add-path`, `add-parameter`,
// `set-body` or `set-route`, never by the page being edited.
func CollectRouteDocs(sandbox *api.Sandbox, io *smartio.SmartIO) ([]RouteDocGroup, error) {
	var groups []RouteDocGroup
	index := map[string]int{}

	entries := []routeDocEntry{}
	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: %w", name, err)
		}
		entries = append(entries, routeDocEntry{name: name, conf: conf})
	}

	for _, current := range entries {
		if current.conf.Hidden {
			continue
		}

		category := current.conf.Category
		if category == "" {
			category = routeDocOther
		}

		position, seen := index[category]
		if !seen {
			position = len(groups)
			index[category] = position
			groups = append(groups, RouteDocGroup{Category: category})
		}

		groups[position].Routes = append(groups[position].Routes, routeDoc(sandbox, current, entries))
	}

	return groups, nil
}

// routeDoc turns one parsed declaration into its page, crossed with every
// other route of the project.
func routeDoc(sandbox *api.Sandbox, current routeDocEntry, entries []routeDocEntry) RouteDoc {
	conf := current.conf
	doc := RouteDoc{
		Name:            current.name,
		Method:          sandbox.Deps.Stringsdeps.Join(conf.Methods, ", "),
		Pattern:         conf.Pattern(),
		Help:            docCell(sandbox, conf.Help),
		LongDescription: docText(sandbox, conf.LongDescription),
		Body:            routeDocBody(sandbox, conf.Body),
	}

	for _, path := range conf.Paths {
		if row, shown := routeDocPath(sandbox, path); shown {
			doc.Address = append(doc.Address, row)
		}
	}
	for _, parameter := range conf.Parameters {
		doc.Parameters = append(doc.Parameters, routeDocParameter(sandbox, parameter, true))
	}

	inherited := []routeDocInherited{}
	for _, other := range entries {
		if other.name == current.name || other.conf.Hidden {
			continue
		}
		reach, condition := utils.RouteMiddlewareReach(sandbox, other.conf, conf)
		if reach == utils.NoReach {
			continue
		}
		name := utils.RouteIdentifier(sandbox, other.name)
		page := other.name + docPageExt
		doc.Middlewares = append(doc.Middlewares, RouteDocReach{
			Name:      name,
			Page:      page,
			Condition: routeReachCondition(reach, condition),
		})
		if reach != utils.Runs {
			continue
		}
		for _, parameter := range other.conf.Parameters {
			field := routeDocParameter(sandbox, parameter, false)
			field.From, field.FromPage = name, page
			doc.Parameters = append(doc.Parameters, field)
			inherited = append(inherited, routeDocInherited{parameter: parameter, always: condition == ""})
		}
	}

	doc.Requests = routeDocRequests(sandbox, conf, inherited)
	for _, example := range conf.Examples {
		if !routeDocGenerated(doc.Requests, example) {
			doc.Examples = append(doc.Examples, example)
		}
	}
	doc.Statuses = routeDocStatuses(sandbox, conf, inherited)
	return doc
}

// routeReachCondition is the cell a crossing is worded with: "" when the
// route in front always runs.
func routeReachCondition(reach utils.Reach, condition string) string {
	switch {
	case reach == utils.MayRun && condition != "":
		return routeDocMayRun + "; " + condition
	case reach == utils.MayRun:
		return routeDocMayRun
	}
	return condition
}

// routeDocGenerated reports whether a declared example is one of the requests
// the page already generates.
func routeDocGenerated(requests []RouteDocRequest, example string) bool {
	for _, request := range requests {
		if request.Command == example {
			return true
		}
	}
	return false
}

// routeDocPath is one path as its row of the address table, and whether it
// gets one: a literal part — an equal or prefix trigger — is already spelled
// by the pattern.
func routeDocPath(sandbox *api.Sandbox, path routeconf.Path) (RouteDocField, bool) {
	row := RouteDocField{Description: docCell(sandbox, path.Description)}

	if !path.Trigger.Exists {
		row.Name = "`" + routeDocPathLabel(path) + "`"
		row.Type = routeDocPathWords(path)
		row.Example = "`" + routeDocPathSample(sandbox, path) + "`"
		return row, true
	}

	literal := path.Trigger.Type == "equal" || path.Trigger.Type == "prefix"
	if literal && !path.Trigger.Negate {
		return row, false
	}

	row.Name = routeDocPlace(sandbox, path)
	row.Type = docCell(sandbox, "text that "+routeDocRule(sandbox, path.Trigger))
	row.Example = "—"
	if sample, ok := routeDocTriggerSample(sandbox, path.Trigger, routeDocSampleFile); ok && !path.Trigger.Negate {
		row.Example = "`" + docCell(sandbox, sandbox.Deps.Stringsdeps.TrimLeft(sample, "/")) + "`"
	}
	return row, true
}

// routeDocPathLabel is how the pattern draws a path that captures its slice.
func routeDocPathLabel(path routeconf.Path) string {
	if path.End == routeconf.LastSegment {
		return "{*" + path.Id + "}"
	}
	if path.Type != routeconf.DefaultPathType {
		return "{" + path.Id + ":" + path.Type + "}"
	}
	return "{" + path.Id + "}"
}

// routeDocPathWords is what a captured path holds, in plain words.
func routeDocPathWords(path routeconf.Path) string {
	if path.End == routeconf.LastSegment || path.End > path.Start {
		return "the rest of the address — one part or more, like `a/b.png`"
	}
	return routeDocTypeWords(path.Type)
}

// routeDocPathSample is the value a generated request puts in a captured
// path.
func routeDocPathSample(sandbox *api.Sandbox, path routeconf.Path) string {
	return routeDocTypeSample(sandbox, path.Id, path.Type)
}

// routeDocPlace names the slice a triggered path reads, counting parts from
// one.
func routeDocPlace(sandbox *api.Sandbox, path routeconf.Path) string {
	first := sandbox.Deps.Stringsdeps.FormatInt(int64(path.Start+1), 10)
	switch {
	case path.Start == 0 && path.End == routeconf.LastSegment:
		return "the whole address"
	case path.End == routeconf.LastSegment:
		return "from part " + first + " on"
	case path.End == path.Start:
		return "part " + first
	}
	return "parts " + first + " to " + sandbox.Deps.Stringsdeps.FormatInt(int64(path.End+1), 10)
}

// routeDocParameter is one parameter as its table row. own is false for a
// parameter a route in front reads: its trigger is a condition on that
// route, not on this one, so it never makes the value required here.
func routeDocParameter(sandbox *api.Sandbox, parameter routeconf.Parameter, own bool) RouteDocField {
	words := routeDocTypeWords(parameter.Type)
	if parameter.Trigger.Exists {
		words += " that " + routeDocRule(sandbox, parameter.Trigger)
	}

	required := "no"
	switch {
	case parameter.Required:
		required = "yes"
	case own && routeDocTriggered(parameter):
		required = "yes — without it this route does not run"
	case parameter.HasDefault:
		required = "no — `" + parameter.Default + "` when left out"
	}

	return RouteDocField{
		Name:        "`" + parameter.Key + "`",
		Where:       routeDocWhere(sandbox, parameter),
		Type:        docCell(sandbox, words),
		Required:    required,
		Example:     "`" + docCell(sandbox, routeDocParameterSample(sandbox, parameter)) + "`",
		Description: docCell(sandbox, parameter.Description),
	}
}

// routeDocTriggered reports whether a parameter carries a trigger its value
// has to meet — one a missing value fails.
func routeDocTriggered(parameter routeconf.Parameter) bool {
	return parameter.Trigger.Exists && !parameter.Trigger.Negate
}

// routeDocFonts is where a parameter is read from, in order; a declaration
// naming none reads the query string alone.
func routeDocFonts(parameter routeconf.Parameter) []string {
	if len(parameter.Fonts) == 0 {
		return []string{"query"}
	}
	return parameter.Fonts
}

// routeDocWhere is where a parameter goes, in plain words.
func routeDocWhere(sandbox *api.Sandbox, parameter routeconf.Parameter) string {
	places := []string{}
	for _, font := range routeDocFonts(parameter) {
		switch font {
		case "query":
			places = append(places, "query string")
		default:
			places = append(places, font)
		}
	}
	return sandbox.Deps.Stringsdeps.Join(places, " or ")
}

// routeDocTypeWords is a path or parameter type in plain words.
func routeDocTypeWords(kind string) string {
	switch kind {
	case "integer":
		return "whole number"
	case "number":
		return "number"
	case "boolean":
		return "`true` or `false`"
	case "uuid":
		return "UUID"
	case "datetime":
		return "date and time, like `" + routeDocSampleDateTime + "`"
	case "string-array":
		return "list of texts"
	case "integer-array":
		return "list of whole numbers"
	}
	return "text"
}

// routeDocTypeSample is the value a generated request sends for a type; a
// text is named after what it is, so the request reads as one.
func routeDocTypeSample(sandbox *api.Sandbox, name string, kind string) string {
	switch kind {
	case "integer", "integer-array":
		return "1"
	case "number":
		return "1.5"
	case "boolean":
		return "true"
	case "uuid":
		return routeDocSampleUuid
	case "datetime":
		return routeDocSampleDateTime
	}
	return "my-" + routeDocWord(sandbox, name)
}

// routeDocParameterSample is the value a generated request sends for a
// parameter: one its trigger accepts, its default, or a sample of its type.
func routeDocParameterSample(sandbox *api.Sandbox, parameter routeconf.Parameter) string {
	fallback := routeDocTypeSample(sandbox, parameter.Key, parameter.Type)
	if routeDocTriggered(parameter) {
		if sample, ok := routeDocTriggerSample(sandbox, parameter.Trigger, fallback); ok {
			return sample
		}
	}
	if parameter.HasDefault {
		return parameter.Default
	}
	return fallback
}

// routeDocTriggerSample is a value a trigger accepts, built around fallback
// when the trigger only bounds one end; false for a regex, which no sample
// is derived from.
func routeDocTriggerSample(sandbox *api.Sandbox, trigger routeconf.Trigger, fallback string) (string, bool) {
	if trigger.Negate {
		return fallback, true
	}
	switch trigger.Type {
	case "equal", "prefix", "text-prefix":
		return trigger.Value, true
	case "suffix":
		return fallback + trigger.Value, true
	case triggerconf.OneOf:
		if len(trigger.Values) > 0 {
			return trigger.Values[0], true
		}
	}
	return "", false
}

// routeDocRule is what a trigger asks of a value, in plain words.
func routeDocRule(sandbox *api.Sandbox, trigger routeconf.Trigger) string {
	verb := "must"
	if trigger.Negate {
		verb = "must not"
	}

	rule := ""
	switch trigger.Type {
	case "equal":
		rule = "be exactly `" + trigger.Value + "`"
	case "prefix", "text-prefix":
		rule = "start with `" + trigger.Value + "`"
	case "suffix":
		rule = "end with `" + trigger.Value + "`"
	case triggerconf.Regex:
		rule = "match the pattern `" + trigger.Value + "`"
	case triggerconf.OneOf:
		rule = "be one of " + routeDocCodeList(sandbox, trigger.Values)
	}

	text := verb + " " + rule
	if trigger.IgnoreCase {
		text += ", in any case"
	}
	return text
}

// routeDocCodeList joins values as code spans: "`a`, `b`".
func routeDocCodeList(sandbox *api.Sandbox, values []string) string {
	quoted := []string{}
	for _, value := range values {
		quoted = append(quoted, "`"+value+"`")
	}
	return sandbox.Deps.Stringsdeps.Join(quoted, ", ")
}

// routeDocWord is a Go id or a header name as the lower-case, dash-separated
// word a sample is spelled with: "UserId" -> "user-id".
func routeDocWord(sandbox *api.Sandbox, name string) string {
	word := []byte{}
	for index := 0; index < len(name); index++ {
		char := name[index]
		if char >= 'A' && char <= 'Z' {
			if index > 0 && routeDocLowerOrDigit(name[index-1]) {
				word = append(word, '-')
			}
			char += 'a' - 'A'
		}
		if char == '_' {
			char = '-'
		}
		word = append(word, char)
	}
	return string(word)
}

// routeDocLowerOrDigit reports whether a byte ends a word a capital starts a
// new one after.
func routeDocLowerOrDigit(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')
}

// routeDocSamplePath is a request path every path of the route accepts,
// built slot by slot the way the pattern is drawn: the segments an equal,
// prefix, text-prefix or one-of trigger spells first — the first to reach a
// slot wins — then a sample of each capture's type in the slots left, then
// each suffix appended to the last segment of its slice. A slot nothing fills
// reads "anything", and a route declaring a segment count gets that many. A
// regex or a negated trigger is not sampled.
func routeDocSamplePath(sandbox *api.Sandbox, conf *routeconf.RouteConf) string {
	strings := sandbox.Deps.Stringsdeps
	slots := map[int]string{}
	last := -1
	place := func(index int, text string) {
		if _, taken := slots[index]; taken {
			return
		}
		slots[index] = text
		if index > last {
			last = index
		}
	}
	spell := func(start int, value string) {
		index := start
		for _, segment := range strings.Split(value, "/") {
			if segment == "" {
				continue
			}
			place(index, segment)
			index++
		}
	}

	for _, path := range conf.Paths {
		trigger := path.Trigger
		if !trigger.Exists || trigger.Negate {
			continue
		}
		switch trigger.Type {
		case "equal", "prefix", "text-prefix":
			spell(path.Start, trigger.Value)
		case triggerconf.OneOf:
			if len(trigger.Values) > 0 {
				spell(path.Start, trigger.Values[0])
			}
		}
	}
	for _, path := range conf.Paths {
		if path.Trigger.Exists {
			continue
		}
		end := path.End
		if end == routeconf.LastSegment {
			end = path.Start
		}
		for index := path.Start; index <= end; index++ {
			place(index, routeDocPathSample(sandbox, path))
		}
	}
	for _, path := range conf.Paths {
		trigger := path.Trigger
		if !trigger.Exists || trigger.Negate || trigger.Type != "suffix" {
			continue
		}
		end := path.End
		if end == routeconf.LastSegment {
			end = path.Start
			if last > end {
				end = last
			}
		}
		place(end, routeDocSampleFile)
		if !strings.HasSuffix(slots[end], trigger.Value) {
			slots[end] += trigger.Value
		}
	}

	count := last + 1
	if conf.HasSegments && conf.Segments > count {
		count = conf.Segments
	}
	segments := []string{}
	for index := 0; index < count; index++ {
		text, taken := slots[index]
		if !taken {
			text = "anything"
		}
		segments = append(segments, text)
	}
	return "/" + strings.Join(segments, "/")
}

// routeDocRequests is the page's ready-to-run calls: the smallest request the
// route answers and, when the route reads more than that, one sending
// everything it reads.
func routeDocRequests(sandbox *api.Sandbox, conf *routeconf.RouteConf, inherited []routeDocInherited) []RouteDocRequest {
	smallest := routeDocCurl(sandbox, conf, inherited, false)
	complete := routeDocCurl(sandbox, conf, inherited, true)
	if smallest == complete {
		return []RouteDocRequest{{Command: complete}}
	}
	return []RouteDocRequest{
		{Title: "Only what is required:", Command: smallest},
		{Title: "With every value it reads:", Command: complete},
	}
}

// routeDocCurl spells one request as a curl command. complete sends every
// parameter and every body property; otherwise only the ones the route
// cannot run without.
func routeDocCurl(sandbox *api.Sandbox, conf *routeconf.RouteConf, inherited []routeDocInherited, complete bool) string {
	strings := sandbox.Deps.Stringsdeps

	head := "curl"
	method := routeconf.DefaultMethod
	if len(conf.Methods) > 0 {
		method = conf.Methods[0]
	}
	switch method {
	case routeconf.AnyMethod, "GET":
	case "HEAD":
		head += " -I"
	default:
		head += " -X " + method
	}

	query, cookies, lines := []string{}, []string{}, []string{}
	send := func(parameter routeconf.Parameter, needed bool) {
		if !complete && !needed {
			return
		}
		value := routeDocParameterSample(sandbox, parameter)
		switch routeDocFonts(parameter)[0] {
		case "header":
			lines = append(lines, "-H "+routeDocShellQuote(sandbox, parameter.Key+": "+value))
		case "cookie":
			cookies = append(cookies, parameter.Key+"="+value)
		default:
			query = append(query, routeDocUrlEscape(sandbox, parameter.Key)+"="+routeDocUrlEscape(sandbox, value))
		}
	}
	for _, parameter := range conf.Parameters {
		send(parameter, parameter.Required || routeDocTriggered(parameter))
	}
	for _, other := range inherited {
		send(other.parameter, other.always && other.parameter.Required)
	}
	if len(cookies) > 0 {
		lines = append(lines, "-b "+routeDocShellQuote(sandbox, strings.Join(cookies, "; ")))
	}

	body := conf.Body
	if body.Type != routeconf.BodyNone && (complete || body.Required) {
		if body.ContentType != "" {
			lines = append(lines, "-H "+routeDocShellQuote(sandbox, "Content-Type: "+body.ContentType))
		}
		switch body.Type {
		case "json":
			sample := "{}"
			if body.HasSchema {
				sample = routeDocJsonText(sandbox, routeDocJsonSample(sandbox, body.Schema, complete), "", false)
			}
			lines = append(lines, "-d "+routeDocShellQuote(sandbox, sample))
		case "form":
			lines = append(lines, "-d "+routeDocShellQuote(sandbox, "name=value"))
		case "text":
			lines = append(lines, "-d "+routeDocShellQuote(sandbox, "some text"))
		default:
			lines = append(lines, "--data-binary @file.bin")
		}
	}

	url := routeDocHost + routeDocSamplePath(sandbox, conf)
	if len(query) > 0 {
		url += "?" + strings.Join(query, "&")
	}
	head += " " + routeDocShellWord(sandbox, url)

	if len(lines) == 0 {
		return head
	}
	return head + " \\\n  " + strings.Join(lines, " \\\n  ")
}

// routeDocShellQuote wraps a value in single quotes for a shell, a quote
// inside it closed, escaped and reopened.
func routeDocShellQuote(sandbox *api.Sandbox, value string) string {
	return "'" + sandbox.Deps.Stringsdeps.ReplaceAll(value, "'", `'\''`) + "'"
}

// routeDocShellWord quotes a url only when a shell would read something in
// it: `curl localhost:3000/health` stays as a person types it.
func routeDocShellWord(sandbox *api.Sandbox, value string) string {
	for index := 0; index < len(value); index++ {
		char := value[index]
		plain := routeDocLowerOrDigit(char) || (char >= 'A' && char <= 'Z') ||
			sandbox.Deps.Stringsdeps.ContainsAny(string(char), "-._/:%@,+")
		if !plain {
			return routeDocShellQuote(sandbox, value)
		}
	}
	return value
}

// routeDocUrlEscape escapes the characters that would end or split a query
// value.
func routeDocUrlEscape(sandbox *api.Sandbox, value string) string {
	strings := sandbox.Deps.Stringsdeps
	for _, pair := range [][2]string{{"%", "%25"}, {" ", "%20"}, {"&", "%26"}, {"#", "%23"}, {"+", "%2B"}, {"=", "%3D"}} {
		value = strings.ReplaceAll(value, pair[0], pair[1])
	}
	return value
}

// routeDocBody is what the page says about a route's body, nil when it
// declares none.
func routeDocBody(sandbox *api.Sandbox, body routeconf.Body) *RouteDocBody {
	if body.Type == routeconf.BodyNone {
		return nil
	}

	kind := "any data — a file, for example"
	switch body.Type {
	case "json":
		kind = "JSON"
	case "text":
		kind = "plain text"
	case "form":
		kind = "form fields, like `name=value&other=value`"
	}

	summary := "Send " + kind
	if body.ContentType != "" {
		summary += " with the header `Content-Type: " + body.ContentType + "`"
	}
	summary += ", up to " + routeDocSize(sandbox, body.MaxBytes) + "."
	if body.Required {
		summary += " The body is required."
	} else {
		summary += " The body is optional."
	}

	doc := &RouteDocBody{Summary: summary}
	if body.Type == "json" && body.HasSchema {
		doc.Fields = routeDocBodyFields(sandbox, body.Schema, "", nil)
		doc.Sample = routeDocJsonText(sandbox, routeDocJsonSample(sandbox, body.Schema, true), "", true)
		doc.SampleLang = "json"
	}
	return doc
}

// routeDocSize is a byte count as a person reads it: 1048576 -> "1 MB".
func routeDocSize(sandbox *api.Sandbox, bytes int) string {
	format := sandbox.Deps.Stringsdeps.FormatInt
	switch {
	case bytes >= 1048576 && bytes%1048576 == 0:
		return format(int64(bytes/1048576), 10) + " MB"
	case bytes >= 1024 && bytes%1024 == 0:
		return format(int64(bytes/1024), 10) + " KB"
	}
	return format(int64(bytes), 10) + " bytes"
}

// routeDocBodyFields flattens a schema into one row per property, a nested
// one under its dotted path and a property of the objects of a list under
// `list[].property`. A body that is not an object is one row, the whole body.
func routeDocBodyFields(sandbox *api.Sandbox, schema *routeconf.Schema, prefix string, rows []RouteDocBodyField) []RouteDocBodyField {
	if prefix == "" && !routeDocIsObject(schema) {
		return append(rows, RouteDocBodyField{
			Name:     "the whole body",
			Type:     routeDocSchemaWords(schema),
			Required: "yes",
			Rules:    routeDocSchemaRules(sandbox, schema),
		})
	}

	for _, property := range schema.Properties {
		name := prefix + property.Name
		required := "no"
		if routeDocContains(schema.Required, property.Name) {
			required = "yes"
		}
		rows = append(rows, RouteDocBodyField{
			Name:     "`" + name + "`",
			Type:     routeDocSchemaWords(property.Schema),
			Required: required,
			Rules:    routeDocSchemaRules(sandbox, property.Schema),
		})

		child := property.Schema
		switch {
		case routeDocIsObject(child):
			rows = routeDocBodyFields(sandbox, child, name+".", rows)
		case child.Type == "array" && child.Items != nil && routeDocIsObject(child.Items):
			rows = routeDocBodyFields(sandbox, child.Items, name+"[].", rows)
		}
	}
	return rows
}

// routeDocIsObject reports whether a schema node is an object: declared one,
// or carrying properties with no type at all.
func routeDocIsObject(schema *routeconf.Schema) bool {
	return schema.Type == "object" || (schema.Type == "" && len(schema.Properties) > 0)
}

// routeDocContains reports whether values holds value.
func routeDocContains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

// routeDocSchemaWords is what a schema node holds, in plain words.
func routeDocSchemaWords(schema *routeconf.Schema) string {
	one, _ := routeDocSchemaNoun(schema)
	return one
}

// routeDocSchemaNoun is a schema node's noun, singular and plural, so a list
// reads "list of whole numbers".
func routeDocSchemaNoun(schema *routeconf.Schema) (string, string) {
	switch {
	case schema.Type == "array" && schema.Items != nil:
		_, many := routeDocSchemaNoun(schema.Items)
		return "list of " + many, "lists"
	case schema.Type == "array":
		return "list", "lists"
	case routeDocIsObject(schema):
		return "object", "objects"
	case schema.Type == "integer":
		return "whole number", "whole numbers"
	case schema.Type == "number":
		return "number", "numbers"
	case schema.Type == "boolean":
		return "`true` or `false`", "`true` or `false` values"
	case schema.Type == "null":
		return "`null`", "`null` values"
	}

	switch schema.Format {
	case "email":
		return "e-mail address", "e-mail addresses"
	case "uuid":
		return "UUID", "UUIDs"
	case "date-time":
		return "date and time", "dates and times"
	case "uri":
		return "link (URL)", "links (URL)"
	}
	return "text", "texts"
}

// routeDocSchemaRules is every bound of a schema node, in plain words.
func routeDocSchemaRules(sandbox *api.Sandbox, schema *routeconf.Schema) string {
	number := func(value float64) string {
		return sandbox.Deps.Stringsdeps.FormatFloat(value, 'f', -1, 64)
	}
	count := func(value int) string {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}

	rules := []string{}
	if schema.HasConst {
		rules = append(rules, "always `"+schema.Const+"`")
	}
	if len(schema.Enum) > 0 {
		rules = append(rules, "one of "+routeDocCodeList(sandbox, schema.Enum))
	}

	switch {
	case schema.HasMinimum && schema.HasMaximum:
		rules = append(rules, "from "+number(schema.Minimum)+" to "+number(schema.Maximum))
	case schema.HasMinimum:
		rules = append(rules, "at least "+number(schema.Minimum))
	case schema.HasMaximum:
		rules = append(rules, "at most "+number(schema.Maximum))
	}
	if schema.HasExclusiveMinimum {
		rules = append(rules, "more than "+number(schema.ExclusiveMinimum))
	}
	if schema.HasExclusiveMaximum {
		rules = append(rules, "less than "+number(schema.ExclusiveMaximum))
	}

	switch {
	case schema.HasMinLength && schema.HasMaxLength:
		rules = append(rules, count(schema.MinLength)+" to "+count(schema.MaxLength)+" characters")
	case schema.HasMinLength:
		rules = append(rules, "at least "+count(schema.MinLength)+" characters")
	case schema.HasMaxLength:
		rules = append(rules, "at most "+count(schema.MaxLength)+" characters")
	}
	if schema.Pattern != "" {
		rules = append(rules, "must match the pattern `"+schema.Pattern+"`")
	}

	switch {
	case schema.HasMinItems && schema.HasMaxItems:
		rules = append(rules, count(schema.MinItems)+" to "+count(schema.MaxItems)+" items")
	case schema.HasMinItems:
		rules = append(rules, "at least "+count(schema.MinItems)+" items")
	case schema.HasMaxItems:
		rules = append(rules, "at most "+count(schema.MaxItems)+" items")
	}
	if schema.UniqueItems {
		rules = append(rules, "no item repeated")
	}
	if schema.Nullable {
		rules = append(rules, "may be `null`")
	}
	if schema.HasAdditionalProperties && !schema.AdditionalProperties {
		rules = append(rules, "no field beyond the ones listed")
	}

	return docCell(sandbox, sandbox.Deps.Stringsdeps.Join(rules, "; "))
}

// routeDocJson is one node of a sample json value: a leaf holding its literal
// text, or an object or a list of nodes.
type routeDocJson struct {
	Literal  string
	Object   bool
	List     bool
	Keys     []string
	Children []*routeDocJson
}

// routeDocJsonSample is a value the schema accepts. complete fills every
// property of an object; otherwise only the required ones.
func routeDocJsonSample(sandbox *api.Sandbox, schema *routeconf.Schema, complete bool) *routeDocJson {
	switch {
	case routeDocIsObject(schema):
		node := &routeDocJson{Object: true}
		for _, property := range schema.Properties {
			if !complete && !routeDocContains(schema.Required, property.Name) {
				continue
			}
			node.Keys = append(node.Keys, property.Name)
			node.Children = append(node.Children, routeDocJsonSample(sandbox, property.Schema, complete))
		}
		return node
	case schema.Type == "array":
		node := &routeDocJson{List: true}
		if schema.Items == nil {
			return node
		}
		items := schema.MinItems
		if items < 1 {
			items = 1
		}
		for index := 0; index < items; index++ {
			node.Children = append(node.Children, routeDocJsonSample(sandbox, schema.Items, complete))
		}
		return node
	}
	return &routeDocJson{Literal: routeDocJsonLeaf(sandbox, schema)}
}

// routeDocJsonLeaf is the literal json text of a scalar the schema accepts:
// its const, its first enum value, or a sample of its type inside its bounds.
func routeDocJsonLeaf(sandbox *api.Sandbox, schema *routeconf.Schema) string {
	strings := sandbox.Deps.Stringsdeps
	quoted := schema.Type == "string" || schema.Type == ""
	literal := func(value string) string {
		if quoted {
			return strings.Quote(value)
		}
		return value
	}

	if schema.HasConst {
		return literal(schema.Const)
	}
	if len(schema.Enum) > 0 {
		return literal(schema.Enum[0])
	}

	switch schema.Type {
	case "integer":
		return strings.FormatInt(int64(routeDocBounded(schema, 1, 1)), 10)
	case "number":
		return strings.FormatFloat(routeDocBounded(schema, 1.5, 0.5), 'f', -1, 64)
	case "boolean":
		return "true"
	case "null":
		return "null"
	}

	text := routeDocSampleText
	switch schema.Format {
	case "email":
		text = routeDocSampleEmail
	case "uuid":
		text = routeDocSampleUuid
	case "date-time":
		text = routeDocSampleDateTime
	case "uri":
		text = routeDocSampleUri
	}
	if schema.HasMinLength && len(text) < schema.MinLength {
		text += strings.Repeat("x", schema.MinLength-len(text))
	}
	if schema.HasMaxLength && len(text) > schema.MaxLength {
		text = text[:schema.MaxLength]
	}
	return strings.Quote(text)
}

// routeDocBounded moves value inside a number's bounds; step is how far past
// an exclusive bound it lands.
func routeDocBounded(schema *routeconf.Schema, value float64, step float64) float64 {
	if schema.HasMinimum && value < schema.Minimum {
		value = schema.Minimum
	}
	if schema.HasExclusiveMinimum && value <= schema.ExclusiveMinimum {
		value = schema.ExclusiveMinimum + step
	}
	if schema.HasMaximum && value > schema.Maximum {
		value = schema.Maximum
	}
	if schema.HasExclusiveMaximum && value >= schema.ExclusiveMaximum {
		value = schema.ExclusiveMaximum - step
	}
	return value
}

// routeDocJsonText writes a sample json value: indented two spaces a level
// when pretty, on one line otherwise.
func routeDocJsonText(sandbox *api.Sandbox, node *routeDocJson, indent string, pretty bool) string {
	if !node.Object && !node.List {
		return node.Literal
	}

	opening, closing := "[", "]"
	if node.Object {
		opening, closing = "{", "}"
	}
	if len(node.Children) == 0 {
		return opening + closing
	}

	inner, separator, colon, newline := indent+"  ", ",", ":", ""
	if pretty {
		newline = "\n"
		colon = ": "
	} else {
		inner, indent = "", ""
	}

	items := []string{}
	for index, child := range node.Children {
		text := routeDocJsonText(sandbox, child, inner, pretty)
		if node.Object {
			text = sandbox.Deps.Stringsdeps.Quote(node.Keys[index]) + colon + text
		}
		items = append(items, inner+text)
	}
	return opening + newline + sandbox.Deps.Stringsdeps.Join(items, separator+newline) + newline + indent + closing
}

// routeDocStatuses is every status this route itself answers with, in plain
// words; the ones any route may answer are on docs/Routes' own page.
func routeDocStatuses(sandbox *api.Sandbox, conf *routeconf.RouteConf, inherited []routeDocInherited) []RouteDocStatus {
	success := "It worked."
	if conf.ResponseType != "" {
		success += " The answer comes as `" + conf.ResponseType + "`."
	}
	statuses := []RouteDocStatus{{Code: "`200`", Meaning: success}}

	bad := conf.Body.Type != routeconf.BodyNone
	for _, parameter := range conf.Parameters {
		bad = bad || parameter.Required || parameter.Type != "string"
	}
	for _, other := range inherited {
		bad = bad || other.parameter.Required || other.parameter.Type != "string"
	}
	if bad {
		statuses = append(statuses, RouteDocStatus{
			Code:    "`400`",
			Meaning: "Something you sent is missing or has the wrong type or format. The answer's `field` names it.",
		})
	}

	for _, path := range conf.Paths {
		if !path.Trigger.Exists && path.Type != routeconf.DefaultPathType {
			statuses = append(statuses, RouteDocStatus{
				Code: "`404`",
				Meaning: "`" + routeDocPathLabel(path) + "` is not a " + routeDocTypeWords(path.Type) +
					", so this route does not answer the address.",
			})
			break
		}
	}

	if conf.Body.Type != routeconf.BodyNone {
		statuses = append(statuses, RouteDocStatus{
			Code:    "`413`",
			Meaning: "The body is larger than " + routeDocSize(sandbox, conf.Body.MaxBytes) + ".",
		})
		if conf.Body.ContentType != "" {
			statuses = append(statuses, RouteDocStatus{
				Code:    "`415`",
				Meaning: "The body was not sent with `Content-Type: " + conf.Body.ContentType + "`.",
			})
		}
	}
	return statuses
}
