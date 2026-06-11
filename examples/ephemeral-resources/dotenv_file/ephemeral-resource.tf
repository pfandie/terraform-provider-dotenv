# Values are not persisted in the state file or plan
ephemeral "dotenv_file" "secrets" {
  env_path = "../../.."
  env_file = ".env.secrets"
}

resource "aws_secretsmanager_secret_version" "app" {
  secret_id                = aws_secretsmanager_secret.app.id
  secret_string_wo         = jsonencode(ephemeral.dotenv_file.secrets.values)
  secret_string_wo_version = 1
}
