package conf

type Config struct {
	Server Server
	Log    Logger
	Vars   map[string]*Vars
}

type Server struct {
	Listen string
}

type Logger struct {
	Path  string
	Level int32
}

type Vars struct {
	Vars map[string]*Var
}

type Var struct {
	Value    string
	Readonly bool
	Patterns []string
	Expand   bool
}
