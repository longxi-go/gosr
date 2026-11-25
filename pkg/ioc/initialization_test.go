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
	"testing"

	"github.com/longxi-go/gosr/pkg/ioc"
)

func TestInitializationReferenceValidation(t *testing.T) {
	cases := []struct {
		name  string
		names []string
		valid bool
	}{
		{
			name:  "absent",
			valid: true,
		},
		{
			name:  "unresolved names",
			names: []string{"migration", "warmup"},
			valid: true,
		},
		{
			name:  "blank",
			names: []string{" "},
		},
		{
			name:  "duplicate",
			names: []string{"migration", "migration"},
		},
		{
			name:  "self name",
			names: []string{"repository"},
		},
		{
			name:  "self alias",
			names: []string{"main"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := validDefinition()
			d.Aliases = []string{"main"}
			d.DependsOn = tc.names
			err := d.Validate()
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
			} else if !errors.Is(err, &ioc.Error{
				Code: ioc.InvalidDefinition,
			}) {
				t.Fatalf("want InvalidDefinition, got %v", err)
			}
		})
	}
}

func TestInitializationReferenceCloneOwnership(t *testing.T) {
	d := validDefinition()
	d.DependsOn = []string{"migration"}
	copy := d.Clone()
	copy.DependsOn[0] = "changed"
	if d.DependsOn[0] != "migration" || (ioc.BeanDefinition{}).Clone().DependsOn != nil {
		t.Fatal("initialization metadata ownership or nil shape changed")
	}
	for _, registry := range []ioc.DefinitionRegistry{&sliceRegistry{}, &mapRegistry{}} {
		if err := registry.Register(d); err != nil {
			t.Fatal(err)
		}
		snapshot := registry.Definitions()
		snapshot[0].DependsOn[0] = "snapshot changed"
		d.DependsOn[0] = "source changed"
		if registry.Definitions()[0].DependsOn[0] != "migration" {
			t.Fatal("registry shares initialization references")
		}
		d.DependsOn[0] = "migration"
	}
}
