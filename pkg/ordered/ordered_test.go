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

package ordered_test

import (
	"testing"

	"github.com/longxi-go/gosr/pkg/ordered"
)

type ordinary int

func (v ordinary) Order() int {
	return int(v)
}

type priority struct {
	ordinary
}

func (priority) PriorityOrder() {
}

func TestPriorityMembershipRequiresExplicitMarker(t *testing.T) {
	var normal ordered.Ordered = ordinary(-100)
	if _, ok := normal.(ordered.PriorityOrdered); ok {
		t.Fatal("ordinary order unexpectedly confers priority membership")
	}
	var first ordered.PriorityOrdered = priority{
		ordinary(100),
	}
	if first.Order() != 100 || normal.Order() != -100 {
		t.Fatal("membership must not rewrite numeric order")
	}
}

func TestBoundariesCoverNativeIntOrders(t *testing.T) {
	minimum := -int(^uint(0)>>1) - 1
	maximum := int(^uint(0) >> 1)
	if ordinary(minimum).Order() < ordered.HighestPrecedence || ordinary(maximum).Order() > ordered.LowestPrecedence {
		t.Fatal("published bounds exclude a valid native-int order")
	}
	if ordered.HighestPrecedence >= 0 || ordered.LowestPrecedence <= 0 {
		t.Fatal("order boundaries do not span negative and positive extension points")
	}
}
