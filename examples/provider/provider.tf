terraform {
  required_providers {
    dotenv = {
      source  = "pfandie/dotenv"
      version = ">= 0.1.0"
    }
  }
}

provider "dotenv" {}
