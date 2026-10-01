// Adapted from Traefik v2.11.57 (pkg/rules/parser.go), https://github.com/traefik/traefik.
// Copyright (c) 2016-2020 Containous SAS; Copyright (c) 2020-2026 Traefik Labs.
// Licensed under the Apache License, Version 2.0; see NOTICE in this directory.
// Changes: moved into this package, logging via logrus, the Host matcher
// reads the request's own Host (the forward-auth server has no request
// decorator), CNAME flattening removed.

package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/vulcand/predicate"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	and = "and"
	or  = "or"
)

// TreeBuilder defines the type for a Tree builder.
type TreeBuilder func() *Tree

// Tree represents the rules' tree structure.
type Tree struct {
	Matcher   string
	Not       bool
	Value     []string
	RuleLeft  *Tree
	RuleRight *Tree
}

// NewParser constructs a parser for the given matchers.
func NewParser(matchers []string) (predicate.Parser, error) {
	parserFuncs := make(map[string]any)

	for _, matcherName := range matchers {
		fn := func(value ...string) TreeBuilder {
			return func() *Tree {
				return &Tree{
					Matcher: matcherName,
					Value:   value,
				}
			}
		}
		parserFuncs[matcherName] = fn
		parserFuncs[strings.ToLower(matcherName)] = fn
		parserFuncs[strings.ToUpper(matcherName)] = fn
		parserFuncs[cases.Title(language.Und).String(strings.ToLower(matcherName))] = fn
	}

	return predicate.NewParser(predicate.Def{
		Operators: predicate.Operators{
			AND: andFunc,
			OR:  orFunc,
			NOT: notFunc,
		},
		Functions: parserFuncs,
	})
}

func andFunc(left, right TreeBuilder) TreeBuilder {
	return func() *Tree {
		return &Tree{
			Matcher:   and,
			RuleLeft:  left(),
			RuleRight: right(),
		}
	}
}

func orFunc(left, right TreeBuilder) TreeBuilder {
	return func() *Tree {
		return &Tree{
			Matcher:   or,
			RuleLeft:  left(),
			RuleRight: right(),
		}
	}
}

func invert(t *Tree) *Tree {
	switch t.Matcher {
	case or:
		t.Matcher = and
		t.RuleLeft = invert(t.RuleLeft)
		t.RuleRight = invert(t.RuleRight)
	case and:
		t.Matcher = or
		t.RuleLeft = invert(t.RuleLeft)
		t.RuleRight = invert(t.RuleRight)
	default:
		t.Not = !t.Not
	}

	return t
}

func notFunc(elem TreeBuilder) TreeBuilder {
	return func() *Tree {
		return invert(elem())
	}
}

// ParseMatchers returns the subset of matchers in the Tree matching the given matchers.
func (tree *Tree) ParseMatchers(matchers []string) []string {
	switch tree.Matcher {
	case and, or:
		return append(tree.RuleLeft.ParseMatchers(matchers), tree.RuleRight.ParseMatchers(matchers)...)
	default:
		if slices.Contains(matchers, tree.Matcher) {
			return lower(tree.Value)
		}

		return nil
	}
}

// CheckRule validates the given rule.
func CheckRule(rule *Tree) error {
	if len(rule.Value) == 0 {
		return fmt.Errorf("no args for matcher %s", rule.Matcher)
	}

	for _, v := range rule.Value {
		if len(v) == 0 {
			return fmt.Errorf("empty args for matcher %s, %v", rule.Matcher, rule.Value)
		}
	}

	return nil
}

func lower(slice []string) []string {
	var lowerStrings []string
	for _, value := range slice {
		lowerStrings = append(lowerStrings, strings.ToLower(value))
	}

	return lowerStrings
}
