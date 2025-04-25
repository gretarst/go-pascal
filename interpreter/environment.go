package interpreter

type Environment struct {
	store map[string]int
}

func NewEnviroment() *Environment {
	return &Environment{store: make(map[string]int)}
}

func (e *Environment) Set(name string, value int) {
	e.store[name] = value
}

func (e *Environment) Get(name string) (int, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Exists(name string) bool {
	_, ok := e.store[name]
	return ok
}
