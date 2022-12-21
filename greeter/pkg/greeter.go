package pkg

import "fmt"

type Greeter interface {
	Greet(name string) string
}

type greeter struct{}

func NewGreeter() Greeter {
	return &greeter{}
}

func (g *greeter) Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome.", name)
}
