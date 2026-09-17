# `deps.Hashdeps`

`sandbox/deps/hashdeps`

## `Sandbox`

Sandbox is the hashing library injected whole as the Deps.Hashdeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Sha256Hex` | `func(content []byte) string` | Sha256Hex returns the SHA-256 digest of content, lower-case hexadecimal. It is what every recorded example tree is compared by, so the encoding is part of the golden and may not change. |

[every contract](doc.md)
