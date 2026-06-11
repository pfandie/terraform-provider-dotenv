// Copyright (c) 2026 Hans Mayer (pfandie)
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joho/godotenv"
)

type parseOptions struct {
	includeEnvs   []string
	excludeEnvs   []string
	sensitiveEnvs []string
	includeEmpty  bool
}

func parseEnvFile(ctx context.Context, path string, opts parseOptions) (map[string]string, map[string]string, error) {
	tflog.Debug(ctx, "loading .env file", map[string]interface{}{
		"env_path":      path,
		"include_empty": opts.includeEmpty,
	})

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("env file %q could not be opened: %w", path, err)
	}
	defer f.Close()

	envMap, err := godotenv.Parse(f)
	if err != nil {
		return nil, nil, fmt.Errorf("env file %q could not be parsed: %w", path, err)
	}

	values := map[string]string{}
	sensitiveValues := map[string]string{}

	for key, value := range envMap {
		if value == "" && !opts.includeEmpty {
			continue
		}

		if len(opts.includeEnvs) > 0 && !slices.Contains(opts.includeEnvs, key) {
			continue
		}
		if slices.Contains(opts.excludeEnvs, key) {
			continue
		}

		if slices.Contains(opts.sensitiveEnvs, key) {
			sensitiveValues[key] = value
		} else {
			values[key] = value
		}
	}

	// only key counts are logged, values may contain secrets
	tflog.Debug(ctx, "parsed .env file", map[string]interface{}{
		"env_path":       path,
		"keys":           len(values),
		"sensitive_keys": len(sensitiveValues),
	})

	return values, sensitiveValues, nil
}

// populateFileModel parses the env file and writes computed values back into
// data. Shared by the data source Read and the ephemeral resource Open.
func populateFileModel(ctx context.Context, data *FileModel) diag.Diagnostics {
	var diags diag.Diagnostics

	opts, envPath, envFile, d := buildOptions(ctx, data)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	values, sensitiveValues, err := parseEnvFile(ctx, filepath.Join(envPath, envFile), opts)
	if err != nil {
		diags.AddError("Unable to Load .env File", err.Error())
		return diags
	}

	data.EnvPath = types.StringValue(envPath)
	data.EnvFile = types.StringValue(envFile)

	valMap, d := types.MapValueFrom(ctx, types.StringType, values)
	diags.Append(d...)
	sensMap, d := types.MapValueFrom(ctx, types.StringType, sensitiveValues)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	data.Values = valMap
	data.SensitiveValues = sensMap

	return diags
}

// buildOptions normalizes the config: env_path defaults to the current
// working directory, env_file to ".env". The returned strings are the
// normalized directory and file name, ready for filepath.Join.
func buildOptions(ctx context.Context, data *FileModel) (parseOptions, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	opts := parseOptions{}

	envPath := "."
	if !data.EnvPath.IsNull() && data.EnvPath.ValueString() != "" {
		envPath = data.EnvPath.ValueString()
	}

	envFile := ".env"
	if !data.EnvFile.IsNull() && data.EnvFile.ValueString() != "" {
		envFile = data.EnvFile.ValueString()
	}

	if !data.IncludeEnvs.IsNull() {
		diags.Append(data.IncludeEnvs.ElementsAs(ctx, &opts.includeEnvs, false)...)
	}

	if !data.ExcludeEnvs.IsNull() {
		diags.Append(data.ExcludeEnvs.ElementsAs(ctx, &opts.excludeEnvs, false)...)
	}

	if !data.SensitiveEnvs.IsNull() {
		diags.Append(data.SensitiveEnvs.ElementsAs(ctx, &opts.sensitiveEnvs, false)...)
	}

	opts.includeEmpty = data.IncludeEmpty.ValueBool()

	return opts, envPath, envFile, diags
}
