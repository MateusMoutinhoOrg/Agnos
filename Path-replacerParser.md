
### Path Parser 

Create a parser for the paths.yaml config.

### paths.yaml
its a document, similar to project.yaml and themes.yaml, that contains a list of paths to be replaced by the builder.
- 
#### sample:
~~~yaml
docs: Documentation
docs/index/cliUsage/: docs/cli
~~~

### Parser
the parser must have the following assignature:

~~~go
type PathReplacerConf interface {

    // recives a path, iterate over the elements, and returns a formated string apllying the replacment the user wants.
    Format func (path string) string

}

~~~


## CONTEX:
- Read : sandbox/parsables/themes.go to understand the pattern

