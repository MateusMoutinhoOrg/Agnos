
### Ignorable Parser 

Create a parser for the ignore.yaml config.

### ignore.yaml
its a document, similar to project.yaml and themes.yaml, that contains only a list of strings

#### sample:
~~~yaml
- sandbox/contacts/*
- scripts/*
~~~

### Parser
the parser must have the following assignature:

~~~go
type IgnorableItens interface {


    AddPath(path string)
    // returns true if eleent its on ignorable list
	IsIgnorable(path string)bool 
}

func NewIgnorableItens(content []string) (IgnorableItens, error){
    
}
~~~


## CONTEX:
- Read : sandbox/parsables/themes.go to understand the pattern

