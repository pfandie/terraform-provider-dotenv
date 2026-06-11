data "dotenv_file" "test" {
  include_envs = ["A"]
  exclude_envs = ["B"]
}
