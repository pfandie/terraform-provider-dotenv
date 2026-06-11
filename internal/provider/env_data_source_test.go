// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEnvDataSource_basic(t *testing.T) {
	t.Setenv("TEST_DOTENV_FOO", "foo_value")
	t.Setenv("TEST_DOTENV_BAR", "bar_value")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/env_basic.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_env.test",
						tfjsonpath.New("values"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"TEST_DOTENV_FOO": knownvalue.StringExact("foo_value"),
							"TEST_DOTENV_BAR": knownvalue.StringExact("bar_value"),
						}),
					),
				},
			},
		},
	})
}

func TestAccEnvDataSource_sensitive(t *testing.T) {
	t.Setenv("TEST_DOTENV_FOO", "foo_value")
	t.Setenv("TEST_DOTENV_SECRET", "s3cr3t")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/env_sensitive.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_env.test",
						tfjsonpath.New("values"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"TEST_DOTENV_FOO": knownvalue.StringExact("foo_value"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.dotenv_env.test",
						tfjsonpath.New("sensitive_values"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"TEST_DOTENV_SECRET": knownvalue.StringExact("s3cr3t"),
						}),
					),
					statecheck.ExpectSensitiveValue(
						"data.dotenv_env.test",
						tfjsonpath.New("sensitive_values"),
					),
				},
			},
		},
	})
}

func TestAccEnvDataSource_missing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/env_missing.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_env.test",
						tfjsonpath.New("values"),
						knownvalue.MapExact(map[string]knownvalue.Check{}),
					),
				},
			},
		},
	})
}
