# `keepdeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	NewDatabase func(props Props) KeepDatabase
}

type KeepDatabase struct {
	Props     Props
	GetSchema func(name string) (SchemaInstance, bool)
}

type SchemaInstance struct {
	Items     []Item
	Prefix    string
	NewItem   func(fields map[string]any) (SchemaItem, *Error)
	FindByKey func(key string, keyValue any) (SchemaItem, bool)
	ListAll   func() ([]SchemaItem, *Error)
	List      func(position int, chunk int) ([]SchemaItem, *Error)
}

type SchemaItem struct {
	Items             []Item
	Prefix            string
	Id                int64
	Get               func(fieldName string) (any, *Error)
	Update            func(fieldName string, value any) *Error
	Remove            func() *Error
	CheckKeysPresence func(keys []string) bool
	ListAll           func(fieldName string) []SchemaItem
	NewSubItem        func(fieldName string, fields map[string]any) (SchemaItem, *Error)
	String            func() string
}
```

## Description

The sandbox's **copy** of the embedded Keep schema-database library's public api, declared in `sandbox/contracts/deps/keepdeps/` and injected whole as the [`deps.Deps.KeepLib`](/docs/LibUsage/References/PublicApi/deps.Deps.md) field. The sandbox may not import Keep — that would be a third-party import — so it restates the shape it needs, struct for struct; the adapter, which lives outside the sandbox, initializes the real library and fills this copy. See [StructContracts.md](/docs/Development/References/StructContracts.md).

Keep is a schema database over a single-key storage backend: a [`Props`](#supporting-types) description declares collections (`Schema`) and their fields (`Item`), `NewDatabase` binds that description, `GetSchema` hands back one collection, and a collection creates, finds, lists and removes records. A field typed `Key` is unique and indexed — that is what `FindByKey` looks through, and inserting a duplicate fails with a `KeyConflict` error. A field typed `Database` is a nested collection, reached through a record's `NewSubItem` and `ListAll`.

Where Verb's fields return plain values, Keep's return further api structs, so this copy restates the whole tree and the adapter's `KeepLibFactory` wraps each such field in a closure that converts what it returns — nothing of the embedded library ever reaches the sandbox. See [HandleDependencies.md](/docs/Development/Protocols/HandleDependencies.md#injecting-a-whole-library).

The adapter also chooses the backend Keep itself runs on: `standard` wires Keep's filesystem adapter. See [Adapters.md](/docs/LibUsage/References/Adapters.md).

## Fields

### `Lib`

| Field | Description |
| :--- | :--- |
| `NewDatabase func(props Props) KeepDatabase` | Creates a database from a `Props` description. |

### `KeepDatabase`

| Field | Description |
| :--- | :--- |
| `Props Props` | The description the database was created from. |
| `GetSchema func(name string) (SchemaInstance, bool)` | Returns the collection with the given name; `false` when the database declares none. |

### `SchemaInstance`

| Field | Description |
| :--- | :--- |
| `Items []Item` | The fields each record of the collection can hold. |
| `Prefix string` | The collection's key prefix. |
| `NewItem func(fields map[string]any) (SchemaItem, *Error)` | Inserts a record, validating the fields against the schema. |
| `FindByKey func(key string, keyValue any) (SchemaItem, bool)` | Looks a record up through a unique `Key` field; `false` when the field is not indexed or nothing matches. |
| `ListAll func() ([]SchemaItem, *Error)` | Returns every record of the collection. |
| `List func(position int, chunk int) ([]SchemaItem, *Error)` | Returns up to `chunk` records starting at `position` (1-based). |

### `SchemaItem`

| Field | Description |
| :--- | :--- |
| `Items []Item` | The schema fields this record's collection declares. |
| `Prefix string` | The collection's key prefix. |
| `Id int64` | The record's permanent, never-reused identifier. |
| `Get func(fieldName string) (any, *Error)` | Returns the typed value stored for a field. |
| `Update func(fieldName string, value any) *Error` | Writes a new value, re-indexing it when the field is a unique key. |
| `Remove func() *Error` | Deletes the record and everything nested under it. |
| `CheckKeysPresence func(keys []string) bool` | Reports whether every named field has a stored value. |
| `ListAll func(fieldName string) []SchemaItem` | Returns every record of a nested (`Database`) field. |
| `NewSubItem func(fieldName string, fields map[string]any) (SchemaItem, *Error)` | Inserts a record into a nested (`Database`) field. |
| `String func() string` | Renders the record's plain fields. |

### Supporting types

| Type | Description |
| :--- | :--- |
| `Item{Name, Type, Required, Itens}` | One field of a schema. `Type` is `Key`, `Int`, or `Database`; `Itens` holds the nested fields of a `Database` field. |
| `Schema{Name, Itens}` | One collection of records and its fields. |
| `Props{Path, Schemas}` | The declarative description a database is created from. `Path` is the prefix every key is written under. |
| `Error{Type, Key, KeyValue, Message}` | One failure. `Type` is `KeyConflict`, `NotFound`, `MissingField`, `InvalidField`, or `Internal`. A nil `*Error` means success. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoskeep "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
)

func main() {
	// The adapter initializes the embedded Keep library — here over Keep's
	// own filesystem backend — and hands it back as one field of the deps
	// contract.
	d := agnosadapter.New("trackerdata")

	db := d.KeepLib.NewDatabase(agnoskeep.Props{
		Path: "app/",
		Schemas: []agnoskeep.Schema{{
			Name: "user",
			Itens: []agnoskeep.Item{
				{Name: "email", Type: agnoskeep.Key, Required: true},
				{Name: "age", Type: agnoskeep.Int, Required: true},
				{Name: "sessions", Type: agnoskeep.Database, Itens: []agnoskeep.Item{
					{Name: "token", Type: agnoskeep.Key, Required: true},
				}},
			},
		}},
	})

	users, ok := db.GetSchema("user")
	if !ok {
		return
	}

	user, err := users.NewItem(map[string]any{"email": "a@b.c", "age": 27})
	if err != nil {
		fmt.Println("insert failed:", err.Message)
		return
	}

	// A duplicate value for a Key field is a KeyConflict, not a panic.
	if _, err := users.NewItem(map[string]any{"email": "a@b.c", "age": 30}); err != nil {
		fmt.Println(err.Type == agnoskeep.KeyConflict, err.Message) // true ...
	}

	user.NewSubItem("sessions", map[string]any{"token": "t1"})
	fmt.Println(len(user.ListAll("sessions"))) // 1

	found, _ := users.FindByKey("email", "a@b.c")
	fmt.Println(found.String()) // {id: 1, email: a@b.c, age: 27}
}
```
