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

package environment

import "strings"

// ProfileNegationPrefix belongs to the shared predicate grammar rather than
// individual sources, containers or annotation compilers.
const ProfileNegationPrefix = "!"

// ParseProfilePredicate validates a name or one negated name for AcceptsProfiles.
// It deliberately does not evaluate expressions or echo invalid input in errors.
func ParseProfilePredicate(predicate string) (name string, negated bool, err error) {
	name, negated = strings.CutPrefix(predicate, ProfileNegationPrefix)
	if err = ValidateProfile(name); err != nil {
		return "", false, err
	}
	return name, negated, nil
}
