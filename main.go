package main

import (
	"net/http"
)

func main() {
	Ly3("Coca cola tu")
	if err := http.ListenAndServe(fetchPort(), nil); err != nil {
		panic(err)
	}
}
