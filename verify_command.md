## verify Action and command

add the command **agnos verify** 
ideia:
the ideia of agnos project,  its to provide a full factory with harness generation for go language , so it must garantee schematization. conssistency.

Command workflow:

1. Sandbox Verification:

- sandbox/* can only import files inside sandbox
- sandbox dir  can only contains api,binds,deps,internal and new.go 
- sandbox/api/*  cannot import anything.
- sandbox/deps/* cannot import anything
- sandbox/binds/ each file must have a equal name of a equivalent file in sandbox/api 
- sandbox/binds/ must only contains functions 



2. Adapters verififcation
- adapters/ can only have the dir available and libs

