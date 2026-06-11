// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &FileEphemeralResource{}

type FileEphemeralResource struct{}

func NewFileEphemeralResource() ephemeral.EphemeralResource {
	return &FileEphemeralResource{}
}

func (e *FileEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (e *FileEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Loads a `.env` file ephemerally — values are never persisted in the plan or state file. Requires Terraform >= 1.10.",
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
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("exclude_envs")),
				},
			},
			"exclude_envs": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"sensitive_envs": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"include_empty": schema.BoolAttribute{
				Optional: true,
			},
			"values": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"sensitive_values": schema.MapAttribute{
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
		},
	}
}

func (e *FileEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data FileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(populateFileModel(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
