# Loads .env from the current working directory
data "dotenv_file" "local" {}

# All options
data "dotenv_file" "lambda" {
  env_path       = "../../.."
  env_file       = ".env.${var.environment}"
  exclude_envs   = ["NUXT_PUBLIC_FORCE_FASTLY_IMAGE_HOST"]
  sensitive_envs = ["NUXT_PUBLIC_SCRIPTS_USERCENTRICS_ID"]
  include_empty  = false
}

resource "aws_lambda_function" "this" {
  # ...
  environment {
    variables = merge(
      data.dotenv_file.lambda.values,
      data.dotenv_file.lambda.sensitive_values,
    )
  }
}
