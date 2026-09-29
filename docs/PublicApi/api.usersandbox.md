# `sandbox/api/usersandbox.go`

## `UserSandbox`

UserSandbox is the part of the Sandbox the project declares itself: api.Sandbox embeds it, so every field typed here is a field of the Sandbox, reached as sandbox.<Field> by every function handed it, like any contract the build lists: type UserSandbox struct { Greeting string } Fill it from a package of your own under sandbox/constructors/ — sandbox/new.go calls its Constructor(sandbox) like any generated one. A field is written in the builtin types or in a type of this package, like every other of sandbox/api/. Written once by `agnos start` and then yours: no build rewrites this file.

[every contract](doc.md)
