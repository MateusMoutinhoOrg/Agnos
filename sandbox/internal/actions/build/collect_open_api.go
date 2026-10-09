package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// openApiVersion is the OpenAPI release the document is written in: the one
// Postman, every Swagger UI and the code generators read.
const openApiVersion = "3.0.3"

// openApiUnversioned is the document's version while the project's
// project.yaml sets none: OpenAPI requires one.
const openApiUnversioned = "unversioned"

// openApiErrorSchema is the component every failure status answers with: the
// json the server/errors handlers write unless the project changed them.
const openApiErrorSchema = "Error"

// The media types a body is read as when its route declares no content-type.
const (
	openApiJsonMedia   = "application/json"
	openApiFormMedia   = "application/x-www-form-urlencoded"
	openApiTextMedia   = "text/plain"
	openApiBinaryMedia = "application/octet-stream"
)

// CollectOpenApi renders every visible route.yaml into one OpenAPI 3.0.3
// document — the text docs/Routes/openapi.json holds and the generated openapi
// route answers — so whoever calls the server imports it into Postman or opens
// it in Swagger UI. It reads the declarations docs/Routes reads, crossed the
// same way, and nothing else.
//
// Routes are laid down in run order. A route declaring ANY is no operation —
// OpenAPI has no method that stands for every one, and it is how a middleware
// is declared — but the parameters of one sure to run in front of a route join
// that route's operation, and an authorization header becomes a security
// scheme. When two routes draw the same path and method, the one running last
// is the one kept: it answers when the ones in front decline.
func CollectOpenApi(sandbox *api.Sandbox, io *stagedfs.StagedFS, title string, version string) (string, error) {
	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return "", err
	}

	paths := openApiObject()
	tags := openApiList()
	tagged := map[string]bool{}
	schemes := openApiObject()

	for _, current := range chain {
		conf := current.Conf
		if conf.Hidden || conf.Private || routeDocContains(conf.Methods, routeconf.AnyMethod) {
			continue
		}

		category := routeDocCategory(conf)
		if !tagged[category] {
			tagged[category] = true
			tag := openApiObject()
			openApiSet(tag, "name", openApiString(sandbox, category))
			openApiAdd(tags, tag)
		}

		_, inherited := routeDocCrossings(sandbox, current, chain)
		template, parameters, notes := openApiPath(sandbox, conf)

		item := openApiGet(paths, template)
		if item == nil {
			item = openApiObject()
			openApiSet(paths, template, item)
		}
		for _, method := range conf.Methods {
			operation := openApiOperation(sandbox, current, method, category, parameters, notes, inherited, schemes)
			openApiSet(item, sandbox.Deps.StringsDeps.ToLower(method), operation)
		}
	}

	if version == "" {
		version = openApiUnversioned
	}
	info := openApiObject()
	openApiSet(info, "title", openApiString(sandbox, title))
	openApiSet(info, "description", openApiString(sandbox, "Every address "+title+" answers: what to send and what comes back."))
	openApiSet(info, "version", openApiString(sandbox, version))

	components := openApiObject()
	schemas := openApiObject()
	openApiSet(schemas, openApiErrorSchema, openApiErrorBody(sandbox))
	openApiSet(components, "schemas", schemas)
	if len(schemes.Keys) > 0 {
		openApiSet(components, "securitySchemes", schemes)
	}

	document := openApiObject()
	openApiSet(document, "openapi", openApiString(sandbox, openApiVersion))
	openApiSet(document, "info", info)
	openApiSet(document, "servers", openApiServers(sandbox))
	openApiSet(document, "tags", tags)
	openApiSet(document, "paths", paths)
	openApiSet(document, "components", components)

	return routeDocJsonText(sandbox, document, "", true) + "\n", nil
}

// openApiServers is where the requests go: localhost:3000 by default, the
// address the route pages call, with the scheme and the host left for the
// caller to change.
func openApiServers(sandbox *api.Sandbox) *routeDocJson {
	scheme := openApiObject()
	openApiSet(scheme, "default", openApiString(sandbox, "http"))
	openApiSet(scheme, "enum", openApiStrings(sandbox, []string{"http", "https"}))

	host := openApiObject()
	openApiSet(host, "default", openApiString(sandbox, routeDocHost))
	openApiSet(host, "description", openApiString(sandbox, "the host and port the server listens on"))

	variables := openApiObject()
	openApiSet(variables, "scheme", scheme)
	openApiSet(variables, "host", host)

	server := openApiObject()
	openApiSet(server, "url", openApiString(sandbox, "{scheme}://{host}"))
	openApiSet(server, "variables", variables)

	servers := openApiList()
	openApiAdd(servers, server)
	return servers
}

