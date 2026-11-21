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
	"context"
	"testing"

	"github.com/longxi-go/gosr/pkg/event"
	"github.com/longxi-go/gosr/pkg/ioc"
)

type refreshedReceiver struct{}

func (refreshedReceiver) OnEvent(context.Context, ioc.ContextRefreshed) error {
	return nil
}

type closingReceiver struct{}

func (closingReceiver) OnEvent(context.Context, ioc.ContextClosing) error {
	return nil
}

var (
	_ event.Listener[ioc.ContextRefreshed] = refreshedReceiver{}
	_ event.Listener[ioc.ContextClosing]   = closingReceiver{}
)

func TestContextNotificationTypeBoundaries(t *testing.T) {
	var refreshed event.Event = ioc.ContextRefreshed{}
	var closing event.Event = ioc.ContextClosing{}
	if _, ok := refreshed.(ioc.ContextClosing); ok {
		t.Fatal("startup barrier also matches closing notification")
	}
	if _, ok := closing.(ioc.ContextRefreshed); ok {
		t.Fatal("closing notification also matches startup barrier")
	}
	for _, value := range []event.Event{&ioc.ContextRefreshed{}, &ioc.ContextClosing{}} {
		if _, ok := value.(ioc.ContextRefreshed); ok {
			t.Fatal("pointer notification implicitly matches a value receiver")
		}
	}
}
