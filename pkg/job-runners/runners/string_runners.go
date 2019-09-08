package runners

// ToStringRunner is responsible for taking a function (with no additional arguments) and providing an interface
// by means of implementation transforms the strings (once primed for execution).
type (
	ToStringJob    func() (string, error)
	ToStringRunner interface {
		AddJob(transformer ToStringJob) []ToStringJob
		AddJobs(transformations []ToStringJob) []ToStringJob
		prime()
		Run() ([]string, error)
	}
)