// openApiErrorBody is the json every failure is answered with: what went wrong
// and, when it is one value, which one.
func openApiErrorBody(sandbox *api.Sandbox) *routeDocJson {
	field := func(description string) *routeDocJson {
		node := openApiObject()
		openApiSet(node, "type", openApiString(sandbox, "string"))
		openApiSet(node, "description", openApiString(sandbox, description))
		return node
	}
	properties := openApiObject()
	openApiSet(properties, "error", field("what went wrong"))
	openApiSet(properties, "field", field("the value that was wrong, when it is one"))

	node := openApiObject()
	openApiSet(node, "type", openApiString(sandbox, "object"))
	openApiSet(node, "properties", properties)
	openApiSet(node, "required", openApiStrings(sandbox, []string{"error"}))
	return node
}

// openApiPath is a route's address as an OpenAPI path template, the path
// parameters it captures and the notes on what the template cannot draw. It is
// laid down slot by slot the way routeconf.Pattern draws it: the segments an
// equal or a prefix trigger spells first, then one {Id} per capture, one-of,
// text-prefix, suffix or regex at the slot it starts on. A prefix lists the
// route at the prefix itself, a negated trigger is a note, and a slot nothing
// reaches is a {segmentN} of its own.
func openApiPath(sandbox *api.Sandbox, conf *routeconf.RouteConf) (string, []*routeDocJson, []string) {
	strings := sandbox.Deps.StringsDeps
	slots := map[int]string{}
	last := -1
	parameters := []*routeDocJson{}
	notes := []string{}

	place := func(index int, text string) bool {
		if _, taken := slots[index]; taken {
			return false
		}
		slots[index] = text
		if index > last {
			last = index
		}
		return true
	}

	for _, path := range conf.Paths {
		trigger := path.Trigger
		if !trigger.Set || trigger.Negate || (trigger.Type != "equal" && trigger.Type != "prefix") {
			continue
		}
		index := path.Start
		for _, segment := range strings.Split(trigger.Value, "/") {
			if segment == "" {
				continue
			}
			place(index, segment)
			index++
		}
		if trigger.Type == "prefix" {
			notes = append(notes, "Also answers every address under `"+trigger.Value+"`.")
		}
	}

	for _, path := range conf.Paths {
		trigger := path.Trigger
		if trigger.Set && trigger.Negate {
			notes = append(notes, openApiCapital(sandbox, routeDocPlace(sandbox, path)+" "+routeDocRule(sandbox, trigger)+"."))
			continue
		}
		if trigger.Set && (trigger.Type == "equal" || trigger.Type == "prefix") {
			continue
		}
		if !place(path.Start, "{"+path.Id+"}") {
			continue
		}
		if !trigger.Set && path.End != routeconf.LastSegment {
			for index := path.Start + 1; index <= path.End; index++ {
				place(index, "")
			}
		}
		parameters = append(parameters, openApiPathParameter(sandbox, path))
	}

	if conf.HasSegments && conf.Segments-1 > last {
		last = conf.Segments - 1
	}
	segments := []string{}
	for index := 0; index <= last; index++ {
		text, taken := slots[index]
		if !taken {
			name := "segment" + strings.FormatInt(int64(index+1), 10)
			text = "{" + name + "}"
			parameters = append(parameters, openApiFreeSegment(sandbox, name, index))
		}
		if text != "" {
			segments = append(segments, text)
		}
	}
	return "/" + strings.Join(segments, "/"), parameters, notes
}

