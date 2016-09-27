package exec

type Exec interface {
	Exec(cntx ExecutionContext) error
	String() string
}
