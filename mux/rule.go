package mux

import "net/http"

type Rule struct {
	matcher  *MatchTree
	priority int
}

func newRule(rule string, priority int) (*Rule, error) {
	fns := []string{}
	for fn := range funcs {
		fns = append(fns, fn)
	}
	parse, err := MatchTreeParse(fns)
	if err != nil {
		panic(err)
	}
	expr, err := parse.Parse(rule)
	if err != nil {
		return nil, err
	}
	mtree, err := newMatchTree(expr.(TreeBuilder)())
	if err != nil {
		return nil, err
	}
	return &Rule{
		priority: priority,
		matcher:  mtree,
	}, nil
}

func (r *Rule) Match(req http.Request) bool {
	return r.matcher.Match(req)
}
