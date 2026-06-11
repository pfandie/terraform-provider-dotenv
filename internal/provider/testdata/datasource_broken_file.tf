data "dotenv_file" "test" {
  env_path = "testdata"
  env_file = "broken.env"
}
