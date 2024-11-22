package main

import (
	"fmt"
	"net/http"

	"github.inclee.matchtree/mux"
)

func main() {
	exprStr := "Path(`/api/v1`) && Method(`GET`) && Header(`A`,`C`)"
	router := mux.Router{}
	router.AddRule(exprStr, 1)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1", nil)
	req.Header.Add("A", "C")
	rule := router.Match(*req)
	if rule != nil {
		fmt.Println("matched: ", rule)
	}
}
