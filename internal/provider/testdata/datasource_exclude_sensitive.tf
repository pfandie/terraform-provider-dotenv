data "dotenv_file" "test" {
  env_path       = "testdata"
  env_file       = "valid.env"
  exclude_envs   = ["EXPORTED_KEY"]
  sensitive_envs = ["SECRET_TOKEN"]
}
