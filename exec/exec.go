package exec

// Exec is the central interface for executable
// processing steps
type Exec interface {
	// Exec executes the step
	// - cntx the data and reporting context
	Exec(cntx Context) error

	// String returns the description for the step
	String(cntx Context) string
}
