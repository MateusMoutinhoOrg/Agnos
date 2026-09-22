

## Refs:
-  [samples](New-Server-Mode/samples) 


### New Handler Assignature:
```go 

func RouteHandler(sandbox *api.Sandbox, route *api.Route, *response serverdeps.Response) error  {
	
	return
}

```
### Priority Mechanic:
the prioriy system workss rom button to top, if a  route has priority 0 and other has priority 5, it will run first priority 0 if it not made a write on body or not returned a err, it will run priority 5 until it makes a write on body or return a err

### Identifiers Logic
the indentifiers must work on params,routes,and headers, they must works as a pattern matching system, the root handler will list all routes, and then for each route, verifiy all **starts-with-identifiers** and **equal-identifiers** that match the current request, and if all that identifiers match, it will add the route to a list of routes to run.

### Lopp Logic
the root handler must loop thought all routes, and list all routes tha all identifiers of that routes match, and then,it will run the routes sorted by priority, in case the first,second,third, route seted a StatusCode,it will return that route, in case the last route not seted a StatusCode,it will run the Handler404 , in case of a first.. route returns a error, it will return that error.

## Handler404 
its a special file that wil lbe created on init of server-init command,these will contain a function with a default 404 reesponse, but in these way, allows the user to creates his own response:

### sandbox/interna/server/handler404.go:
```go 

func RouteHandler(sandbox *api.Sandbox, route *api.Route, *response serverdeps.Response) error  {
	
	return
}

```
