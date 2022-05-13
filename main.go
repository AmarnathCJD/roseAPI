package main

import (
	"net/http"
)

func main() {
	if err := http.ListenAndServe(fetchPort(), nil); err != nil {
		panic(err)
	}
}