// openApiPathParameter is one path the template draws as {Id}: its type, the
// values a one-of allows, the pattern a regex holds it to, and every other
// rule in words.
func openApiPathParameter(sandbox *api.Sandbox, path routeconf.Path) *routeDocJson {
	trigger := path.Trigger
	schema := openApiTypeSchema(sandbox, path.Type)
	words := ""
	example := openApiTyped(sandbox, path.Type, routeDocPathSample(sandbox, path))

	switch {
	case !trigger.Set && (path.End == routeconf.LastSegment || path.End > path.Start):
		words = openApiCapital(sandbox, routeDocPathWords(path)) + "."
	case trigger.Type == triggerconf.OneOf && !trigger.IgnoreCase:
		values := []string{}
		for _, value := range trigger.Values {
			values = append(values, sandbox.Deps.StringsDeps.TrimLeft(value, "/"))
		}
		openApiSet(schema, "enum", openApiStrings(sandbox, values))
		if len(values) > 0 {
			example = openApiString(sandbox, values[0])
		}
	case trigger.Set:
		words = openApiCapital(sandbox, routeDocPlace(sandbox, path)+" "+routeDocRule(sandbox, trigger)) + "."
		if trigger.Type == triggerconf.Regex && !trigger.IgnoreCase {
			openApiSet(schema, "pattern", openApiString(sandbox, trigger.Value))
		}
		example = nil
		if sample, ok := routeDocTriggerSample(sandbox, trigger, routeDocSampleFile); ok {
			example = openApiString(sandbox, sandbox.Deps.StringsDeps.TrimLeft(sample, "/"))
		}
	}

	node := openApiObject()
	openApiSet(node, "name", openApiString(sandbox, path.Id))
	openApiSet(node, "in", openApiString(sandbox, "path"))
	if description := openApiText(path.Description, words); description != "" {
		openApiSet(node, "description", openApiString(sandbox, description))
	}
	openApiSet(node, "required", openApiBool(true))
	openApiSet(node, "schema", schema)
	if example != nil {
		openApiSet(node, "example", example)
	}
	return node
}

// openApiFreeSegment is a part of the address no path of the route reaches:
// any text is answered there.
func openApiFreeSegment(sandbox *api.Sandbox, name string, index int) *routeDocJson {
	node := openApiObject()
	openApiSet(node, "name", openApiString(sandbox, name))
	openApiSet(node, "in", openApiString(sandbox, "path"))
	openApiSet(node, "description", openApiString(sandbox, "Part "+sandbox.Deps.StringsDeps.FormatInt(int64(index+1), 10)+" of the address: any text."))
	openApiSet(node, "required", openApiBool(true))
	openApiSet(node, "schema", openApiTypeSchema(sandbox, routeconf.DefaultPathType))
	return node
}

// openApiOperation is one method of a route: what its page says, the values
// it and the routes always in front of it read, its body and its statuses.
func openApiOperation(sandbox *api.Sandbox, current utils.RouteChainEntry, method string, category string, pathParameters []*routeDocJson, notes []string, inherited []routeDocInherited, schemes *routeDocJson) *routeDocJson {
	strings := sandbox.Deps.StringsDeps
	conf := current.Conf
	name := utils.RouteName(sandbox, current.Name)

	operation := openApiObject()
	id := name
	if len(conf.Methods) > 1 {
		id += "-" + strings.ToLower(method)
	}
	openApiSet(operation, "operationId", openApiString(sandbox, id))
	openApiSet(operation, "tags", openApiStrings(sandbox, []string{category}))
	if conf.Summary != "" {
		openApiSet(operation, "summary", openApiString(sandbox, conf.Summary))
	}
	if description := openApiText(append([]string{conf.Description}, notes...)...); description != "" {
		openApiSet(operation, "description", openApiString(sandbox, description))
	}

	parameters := openApiList()
	parameters.Children = append(parameters.Children, pathParameters...)
	seen := map[string]bool{}
	requirement := openApiObject()
	required := false

	read := func(parameter routeconf.Parameter, needed bool, from string, own bool) {
		in := routeDocSources(parameter)[0]
		key := parameter.Key
		if in == "header" {
			key = strings.ToLower(key)
		}
		if seen[in+":"+key] {
			return
		}
		seen[in+":"+key] = true

		if in == "header" && key == "authorization" {
			scheme := openApiObject()
			openApiSet(scheme, "type", openApiString(sandbox, "apiKey"))
			openApiSet(scheme, "in", openApiString(sandbox, "header"))
			openApiSet(scheme, "name", openApiString(sandbox, parameter.Key))
			if parameter.Description != "" {
				openApiSet(scheme, "description", openApiString(sandbox, parameter.Description))
			}
			openApiSet(schemes, from, scheme)
			openApiSet(requirement, from, openApiList())
			required = required || needed
			return
		}
		if in == "header" && (key == "accept" || key == "content-type") {
			return
		}
		openApiAdd(parameters, openApiParameter(sandbox, parameter, in, needed, from, own))
	}
	for _, parameter := range conf.Parameters {
		read(parameter, parameter.Required || routeDocTriggered(parameter), name, true)
	}
	for _, other := range inherited {
		read(other.parameter, other.always && other.parameter.Required, other.from, false)
	}

	if len(parameters.Children) > 0 {
		openApiSet(operation, "parameters", parameters)
	}
	if body := openApiBody(sandbox, conf.Body); body != nil {
		openApiSet(operation, "requestBody", body)
	}
	openApiSet(operation, "responses", openApiResponses(sandbox, conf, inherited))

	if len(requirement.Keys) > 0 {
		security := openApiList()
		if !required {
			openApiAdd(security, openApiObject())
		}
		openApiAdd(security, requirement)
		openApiSet(operation, "security", security)
	}
	return operation
}

