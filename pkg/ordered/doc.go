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

// Package ordered defines reusable precedence contracts, independent of IoC.
// PriorityOrdered forms the first group; ordinary Ordered and unordered objects
// share the second. Smaller Order values run first within a group; an unordered
// object uses LowestPrecedence. Feature owners define meaningful local anchors.
//
// Consumers resolve objects and sample keys once per ordering round, then compare
// detached values without subtraction overflow. The package owns no comparator,
// object identity, cache or scheduling policy. Mutable order implementations own
// synchronization; changes can take effect at the next ordering round.
package ordered
