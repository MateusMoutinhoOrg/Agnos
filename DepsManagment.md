## dep-install,dep-remove ,dep-list mechanic
create a mechanic for auto installed embed deps in the project


## Assets structs:
all the deps , must be retrived from assets/deplist/ 
each dir inside assets/deplist its a dep. 

## dep-install

command: agnos dep-install <dep>

this command will install the embeddeps inside the project

command:
~~~bash
agnos dep-install "embed" 
~~~

action:

~~~go
func DepInstallInternal(deps *deps.Deps, io *smartio.SmartIO, path string,dep string) error 

func DepInstall(deps *deps.Deps, path string,dep string) error
~~~


## dep-remove 

this command will remove the embeddeps inside the project
command:
~~~bash 
agnos dep-remove "embed" 
~~~

action:

~~~go
func DepRemoveInternal(deps *deps.Deps, io *smartio.SmartIO, path string,dep string) error 

func DepRemove(deps *deps.Deps, path string,dep string) error
~~~


## dep-list

will list all available deps and their description

command:
~~~bash
agnos dep-list 
~~~

action:

~~~go 
func DepListInternal(deps *deps.Deps, io *smartio.SmartIO, path string) ([]string, error) 

func DepList(deps *deps.Deps, path string) ([]string, error)
~~~ 
