package routeprops

// RouteProps is what one request carries from route to route of its chain: the
// dispatch builds one, empty, per request and hands the same one to every
// InternalPureHandler that runs for it, as its first argument. A middleware
// sets what it learned — the user it authenticated — and every route after it
// reads it, typed:
//
//	type RouteProps struct {
//		// User is the caller the auth middleware authenticated, nil when none.
//		User *maindatabase.UserItem
//	}
//
// It lives under sandbox/internal, not in sandbox/api, so a field may name any
// type of the project — a record of one of its databases, a type of one of its
// own packages — as long as that package imports no route.
//
// Written once by `{{.GeneratorName}} build` and the project's from then on:
// declare here whatever the routes of this project hand each other.
type RouteProps struct {
}
