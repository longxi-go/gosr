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

// These two deliberately small metadata fixtures prove that the registry
// contract does not require a particular storage engine or a Bean constructor.
type sliceRegistry struct {
	definitions []ioc.BeanDefinition
}

func (r *sliceRegistry) Register(batch ...ioc.BeanDefinition) error {
	seen := map[string]bool{}
	for _, d := range r.definitions {
		seen[d.ID] = true
	}
	for _, d := range batch {
		if err := d.Validate(); err != nil {
			return err
		}
		if seen[d.ID] {
			return &ioc.Error{
				Code: ioc.DuplicateDefinition,
				Bean: d.ID,
			}
		}
		seen[d.ID] = true
	}
	for _, d := range batch {
		r.definitions = append(r.definitions, d.Clone())
	}
	return nil
}

func (r *sliceRegistry) Definitions() []ioc.BeanDefinition {
	out := make([]ioc.BeanDefinition, len(r.definitions))
	for i, d := range r.definitions {
		out[i] = d.Clone()
	}
	return out
}

type mapRegistry struct {
	definitions map[string]ioc.BeanDefinition
}

func (r *mapRegistry) Register(batch ...ioc.BeanDefinition) error {
	next := map[string]ioc.BeanDefinition{}
	for k, v := range r.definitions {
		next[k] = v
	}
	for _, d := range batch {
		if err := d.Validate(); err != nil {
			return err
		}
		if _, ok := next[d.ID]; ok {
			return &ioc.Error{
				Code: ioc.DuplicateDefinition,
				Bean: d.ID,
			}
		}
		next[d.ID] = d.Clone()
	}
	r.definitions = next
	return nil
}

func (r *mapRegistry) Definitions() []ioc.BeanDefinition {
	var out []ioc.BeanDefinition
	for _, d := range r.definitions {
		out = append(out, d.Clone())
	}
	return out
}

func TestIndependentMetadataRegistryContracts(t *testing.T) {
	for _, r := range []ioc.DefinitionRegistry{&sliceRegistry{}, &mapRegistry{}} {
		d := validDefinition()
		if err := r.Register(d); err != nil {
			t.Fatal(err)
		}
		bad := validDefinition()
		bad.ID = "new"
		duplicate := validDefinition()
		if err := r.Register(bad, duplicate); !errors.Is(err, &ioc.Error{
			Code: ioc.DuplicateDefinition,
		}) {
			t.Fatal(err)
		}
		if len(r.Definitions()) != 1 {
			t.Fatal("failed batch partially published")
		}
		d.Contracts[0] = "changed"
		snapshot := r.Definitions()
		snapshot[0].Contracts[0] = "also changed"
		if r.Definitions()[0].Contracts[0] != "example.Repository" {
			t.Fatal("metadata ownership violated")
		}
	}
}
