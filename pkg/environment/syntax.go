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

const (
	// PlaceholderPrefix opens a property reference under PropertyResolver's grammar.
	PlaceholderPrefix = "${"
	// PlaceholderSuffix closes a reference, including a nested default.
	PlaceholderSuffix = "}"
	// PlaceholderDefaultSeparator separates a key from its missing-value fallback.
	PlaceholderDefaultSeparator = ":"
	// PlaceholderEscapePrefix protects the following placeholder from expansion.
	PlaceholderEscapePrefix = "\\"
	// EscapedPlaceholderPrefix is derived so parsers share the same escape grammar.
	EscapedPlaceholderPrefix = PlaceholderEscapePrefix + PlaceholderPrefix
	// DefaultProfile is the fallback profile when no defaults are supplied.
	DefaultProfile = "default"
)
