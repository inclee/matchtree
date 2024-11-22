package mux

import (
	"github.com/vulcand/predicate"
)

type Tree struct {
	fn        string
	predicate string
	value     []interface{}
	left      *Tree
	right     *Tree
}

type TreeBuilder func() *Tree

func MatchTreeParse(funs []string) (predicate.Parser, error) {
	logicOperators := predicate.Operators{
		AND: func(left, right TreeBuilder) TreeBuilder {
			return func() *Tree {
				return &Tree{
					predicate: "and",
					left:      left(),
					right:     right(),
				}
			}

		},
		OR: func(left, right TreeBuilder) TreeBuilder {
			return func() *Tree {
				return &Tree{
					predicate: "or",
					left:      left(),
					right:     right(),
				}
			}
		},
	}
	fns := map[string]interface{}{}
	for _, fn := range funs {
		fns[fn] = func(value ...interface{}) TreeBuilder {
			return func() *Tree {
				return &Tree{
					fn:    fn,
					value: value,
				}
			}
		}
	}
	return predicate.NewParser(predicate.Def{
		Functions: fns,
		Operators: logicOperators,
	})
}
