package executor

type FileExecutor interface {
	Apply(file string) error
}

