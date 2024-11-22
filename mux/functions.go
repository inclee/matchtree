package mux

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

var funcs = map[string]func(tree *MatchTree, values ...interface{}) error{
	"Path":   expectedNumArgs(path, 1),
	"Method": expectedNumArgs(method, 1),
	"Header": expectedNumArgs(header, 2),
}

func expectedNumArgs(fn func(tree *MatchTree, values ...interface{}) error, n int) func(tree *MatchTree, values ...interface{}) error {
	return func(tree *MatchTree, values ...interface{}) error {
		if len(values) != n {
			funcName := getFuncName(fn)
			return fmt.Errorf("function [%s] expected %d parameter(s), but got %d", funcName, n, len(values))
		}
		return fn(tree, values...)
	}
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func path(tree *MatchTree, values ...interface{}) error {
	expectedPath, ok := values[0].(string)
	if !ok {
		return fmt.Errorf("path matcher expects a string as the first parameter")
	}
	tree.matcher = func(r http.Request) bool {
		return r.URL.Path == expectedPath
	}
	return nil
}

func method(tree *MatchTree, values ...interface{}) error {
	expectedMethod, ok := values[0].(string)
	if !ok {
		return fmt.Errorf("method matcher expects a string as the first parameter")
	}
	tree.matcher = func(r http.Request) bool {
		return strings.EqualFold(r.Method, strings.TrimSpace(expectedMethod))
	}
	return nil
}

func header(tree *MatchTree, values ...interface{}) error {
	headerKey, ok1 := values[0].(string)
	headerValue, ok2 := values[1].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("header matcher expects both parameters to be strings")
	}
	tree.matcher = func(r http.Request) bool {
		return strings.EqualFold(r.Header.Get(headerKey), headerValue)
	}
	return nil
}
