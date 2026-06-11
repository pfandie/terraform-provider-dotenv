data "dotenv_env" "test" {
  keys           = ["TEST_DOTENV_FOO", "TEST_DOTENV_SECRET"]
  sensitive_envs = ["TEST_DOTENV_SECRET"]
}
