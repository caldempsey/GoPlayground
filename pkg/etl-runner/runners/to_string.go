package runners

// ToStringRunner is responsible for taking a function (with no additional arguments) and providing an interface
// by means of implementation transforms the strings (once primed for execution).
type (
	ToString       func() (string, error)
	ToStringRunner interface {
		AddJob(transformer ToString) []ToString
		AddJobs(transformations []ToString) []ToString
		prime()
		Run() ([]string, error)
	}
)