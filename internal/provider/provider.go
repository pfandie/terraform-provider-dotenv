// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider                       = &DotEnvProvider{}
	_ provider.ProviderWithEphemeralResources = &DotEnvProvider{}
)

type DotEnvProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DotEnvProvider{version: version}
	}
}

func (p *DotEnvProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dotenv"
	resp.Version = p.version
}

func (p *DotEnvProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Loads key/value pairs into Terraform — from .env files (dotenv_file) or directly from environment.",
	}
}

func (p *DotEnvProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *DotEnvProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *DotEnvProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFileDataSource,
		NewEnvDataSource,
	}
}

func (p *DotEnvProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{NewFileEphemeralResource}
}