// openApiParameter is one query, header or cookie value: read from the first
// of its sources, required when the route cannot run without it, and the
// value the route page's requests send as its example. own is false for a
// value a route in front reads, which the description names.
func openApiParameter(sandbox *api.Sandbox, parameter routeconf.Parameter, in string, required bool, from string, own bool) *routeDocJson {
	array := sandbox.Deps.StringsDeps.HasSuffix(parameter.Type, "-array")
	schema := openApiTypeSchema(sandbox, parameter.Type)

	rule := ""
	if parameter.Trigger.Set {
		rule = "The value " + routeDocRule(sandbox, parameter.Trigger) + "."
		trigger := parameter.Trigger
		if !array && !trigger.Negate && !trigger.IgnoreCase {
			switch trigger.Type {
			case "equal":
				openApiSet(schema, "enum", openApiTypedList(sandbox, parameter.Type, []string{trigger.Value}))
			case triggerconf.OneOf:
				openApiSet(schema, "enum", openApiTypedList(sandbox, parameter.Type, trigger.Values))
			case triggerconf.Regex:
				openApiSet(schema, "pattern", openApiString(sandbox, trigger.Value))
			}
		}
	}
	if parameter.HasDefault && !array {
		openApiSet(schema, "default", openApiTyped(sandbox, parameter.Type, parameter.Default))
	}
	origin := ""
	if !own {
		origin = "Read by `" + from + "`, which runs in front of this route."
	}

	node := openApiObject()
	openApiSet(node, "name", openApiString(sandbox, parameter.Key))
	openApiSet(node, "in", openApiString(sandbox, in))
	if description := openApiText(parameter.Description, rule, origin); description != "" {
		openApiSet(node, "description", openApiString(sandbox, description))
	}
	if required {
		openApiSet(node, "required", openApiBool(true))
	}
	openApiSet(node, "schema", schema)

	example := openApiTyped(sandbox, parameter.Type, routeDocParameterSample(sandbox, parameter))
	if array {
		example = openApiList()
		openApiAdd(example, openApiTyped(sandbox, parameter.Type, routeDocParameterSample(sandbox, parameter)))
	}
	openApiSet(node, "example", example)
	return node
}

// openApiTypeSchema is a path or parameter type as an OpenAPI schema.
func openApiTypeSchema(sandbox *api.Sandbox, kind string) *routeDocJson {
	node := openApiObject()
	switch kind {
	case "integer", "number", "boolean":
		openApiSet(node, "type", openApiString(sandbox, kind))
	case "uuid":
		openApiSet(node, "type", openApiString(sandbox, "string"))
		openApiSet(node, "format", openApiString(sandbox, "uuid"))
	case "datetime":
		openApiSet(node, "type", openApiString(sandbox, "string"))
		openApiSet(node, "format", openApiString(sandbox, "date-time"))
	case "string-array", "integer-array":
		openApiSet(node, "type", openApiString(sandbox, "array"))
		openApiSet(node, "items", openApiTypeSchema(sandbox, sandbox.Deps.StringsDeps.TrimSuffix(kind, "-array")))
	default:
		openApiSet(node, "type", openApiString(sandbox, "string"))
	}
	return node
}

