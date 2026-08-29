### Objective
Rename the contracts dir and objects with names more easly to understand for the user

## Dirs:
| Dir  | New Dir | Comment |
| ---- | -------- | ------- |
| `contracts/api` | `contracts/sdk` | Sdk is a set of contracts that are used to communicate between processes |
| `contracts/deps` | `contracts/requirements` | Depedencies of the project |


### Classes:

| Class | Class Location |  New Class | New Location | Comment |
| ----- | ---------      |  --------  | ----------  | ------- |
| api.SandBox | contracts/api/api.go | api.Sdk | contracts/sdk/sdk.go | Core contracts are the base contracts that are used to communicate between processes |
| api.Deps | contracts/deps/deps.go| api.Requirements | contracts/requirements/requirements.go | Dependencies of the project