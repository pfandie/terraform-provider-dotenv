# Read specific keys from the process environment (e.g. set by a CI pipeline)
data "dotenv_env" "pipeline" {
  keys           = ["DATABASE_URL", "NODE_ENV", "API_KEY"]
  sensitive_envs = ["API_KEY"]
}

# Merge with values from a .env file — file takes precedence
data "dotenv_file" "app" {
  env_file = ".env.${var.environment}"
}

resource "aws_lambda_function" "this" {
  # ...
  environment {
    variables = merge(
      data.dotenv_env.pipeline.values,
      data.dotenv_file.app.values,
    )
  }
}
