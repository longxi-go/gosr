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

const (
	expectedASCIIProfileNameMessage                        = "expected ASCII profile name"
	expectedNonemptyKeyWithoutSurroundingWhitespaceMessage = "expected nonempty key without surrounding whitespace"
	expectedNonemptyProfileNameMessage                     = "expected nonempty profile name"
)

// ValidateKey admits nonempty exact addresses without surrounding Unicode
// whitespace. Interior whitespace and punctuation are literal; sources may
// normalize before publication, but resolvers must not infer aliases.
func ValidateKey(key string) error {
	if key == "" || strings.TrimSpace(key) != key {
		return &Error{
			Code:    InvalidKey,
			Message: expectedNonemptyKeyWithoutSurroundingWhitespaceMessage,
		}
	}
	return nil
}

// ValidateProfile admits an ASCII letter/digit followed by ASCII letters/digits,
// dots, underscores or hyphens. Negation belongs to AcceptsProfiles predicates,
// never to stored names. Rejections do not echo the supplied text.
func ValidateProfile(profile string) error {
	for i := 0; i < len(profile); i++ {
		c := profile[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') {
			continue
		}
		if i > 0 &&
			(c == '.' || c == '_' || c == '-') {
			continue
		}
		return &Error{
			Code:    InvalidProfile,
			Message: expectedASCIIProfileNameMessage,
		}
	}
	if profile == "" {
		return &Error{
			Code:    InvalidProfile,
			Message: expectedNonemptyProfileNameMessage,
		}
	}
	return nil
}