// openApiBody is a route's request body, nil when it declares none: its
// media type, whether it may be left out, how large it may be and, for json
// and form, its schema — with the json sample the route page shows.
func openApiBody(sandbox *api.Sandbox, body routeconf.Body) *routeDocJson {
	if body.Type == routeconf.BodyNone {
		return nil
	}

	media := openApiObject()
	schema := openApiObject()
	kind := openApiBinaryMedia
	switch body.Type {
	case "json":
		kind = openApiJsonMedia
		if body.HasSchema {
			schema = openApiSchema(sandbox, body.Schema)
		}
	case "form":
		kind = openApiFormMedia
		if body.HasSchema {
			schema = openApiSchema(sandbox, body.Schema)
		} else {
			openApiSet(schema, "type", openApiString(sandbox, "object"))
		}
	case "text":
		kind = openApiTextMedia
		openApiSet(schema, "type", openApiString(sandbox, "string"))
	default:
		openApiSet(schema, "type", openApiString(sandbox, "string"))
		openApiSet(schema, "format", openApiString(sandbox, "binary"))
	}
	if body.ContentType != "" {
		kind = body.ContentType
	}
	openApiSet(media, "schema", schema)
	if body.Type == "json" && body.HasSchema {
		openApiSet(media, "example", routeDocJsonSample(sandbox, body.Schema, true))
	}

	content := openApiObject()
	openApiSet(content, kind, media)

	node := openApiObject()
	openApiSet(node, "description", openApiString(sandbox, "Up to "+routeDocSize(sandbox, body.MaxBytes)+"."))
	if body.Required {
		openApiSet(node, "required", openApiBool(true))
	}
	openApiSet(node, "content", content)
	return node
}

// openApiSchema is a declared json-schema node as OpenAPI 3.0.3 reads one: the
// same subset, with an enum and a const typed by the node's type, a const
// spelled as a one-value enum, a numeric exclusive bound as the boolean one
// over its minimum or maximum, and a null type as nullable.
func openApiSchema(sandbox *api.Sandbox, schema *routeconf.Schema) *routeDocJson {
	node := openApiObject()
	number := func(value float64) *routeDocJson {
		return openApiLiteral(sandbox.Deps.StringsDeps.FormatFloat(value, 'f', -1, 64))
	}
	count := func(value int) *routeDocJson {
		return openApiLiteral(sandbox.Deps.StringsDeps.FormatInt(int64(value), 10))
	}

	switch {
	case schema.Type == "null":
		openApiSet(node, "nullable", openApiBool(true))
	case schema.Type != "":
		openApiSet(node, "type", openApiString(sandbox, schema.Type))
	}
	if schema.Nullable && schema.Type != "null" {
		openApiSet(node, "nullable", openApiBool(true))
	}
	if schema.Format != "" {
		openApiSet(node, "format", openApiString(sandbox, schema.Format))
	}
	if schema.Pattern != "" {
		openApiSet(node, "pattern", openApiString(sandbox, schema.Pattern))
	}
	switch {
	case schema.HasConst:
		openApiSet(node, "enum", openApiTypedList(sandbox, schema.Type, []string{schema.Const}))
	case len(schema.Enum) > 0:
		openApiSet(node, "enum", openApiTypedList(sandbox, schema.Type, schema.Enum))
	}

	switch {
	case schema.HasExclusiveMinimum && (!schema.HasMinimum || schema.ExclusiveMinimum >= schema.Minimum):
		openApiSet(node, "minimum", number(schema.ExclusiveMinimum))
		openApiSet(node, "exclusiveMinimum", openApiBool(true))
	case schema.HasMinimum:
		openApiSet(node, "minimum", number(schema.Minimum))
	}
	switch {
	case schema.HasExclusiveMaximum && (!schema.HasMaximum || schema.ExclusiveMaximum <= schema.Maximum):
		openApiSet(node, "maximum", number(schema.ExclusiveMaximum))
		openApiSet(node, "exclusiveMaximum", openApiBool(true))
	case schema.HasMaximum:
		openApiSet(node, "maximum", number(schema.Maximum))
	}

	if schema.HasMinLength {
		openApiSet(node, "minLength", count(schema.MinLength))
	}
	if schema.HasMaxLength {
		openApiSet(node, "maxLength", count(schema.MaxLength))
	}
	if schema.Items != nil {
		openApiSet(node, "items", openApiSchema(sandbox, schema.Items))
	}
	if schema.HasMinItems {
		openApiSet(node, "minItems", count(schema.MinItems))
	}
	if schema.HasMaxItems {
		openApiSet(node, "maxItems", count(schema.MaxItems))
	}
	if schema.UniqueItems {
		openApiSet(node, "uniqueItems", openApiBool(true))
	}
	if len(schema.Properties) > 0 {
		properties := openApiObject()
		for _, property := range schema.Properties {
			openApiSet(properties, property.Name, openApiSchema(sandbox, property.Schema))
		}
		openApiSet(node, "properties", properties)
	}
	if len(schema.Required) > 0 {
		openApiSet(node, "required", openApiStrings(sandbox, schema.Required))
	}
	if schema.HasAdditionalProperties {
		openApiSet(node, "additionalProperties", openApiBool(schema.AdditionalProperties))
	}
	return node
}

