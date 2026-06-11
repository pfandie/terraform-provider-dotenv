// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFileEphemeralResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// Ephemeral resources: Terraform >= 1.10, OpenTofu >= 1.11
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		Steps: []resource.TestStep{
			{
				// the config file contains a provider block (echo), therefore the
				// factories have to be declared on the step instead of the test case
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"dotenv": testAccProtoV6ProviderFactories["dotenv"],
					"echo":   echoprovider.NewProviderServer(),
				},
				// Ephemeral resources are re-read on every plan; a non-empty plan after apply is expected.
				ExpectNonEmptyPlan: true,
				ConfigFile:         config.StaticFile("testdata/ephemeral_basic.tf"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"SECRET_TOKEN": knownvalue.StringExact("top_secret"),
						}),
					),
				},
			},
		},
	})
}
