objetivo: crie a classe smartIO 


### smartio.go 
#### dest: 
utils/smartio/smartio.go
#### required Assignatures: 
~~~go 
package smartio

type SmartIO struct {
    // props 

    Ignore  IgnorableConf // in case not exist, starts empty
    Replacers PathReplacerConf  // in case of not exist starts empty
    Transactions map[string][]byte // used to store data that are "in memory" ready to be written



    //worflow:
    // aply the replacer in path
    // verifiy if the item is inside the ignore, if it is,returns a error 
	//returns the byte seq of file
    ReadFile func(path string) ([]byte, error)


    //workflow:
    // aply the replacer in path
    // verifiy if the item is inside the ignore, if it is,returns a error
    // verify if the dest not exist, if exist, returns a error 
    // add the byte seq to a in-memory map
	WriteFile func(path string, content []byte) error

    //equals to WriteFile, but it dont check if the file not exist, before writing, it will overwrite the file if exist
	WriteFileOverwrite func(path string, content []byte) error

    //writes the Transaction in hardware usin the deps io
    Persist func()error 

    
	IsDir func(path string) bool

	IsFile func(path string) bool

	Exist func(path string) bool


	CreateDir func(path string)

	ListDirs func(path string) []string

	ListFiles func(path string) []string


	ListAll func(path string) []string


	ListDirsRecursively func(path string) []string


	ListFilesRecursively func(path string) []string


	ListAllRecursively func(path string) []string
}


func NewSmartIO(path string) *SmartIO
~~~

### Important
for all elements, it needs to first aply the replacers,them check ignore.