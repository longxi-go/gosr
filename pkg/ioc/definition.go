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

package ioc

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	aliasesLabel                               = "aliases: "
	contractsLabel                             = "contracts: "
	qualifiersLabel                            = "qualifiers: "
	aContractOrNameIsRequiredMessage           = "a contract or name is required"
	anAliasMustDifferFromThePrimaryNameMessage = "an alias must differ from the primary name"
	atLeastOneExposedContractIsRequiredMessage = "at least one exposed contract is required"
	blankTokenMessage                          = "blank token"
	dependencyAccessMustBeExplicitFormat       = "dependency %d: access must be explicit"
	dependencyCardinalityMustBeExplicitFormat  = "dependency %d: cardinality must be explicit"
	dependencyFormat                           = "dependency %d: %v"
	dependencyNegativeSourceLineFormat         = "dependency %d: negative source line"
	duplicateTokenFormat                       = "duplicate token %q"
	identityAndNameMustBeNonblankMessage       = "identity and name must be nonblank"
	scopeMustBeNonblankMessage                 = "scope must be nonblank"
	sourceLineMustBeNonnegativeMessage         = "source line must be nonnegative"
	suppliedConstraintsMustBeNonblankMessage   = "supplied constraints must be nonblank"
)

// Contract identifies an explicitly exposed capability. Its encoding is opaque
// to gosr; no runtime type discovery or cross-generator ABI is required.
type Contract string

// ScopeID names an instance ownership policy, not a module visibility boundary.
type ScopeID string

const (
	// Singleton retains one object per application context.
	Singleton ScopeID = "singleton"
	// Prototype creates a fresh object per resolution. A non-resource prototype
	// may be retained by its consumer; a resource requires explicit ownership.
	Prototype ScopeID = "prototype"
	// Request retains an object only for an explicitly opened request lifetime.
	Request ScopeID = "request"
)

// Access distinguishes creation edges from caller-invoked provider edges.
type Access string

const (
	Direct   Access = "direct"
	Deferred Access = "deferred"
)

// Cardinality makes absence and collections explicit; pointers remain required.
// One selects a unique/primary candidate. ZeroOrOne permits zero, not ambiguity.
// Many retains all matches ordered by Order then ID, regardless of primary.
type Cardinality string

const (
	One       Cardinality = "one"
	ZeroOrOne Cardinality = "optional"
	Many      Cardinality = "many"
)

// Source carries declaration provenance. A zero line means unknown provenance.
type Source struct {
	Package, File string
	Line          int
}

// Selector combines all populated constraints; Name matches names or aliases.
type Selector struct {
	Contract   Contract
	Name       string
	Qualifiers []string
}

// Dependency is portable binding intent. It contains no executable resolver.
type Dependency struct {
	Selector    Selector
	Access      Access
	Cardinality Cardinality
	Source      Source
}

// BeanDefinition describes one candidate rather than one Go type. Several
// contracts of the same definition share its instance and cleanup identity.
type BeanDefinition struct {
	ID, Name            string
	Aliases, Qualifiers []string
	Contracts           []Contract
	Dependencies        []Dependency
	// DependsOn names initialization prerequisites through primary names or aliases.
	// They finish construction and Init before this definition's factory runs;
	// they do not occupy injection ports or imply Start/readiness guarantees.
	// Implementations reject unresolved targets, cycles and unsupported ownership
	// relations, and release same-owner dependents before their prerequisites.
	DependsOn     []string
	Scope         ScopeID
	Primary, Lazy bool
	Order         int
	Source        Source
}

