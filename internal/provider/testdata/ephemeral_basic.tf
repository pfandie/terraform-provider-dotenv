ephemeral "dotenv_file" "test" {
  env_path       = "testdata"
  env_file       = "valid.env"
  sensitive_envs = ["SECRET_TOKEN"]
}

# The echo provider transfers ephemeral values into state,
# so the test can assert on them. Test-only construct.
provider "echo" {
  data = ephemeral.dotenv_file.test.sensitive_values
}

resource "echo" "test" {}
