// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &FileDataSource{}

type FileDataSource struct{}

func NewFileDataSource() datasource.DataSource {
	return &FileDataSource{}
}

type FileModel struct {
	EnvPath         types.String `tfsdk:"env_path"`
	EnvFile         types.String `tfsdk:"env_file"`
	IncludeEnvs     types.List   `tfsdk:"include_envs"`
	ExcludeEnvs     types.List   `tfsdk:"exclude_envs"`
	SensitiveEnvs   types.List   `tfsdk:"sensitive_envs"`
	IncludeEmpty    types.Bool   `tfsdk:"include_empty"`
	Values          types.Map    `tfsdk:"values"`
	SensitiveValues types.Map    `tfsdk:"sensitive_values"`
}

func (d *FileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *FileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Loads a `.env` file and exposes its key/value pairs to Terraform.",
		Attributes: map[string]schema.Attribute{
			"env_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Directory containing the env file. Defaults to the current working directory (`.`).",
			},
			"env_file": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of the env file inside `env_path`. Defaults to `.env`.",
			},
			"include_envs": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Only these keys will be included. Exclusive to `exclude_envs`.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("exclude_envs")),
				},
			},
			"exclude_envs": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "These keys won't be included. Exclusive to `include_envs`.",
			},
			"sensitive_envs": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "These keys will be set as `sensitive_values` instead of values.",
			},
			"include_empty": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Is set to true, empty Keys will also be loaded (e.g. `FOO=`). Default is `false`",
			},
			"values": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "All Key/Values pairs without `sensitive_envs`.",
			},
			"sensitive_values": schema.MapAttribute{
				Computed:            true,
				Sensitive:           true,
				ElementType:         types.StringType,
				MarkdownDescription: "Sensitive Key/Values pairs defined in `sensitive_envs`.",
			},
		},
	}
}

func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(populateFileModel(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
