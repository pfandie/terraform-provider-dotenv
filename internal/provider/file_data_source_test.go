// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dotenv": providerserver.NewProtocol6WithError(New("test")()),
}

// knownValueMap converts plain maps into the exact map check used by statecheck.
func knownValueMap(m map[string]string, extra map[string]string) knownvalue.Check {
	checks := map[string]knownvalue.Check{}
	for k, v := range m {
		checks[k] = knownvalue.StringExact(v)
	}
	for k, v := range extra {
		checks[k] = knownvalue.StringExact(v)
	}
	return knownvalue.MapExact(checks)
}

func TestAccFileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/datasource_basic.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_file.test",
						tfjsonpath.New("values"),
						// EMPTY + EMPTY_WITH_QUOTES missing on purpose: include_empty defaults to false
						knownValueMap(validEnvValues, nil),
					),
				},
			},
		},
	})
}

func TestAccFileDataSource_excludeAndSensitive(t *testing.T) {
	values := map[string]string{}
	for k, v := range validEnvValues {
		values[k] = v
	}
	delete(values, "EXPORTED_KEY")
	delete(values, "SECRET_TOKEN")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/datasource_exclude_sensitive.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_file.test",
						tfjsonpath.New("values"),
						knownValueMap(values, nil),
					),
					statecheck.ExpectKnownValue(
						"data.dotenv_file.test",
						tfjsonpath.New("sensitive_values"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"SECRET_TOKEN": knownvalue.StringExact("top_secret"),
						}),
					),
					statecheck.ExpectSensitiveValue(
						"data.dotenv_file.test",
						tfjsonpath.New("sensitive_values"),
					),
				},
			},
		},
	})
}

func TestAccFileDataSource_includeEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile("testdata/datasource_include_empty.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.dotenv_file.test",
						tfjsonpath.New("values"),
						knownValueMap(validEnvValues, map[string]string{
							"EMPTY":             "",
							"EMPTY_WITH_QUOTES": "",
						}),
					),
				},
			},
		},
	})
}

func TestAccFileDataSource_includeExcludeConflict(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile:  config.StaticFile("testdata/datasource_include_exclude_conflict.tf"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccFileDataSource_fileNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile:  config.StaticFile("testdata/datasource_file_not_found.tf"),
				ExpectError: regexp.MustCompile(`Unable to Load .env File`),
			},
		},
	})
}

func TestAccFileDataSource_brokenFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile:  config.StaticFile("testdata/datasource_broken_file.tf"),
				ExpectError: regexp.MustCompile(`Unable to Load .env File`),
			},
		},
	})
}
