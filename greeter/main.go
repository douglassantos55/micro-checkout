package main

import "net/http"

func main() {
	service := new(greeter)
	http.Handle("/greet", MakeHttpHandler(service))
	http.ListenAndServe(":5353", nil)
}
