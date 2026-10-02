package routeprops

// Project is the part of RouteProps this project declares itself, embedded so
// each field is read as props.<Field>. A middleware sets what it learned — the
// user it authenticated — and every route after it reads it, typed:
//
//	type Project struct {
//		// User is the caller the auth middleware authenticated, nil when none.
//		User *maindatabase.UserItem
//	}
//
// It lives under sandbox/internal, not in sandbox/api, so a field may name any
// type of the project — a record of one of its databases, a type of one of its
// own packages — as long as that package imports no route.
//
// Written once by `{{.GeneratorName}} build` and the project's from then on.
type Project struct {
}
