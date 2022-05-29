package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(ImdbTtitle("tt0076759"))
	if err := http.ListenAndServe(fetchPort(), nil); err != nil {
		panic(err)
	}
}