// Validate checks declaration invariants without choosing an implementation.
func (d BeanDefinition) Validate() error {
	invalid := func(message string, source Source) error {
		return &Error{
			Code:    InvalidDefinition,
			Bean:    d.ID,
			Source:  source,
			Message: message,
		}
	}
	if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.Name) == "" {
		return invalid(identityAndNameMustBeNonblankMessage, d.Source)
	}
	if len(d.Contracts) == 0 {
		return invalid(atLeastOneExposedContractIsRequiredMessage, d.Source)
	}
	if err := unique(d.Contracts); err != nil {
		return invalid(contractsLabel+err.Error(), d.Source)
	}
	if err := unique(d.Aliases); err != nil {
		return invalid(aliasesLabel+err.Error(), d.Source)
	}
	if slices.Contains(d.Aliases, d.Name) {
		return invalid(anAliasMustDifferFromThePrimaryNameMessage, d.Source)
	}
	if err := unique(d.Qualifiers); err != nil {
		return invalid(qualifiersLabel+err.Error(), d.Source)
	}
	if strings.TrimSpace(string(d.Scope)) == "" {
		return invalid(scopeMustBeNonblankMessage, d.Source)
	}
	if d.Source.Line < 0 {
		return invalid(sourceLineMustBeNonnegativeMessage, d.Source)
	}
	if err := validateInitializationReferences(d); err != nil {
		return err
	}
	for i, dep := range d.Dependencies {
		if err := dep.Selector.Validate(); err != nil {
			return invalid(fmt.Sprintf(dependencyFormat, i, err), dep.Source)
		}
		if dep.Access != Direct && dep.Access != Deferred {
			return invalid(fmt.Sprintf(dependencyAccessMustBeExplicitFormat, i), dep.Source)
		}
		if dep.Cardinality != One && dep.Cardinality != ZeroOrOne && dep.Cardinality != Many {
			return invalid(fmt.Sprintf(dependencyCardinalityMustBeExplicitFormat, i), dep.Source)
		}
		if dep.Source.Line < 0 {
			return invalid(fmt.Sprintf(dependencyNegativeSourceLineFormat, i), dep.Source)
		}
	}
	return nil
}

// Clone gives the recipient ownership of every mutable metadata slice.
func (d BeanDefinition) Clone() BeanDefinition {
	d.Aliases = slices.Clone(d.Aliases)
	d.Qualifiers = slices.Clone(d.Qualifiers)
	d.Contracts = slices.Clone(d.Contracts)
	d.Dependencies = slices.Clone(d.Dependencies)
	d.DependsOn = slices.Clone(d.DependsOn)
	for i := range d.Dependencies {
		d.Dependencies[i].Selector.Qualifiers = slices.Clone(d.Dependencies[i].Selector.Qualifiers)
	}
	return d
}

// Matches filters candidates; uniqueness and primary selection belong to binding.
func (d BeanDefinition) Matches(s Selector) bool {
	if s.Contract != "" && !slices.Contains(d.Contracts, s.Contract) {
		return false
	}
	if s.Name != "" && s.Name != d.Name && !slices.Contains(d.Aliases, s.Name) {
		return false
	}
	for _, q := range s.Qualifiers {
		if !slices.Contains(d.Qualifiers, q) {
			return false
		}
	}
	return true
}

// Validate rejects unconstrained or malformed runtime lookup requests.
func (s Selector) Validate() error {
	invalid := func(message string) error {
		return &Error{
			Code:    InvalidSelector,
			Message: message,
		}
	}
	if s.Contract == "" && s.Name == "" {
		return invalid(aContractOrNameIsRequiredMessage)
	}
	blankContract := s.Contract != "" && strings.TrimSpace(string(s.Contract)) == ""
	blankName := s.Name != "" && strings.TrimSpace(s.Name) == ""
	if blankContract || blankName {
		return invalid(suppliedConstraintsMustBeNonblankMessage)
	}
	if err := unique(s.Qualifiers); err != nil {
		return invalid(qualifiersLabel + err.Error())
	}
	return nil
}

func unique[T ~string](values []T) error {
	seen := make(map[T]bool, len(values))
	for _, v := range values {
		if strings.TrimSpace(string(v)) == "" {
			return errors.New(blankTokenMessage)
		}
		if seen[v] {
			return fmt.Errorf(duplicateTokenFormat, v)
		}
		seen[v] = true
	}
	return nil
}
