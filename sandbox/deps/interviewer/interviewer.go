package interviewer

type AlterativeOption struct {
	Id  string
	Msg bool
}

type Sandbox struct {
	IntQuestion   func(question string) (int, error)
	StrQuestion   func(question string) (string, error)
	FloatQuestion func(question string) (float64, error)

	//add Yes or No for answer
	BoolQuestion func(question string) (bool, error)
	//returns the id of the answer
	SingleAlternativeQuestion func(question string, alternatives []AlterativeOption) (string, error)
	//returns the ids of the answers
	MultipleAlternativeQuestion func(question string, alternatives []AlterativeOption) ([]string, error)
}
