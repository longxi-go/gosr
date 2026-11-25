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

package environment_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/environment"
)

type privateFailure struct {
	reason string
}

func (e *privateFailure) Error() string {
	return e.reason
}

func TestErrorClassificationCauseAndRedaction(t *testing.T) {
	cause := &privateFailure{
		reason: "password=do-not-print",
	}
	err := &environment.Error{
		Code: environment.SourceFailed,
		Key:  "db.password",
		Origin: environment.Origin{
			Source:   "config",
			Resource: "settings.json",
			Key:      "db.password",
			Line:     3,
			Column:   2,
		},
		Path:    []string{"db.password", "secrets.password"},
		Message: "source lookup failed",
		Cause:   cause,
	}
	wrapped := fmt.Errorf("operation: %w", err)
	if !errors.Is(wrapped, &environment.Error{
		Code: environment.SourceFailed,
	}) || !errors.Is(wrapped, cause) {
		t.Fatal("classification or cause lost")
	}
	for _, target := range []error{&environment.Error{}, &environment.Error{
		Code: environment.InvalidValue,
	}, (*environment.Error)(nil), errors.New("other")} {
		if errors.Is(wrapped, target) {
			t.Errorf("unexpected match: %v", target)
		}
	}
	var detail *environment.Error
	var underlying *privateFailure
	if !errors.As(wrapped, &detail) ||
		detail != err ||
		!errors.As(wrapped, &underlying) ||
		underlying != cause {
		t.Fatal("structured error or cause inaccessible")
	}
	text := err.Error()
	if strings.Contains(text, cause.reason) {
		t.Fatal("error exposed cause values")
	}
	if !strings.Contains(text, string(environment.SourceFailed)) || !strings.Contains(text, "db.password") {
		t.Errorf("missing diagnostic identity: %q", text)
	}
	if detail.Origin.Line != 3 || detail.Path[1] != "secrets.password" {
		t.Fatal("structured provenance lost")
	}
	var nilError *environment.Error
	if nilError.Error() != "<nil>" || nilError.Unwrap() != nil || nilError.Is(err) {
		t.Fatal("nil error behavior")
	}
	if (&environment.Error{}).Is(&environment.Error{}) {
		t.Fatal("empty codes match")
	}
}

func TestValidation(t *testing.T) {
	for _, key := range []string{
		"app.port",
		"a",
		"app",
		"app[0]",
		"a b",
		"配置.端口",
		"a\tb",
		"a\nb",
		"${literal}",
		"!literal",
	} {
		if err := environment.ValidateKey(key); err != nil {
			t.Errorf("valid key %q: %v", key, err)
		}
	}
	for _, key := range []string{"", " ", "\t", " app", "app ", "\u2003app", "app\n", "app\u00a0"} {
		if err := environment.ValidateKey(key); !errors.Is(err, &environment.Error{
			Code: environment.InvalidKey,
		}) {
			t.Errorf("invalid key %q: %v", key, err)
		}
	}
	for _, profile := range []string{"a", "0", "default", "dev.us-east_1", "A-Z_1.0"} {
		if err := environment.ValidateProfile(profile); err != nil {
			t.Errorf("valid profile %q: %v", profile, err)
		}
	}
	for _, profile := range []string{
		"",
		" ",
		" dev",
		"dev ",
		"dev prod",
		"!dev",
		"dev|prod",
		"dev&prod",
		"(dev)",
		".dev",
		"-dev",
		"_dev",
		"développement",
		"生产",
		"dev/qa",
		"dev\n",
	} {
		if err := environment.ValidateProfile(profile); !errors.Is(err, &environment.Error{
			Code: environment.InvalidProfile,
		}) {
			t.Errorf("invalid profile %q: %v", profile, err)
		}
	}
}
