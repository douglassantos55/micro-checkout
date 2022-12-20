package main

import "fmt"

type Greeter interface {
	Greet(name string) map[string]string
}

type greeter struct{}

func (g *greeter) Greet(name string) map[string]string {
	return map[string]string{
		"greeting": fmt.Sprintf("Hello, %s! Welcome.", name),
	}
}
