// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"
)

const (
	validEnvFile  = "testdata/valid.env"
	brokenEnvFile = "testdata/broken.env"
)

// validEnvValues are the non-empty, non-sensitive key/value pairs in testdata/valid.env.
var validEnvValues = map[string]string{
	"THIS_IS_A_ENV_KEY": "THIS_IS_A_ENV_VALUE",
	"ANOTHER_KEY":       "ANOTHER_VALUE",
	"KEY_WITH_SPACES":   "VALUE WITH SPACES",
	"QUOTES":            "quoted value",
	"INLINE_COMMENT":    "value",
	"EXPORTED_KEY":      "true",
	"SECRET_TOKEN":      "top_secret",
}

func TestParseEnvFile_Valid(t *testing.T) {
	values, sensitive, err := parseEnvFile(t.Context(), validEnvFile, parseOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != len(validEnvValues) {
		t.Errorf("amount of keys: got %d, expected: %d (got: %v)", len(values), len(validEnvValues), values)
	}

	for k, v := range validEnvValues {
		if values[k] != v {
			t.Errorf("key %q: got %q, expected: %q", k, values[k], v)
		}
	}
	if len(sensitive) != 0 {
		t.Errorf("expected no sensitive keys, got: %v", sensitive)
	}

	for _, k := range []string{"EMPTY", "EMPTY_WITH_QUOTES"} {
		if _, ok := values[k]; ok {
			t.Errorf("empty key %q loaded with include_empty=false set", k)
		}
	}
}

func TestParseEnvFile_IncludedEmpty(t *testing.T) {
	values, _, err := parseEnvFile(t.Context(), validEnvFile, parseOptions{includeEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, k := range []string{"EMPTY", "EMPTY_WITH_QUOTES"} {
		v, ok := values[k]
		if !ok {
			t.Errorf("key %q is missing with include_empty=true", k)
		}
		if v != "" {
			t.Errorf("key %q: got %q, expected \"\"", k, v)
		}
	}
}

func TestParseEnvFile_Exclude(t *testing.T) {
	values, _, err := parseEnvFile(t.Context(), validEnvFile, parseOptions{excludeEnvs: []string{"ANOTHER_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := values["ANOTHER_KEY"]; ok {
		t.Error("ANOTHER_KEY should be excluded but is present")
	}
	if values["THIS_IS_A_ENV_KEY"] != "THIS_IS_A_ENV_VALUE" {
		t.Errorf("THIS_IS_A_ENV_KEY is missing or wrong: %v", values)
	}
}

func TestParseEnvFile_Include(t *testing.T) {
	values, _, err := parseEnvFile(t.Context(), validEnvFile, parseOptions{includeEnvs: []string{"ANOTHER_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(values) != 1 || values["ANOTHER_KEY"] != "ANOTHER_VALUE" {
		t.Errorf("include is not correct: %v", values)
	}
}

func TestParseEnvFile_Sensitive(t *testing.T) {
	values, sensitive, err := parseEnvFile(t.Context(), validEnvFile, parseOptions{sensitiveEnvs: []string{"SECRET_TOKEN"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := values["SECRET_TOKEN"]; ok {
		t.Error("SECRET_TOKEN must not be in values")
	}
	if sensitive["SECRET_TOKEN"] != "top_secret" {
		t.Errorf("SECRET_TOKEN is missing in sensitive_values: %v", sensitive)
	}
	if values["ANOTHER_KEY"] != "ANOTHER_VALUE" {
		t.Errorf("ANOTHER_KEY is missing in values: %v", values)
	}
}

func TestParseEnvFile_FileNotFound(t *testing.T) {
	_, _, err := parseEnvFile(t.Context(), "/just/an/imaginary/path/.env", parseOptions{})
	if err == nil {
		t.Fatal("expected an error, file does not exist on purpose")
	}
}

func TestParseEnvFile_BrokenContent(t *testing.T) {
	_, _, err := parseEnvFile(t.Context(), brokenEnvFile, parseOptions{})
	if err == nil {
		t.Fatal("expected a parse error for invalid content")
	}
}
