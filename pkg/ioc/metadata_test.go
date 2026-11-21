// Copyright © 2026-present the longxi-go/gosr authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ioc_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/ioc"
)

func validDefinition() ioc.BeanDefinition {
	return ioc.BeanDefinition{
		ID:        "example/repository",
		Name:      "repository",
		Contracts: []ioc.Contract{"example.Repository"},
		Scope:     ioc.Singleton,
	}
}

func TestDefinitionValidation(t *testing.T) {
	cases := map[string]func(*ioc.BeanDefinition){
		"empty id": func(d *ioc.BeanDefinition) {
			d.ID = ""
		},
		"blank name": func(d *ioc.BeanDefinition) {
			d.Name = " "
		},
		"no contracts": func(d *ioc.BeanDefinition) {
			d.Contracts = nil
		},
		"blank contract": func(d *ioc.BeanDefinition) {
			d.Contracts = []ioc.Contract{" "}
		},
		"duplicate contract": func(d *ioc.BeanDefinition) {
			d.Contracts = append(d.Contracts, d.Contracts[0])
		},
		"self alias": func(d *ioc.BeanDefinition) {
			d.Aliases = []string{d.Name}
		},
		"duplicate alias": func(d *ioc.BeanDefinition) {
			d.Aliases = []string{"a", "a"}
		},
		"blank qualifier": func(d *ioc.BeanDefinition) {
			d.Qualifiers = []string{""}
		},
		"duplicate qualifier": func(d *ioc.BeanDefinition) {
			d.Qualifiers = []string{"q", "q"}
		},
		"no scope": func(d *ioc.BeanDefinition) {
			d.Scope = ""
		},
		"negative origin": func(d *ioc.BeanDefinition) {
			d.Source.Line = -1
		},
		"missing dependency selector": func(d *ioc.BeanDefinition) {
			d.Dependencies = []ioc.Dependency{
				{
					Access:      ioc.Direct,
					Cardinality: ioc.One,
				},
			}
		},
		"implicit dependency access": func(d *ioc.BeanDefinition) {
			d.Dependencies = []ioc.Dependency{
				{
					Selector: ioc.Selector{
						Name: "db",
					},
					Cardinality: ioc.One,
				},
			}
		},
		"implicit cardinality": func(d *ioc.BeanDefinition) {
			d.Dependencies = []ioc.Dependency{
				{
					Selector: ioc.Selector{
						Name: "db",
					},
					Access: ioc.Direct,
				},
			}
		},
		"invalid dependency origin": func(d *ioc.BeanDefinition) {
			d.Dependencies = []ioc.Dependency{
				{
					Selector: ioc.Selector{
						Name: "db",
					},
					Access:      ioc.Direct,
					Cardinality: ioc.One,
					Source: ioc.Source{
						Line: -1,
					},
				},
			}
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			d := validDefinition()
			change(&d)
			err := d.Validate()
			if !errors.Is(err, &ioc.Error{
				Code: ioc.InvalidDefinition,
			}) {
				t.Fatalf("invalid declaration accepted: %v", err)
			}
		})
	}
	d := validDefinition()
	d.Scope = "tenant"
	d.Lazy = true
	d.Order = -7
	d.Dependencies = []ioc.Dependency{
		{
			Selector: ioc.Selector{
				Name: "db",
			},
			Access:      ioc.Deferred,
			Cardinality: ioc.ZeroOrOne,
		},
	}
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSelectorValidation(t *testing.T) {
	for _, s := range []ioc.Selector{
		{},
		{
			Contract: " ",
		},
		{
			Name: " ",
		},
		{
			Name:       "db",
			Qualifiers: []string{"a", "a"},
		},
		{
			Contract:   "DB",
			Qualifiers: []string{""},
		},
	} {
		if err := s.Validate(); !errors.Is(err, &ioc.Error{
			Code: ioc.InvalidSelector,
		}) {
			t.Fatalf("accepted %#v: %v", s, err)
		}
	}
	if err := (ioc.Selector{
		Contract: "map[string]example.DB",
	}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDefinitionCloneDetachesEverySlice(t *testing.T) {
	original := validDefinition()
	original.Aliases = []string{"main"}
	original.Qualifiers = []string{"prod"}
	original.Dependencies = []ioc.Dependency{
		{
			Selector: ioc.Selector{
				Name:       "db",
				Qualifiers: []string{"durable"},
			},
			Access:      ioc.Direct,
			Cardinality: ioc.One,
		},
	}
	clone := original.Clone()
	clone.Contracts[0] = "changed"
	clone.Aliases[0] = "changed"
	clone.Qualifiers[0] = "changed"
	clone.Dependencies[0].Selector.Qualifiers[0] = "changed"
	clone.Dependencies[0].Selector.Name = "changed"
	if original.Contracts[0] != "example.Repository" || original.Aliases[0] != "main" || original.Qualifiers[0] != "prod" || original.Dependencies[0].Selector.Qualifiers[0] != "durable" || original.Dependencies[0].Selector.Name != "db" {
		t.Fatal("clone changed source metadata")
	}
	zero := (ioc.BeanDefinition{}).Clone()
	if zero.Contracts != nil || zero.Dependencies != nil {
		t.Fatal("nil shape changed")
	}
}

func TestMatchesUsesConjunctionAndAliasIdentity(t *testing.T) {
	d := validDefinition()
	d.Aliases = []string{"main"}
	d.Qualifiers = []string{"west", "prod"}
	for _, s := range []ioc.Selector{
		{
			Name: "main",
		},
		{
			Contract:   "example.Repository",
			Name:       "repository",
			Qualifiers: []string{"prod", "west"},
		},
	} {
		if !d.Matches(s) {
			t.Fatalf("matching selector rejected: %#v", s)
		}
	}
	for _, s := range []ioc.Selector{
		{
			Name: d.ID,
		},
		{
			Name: "Main",
		},
		{
			Name:     "main",
			Contract: "other",
		},
		{
			Qualifiers: []string{"west", "absent"},
		},
	} {
		if d.Matches(s) {
			t.Fatalf("incorrect match: %#v", s)
		}
	}
}

func TestErrorsKeepCodeCauseAndSource(t *testing.T) {
	cause := errors.New("allocation failed")
	err := &ioc.Error{
		Code: ioc.CreationFailed,
		Bean: "repository",
		Source: ioc.Source{
			File: "repo.go",
			Line: 7,
		},
		Path:    []string{"service", "repository"},
		Message: "construction",
		Cause:   cause,
	}
	if !errors.Is(err, cause) || !errors.Is(err, &ioc.Error{
		Code: ioc.CreationFailed,
	}) || errors.Is(err, &ioc.Error{
		Code: ioc.MissingBean,
	}) || errors.Is(err, &ioc.Error{}) {
		t.Fatal("error identity lost")
	}
	if !strings.Contains(err.Error(), "repo.go:7") || !strings.Contains(err.Error(), "service -> repository") || !strings.Contains(err.Error(), "allocation failed") {
		t.Fatal("missing diagnostics")
	}
	var typed *ioc.Error
	if !errors.As(err, &typed) || typed.Source.Line != 7 {
		t.Fatal("source lost")
	}
	var empty *ioc.Error
	if empty.Unwrap() != nil || empty.Is(&ioc.Error{
		Code: ioc.CreationFailed,
	}) {
		t.Fatal("nil error unsafe")
	}
}
