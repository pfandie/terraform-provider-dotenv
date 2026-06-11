// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"os"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &EnvDataSource{}

type EnvDataSource struct{}

func NewEnvDataSource() datasource.DataSource {
	return &EnvDataSource{}
}

type EnvModel struct {
	Keys            types.List `tfsdk:"keys"`
	SensitiveEnvs   types.List `tfsdk:"sensitive_envs"`
	Values          types.Map  `tfsdk:"values"`
	SensitiveValues types.Map  `tfsdk:"sensitive_values"`
}

func (d *EnvDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_env"
}

func (d *EnvDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the specified keys from the process environment and exposes their values to Terraform. Keys not present in the environment are silently omitted.",
		Attributes: map[string]schema.Attribute{
			"keys": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Environment variable names to read from the process environment.",
			},
			"sensitive_envs": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Keys whose values are returned in `sensitive_values` instead of `values`.",
			},
			"values": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Key/value pairs for all set keys not listed in `sensitive_envs`.",
			},
			"sensitive_values": schema.MapAttribute{
				Computed:            true,
				Sensitive:           true,
				ElementType:         types.StringType,
				MarkdownDescription: "Key/value pairs for keys listed in `sensitive_envs` (marked as sensitive).",
			},
		},
	}
}

func (d *EnvDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnvModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var keys []string
	resp.Diagnostics.Append(data.Keys.ElementsAs(ctx, &keys, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sensitiveEnvs []string
	if !data.SensitiveEnvs.IsNull() {
		resp.Diagnostics.Append(data.SensitiveEnvs.ElementsAs(ctx, &sensitiveEnvs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	values := map[string]string{}
	sensitiveValues := map[string]string{}

	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if !ok {
			tflog.Debug(ctx, "environment variable not set, skipping", map[string]interface{}{"key": key})
			continue
		}
		if slices.Contains(sensitiveEnvs, key) {
			sensitiveValues[key] = val
		} else {
			values[key] = val
		}
	}

	// only key counts are logged, values may contain secrets
	tflog.Debug(ctx, "read environment variables", map[string]interface{}{
		"keys":           len(values),
		"sensitive_keys": len(sensitiveValues),
	})

	valMap, diags := types.MapValueFrom(ctx, types.StringType, values)
	resp.Diagnostics.Append(diags...)
	sensMap, diags := types.MapValueFrom(ctx, types.StringType, sensitiveValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Values = valMap
	data.SensitiveValues = sensMap

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
