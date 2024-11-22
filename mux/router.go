package mux

import (
	"net/http"
	"sort"
)

type Router struct {
	rules []*Rule
}

func (r *Router) AddRule(expr string, priority int) error {
	rule, err := newRule(expr, priority)
	if err != nil {
		return err
	}
	r.rules = append(r.rules, rule)
	sort.Slice(r.rules, func(i, j int) bool {
		return r.rules[i].priority < r.rules[j].priority
	})
	return nil
}

func (r *Router) Match(req http.Request) *Rule {
	for _, rule := range r.rules {
		if rule.Match(req) {
			return rule
		}
	}
	return nil
}