// openApiResponses is every status the route page lists, in the same words:
// a 200 in the route's response-type, and every failure as the Error json.
func openApiResponses(sandbox *api.Sandbox, conf *routeconf.RouteConf, inherited []routeDocInherited) *routeDocJson {
	responses := openApiObject()
	for _, status := range routeDocStatuses(sandbox, conf, inherited) {
		code := sandbox.Deps.StringsDeps.Trim(status.Code, "`")
		response := openApiObject()
		openApiSet(response, "description", openApiString(sandbox, status.Meaning))

		content := openApiObject()
		if code == "200" {
			if conf.ResponseType != "" {
				openApiSet(content, conf.ResponseType, openApiObject())
			}
		} else {
			reference := openApiObject()
			openApiSet(reference, "$ref", openApiString(sandbox, "#/components/schemas/"+openApiErrorSchema))
			media := openApiObject()
			openApiSet(media, "schema", reference)
			openApiSet(content, openApiJsonMedia, media)
		}
		if len(content.Keys) > 0 {
			openApiSet(response, "content", content)
		}
		openApiSet(responses, code, response)
	}
	return responses
}

// openApiTyped is a declared value as the json literal its type reads as: a
// number or a boolean bare, anything else — or a value that does not parse as
// its type — quoted.
func openApiTyped(sandbox *api.Sandbox, kind string, value string) *routeDocJson {
	switch sandbox.Deps.StringsDeps.TrimSuffix(kind, "-array") {
	case "integer", "number":
		parsed, err := sandbox.Deps.StringsDeps.ParseFloat(value, 64)
		if err == nil && parsed-parsed == 0 {
			return openApiLiteral(sandbox.Deps.StringsDeps.FormatFloat(parsed, 'f', -1, 64))
		}
	case "boolean":
		if value == "true" || value == "false" {
			return openApiLiteral(value)
		}
	}
	return openApiString(sandbox, value)
}

// openApiTypedList is openApiTyped over a list of values.
func openApiTypedList(sandbox *api.Sandbox, kind string, values []string) *routeDocJson {
	list := openApiList()
	for _, value := range values {
		openApiAdd(list, openApiTyped(sandbox, kind, value))
	}
	return list
}

// openApiText joins the parts of a description that say something, one
// paragraph each.
func openApiText(parts ...string) string {
	text := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		if text != "" {
			text += "\n\n"
		}
		text += part
	}
	return text
}

// openApiCapital is a sentence with its first letter upper-cased.
func openApiCapital(sandbox *api.Sandbox, text string) string {
	if text == "" {
		return text
	}
	return sandbox.Deps.StringsDeps.ToUpper(text[:1]) + text[1:]
}

// The document is built as a routeDocJson tree, the node docs/Routes writes
// its json samples with: an object keeps its keys in the order they are set,
// so the same declarations always print the same bytes.

func openApiObject() *routeDocJson {
	return &routeDocJson{Object: true}
}

func openApiList() *routeDocJson {
	return &routeDocJson{List: true}
}

func openApiLiteral(text string) *routeDocJson {
	return &routeDocJson{Literal: text}
}

func openApiString(sandbox *api.Sandbox, text string) *routeDocJson {
	return openApiLiteral(routeDocJsonQuote(sandbox, text))
}

func openApiBool(value bool) *routeDocJson {
	if value {
		return openApiLiteral("true")
	}
	return openApiLiteral("false")
}

func openApiStrings(sandbox *api.Sandbox, values []string) *routeDocJson {
	list := openApiList()
	for _, value := range values {
		openApiAdd(list, openApiString(sandbox, value))
	}
	return list
}

// openApiSet sets key on an object, replacing the child it held.
func openApiSet(node *routeDocJson, key string, child *routeDocJson) {
	for index, existing := range node.Keys {
		if existing == key {
			node.Children[index] = child
			return
		}
	}
	node.Keys = append(node.Keys, key)
	node.Children = append(node.Children, child)
}

// openApiGet is the child an object holds under key, nil when it holds none.
func openApiGet(node *routeDocJson, key string) *routeDocJson {
	for index, existing := range node.Keys {
		if existing == key {
			return node.Children[index]
		}
	}
	return nil
}

// openApiAdd appends a child to a list.
func openApiAdd(list *routeDocJson, child *routeDocJson) {
	list.Children = append(list.Children, child)
}
