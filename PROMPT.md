## Doc Refactor


### New Doc Tree
- docs/index/<theme>.md
- docs/docs/<doc-name>/<doc-name>.md
- docs/docs/<doc-name>/props.yaml
- docs/docs/<doc-name>/sub-docs/

### docs/index/<theme>.md
#### automated: yes
#### description:
for each theme listed in themes.yaml , needs to create a file containing a table, listing all the docs that are flaged in props.yaml.

#### model:
```
#  <Theme Name>
<Theme Description>
| Doc | Description |
| --- | --- |
| [Doc Name](/docs/docs/<doc-name>/<doc-name>.md) | <doc description> |
```

### docs/docs/<doc-name>/<doc-name>.md
#### automated: no
#### description:
its the documentation the user created 



### docs/docs/<doc-name>/props.yaml
#### automated: no
#### description:
contains the themes id that documentation is part 
#### sample:
```yaml
- cli-usage
- lib-usage

```

### docs/docs/<doc-name>/sub-docs/ 
#### automated: no
#### description:
docs contained in here should be related to this doc, for example, in public api, only the root public api , will be **docs/docs/Public-Api/Public-Api.md** the other docs of public api, will be inside sub-docs, and indexable for  **docs/docs/Public-Api/Public-Api.md**




### Task:
- Refactor all the documentation to match these patten
- Adapt the code, to generate the all the docs/index 
