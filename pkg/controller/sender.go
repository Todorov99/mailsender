package controller

import (
	"fmt"
	"net/http"
)

func Send(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Testing")
}
