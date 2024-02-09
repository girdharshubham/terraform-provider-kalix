terraform {
  required_providers {
    kalix = {
      source = "girdharshubham/kalix"
    }
  }
}

provider "kalix" {
  cli_path = "/usr/local/bin/kalix"
  # or KALIX_REFRESH_TOKEN
  #  refresh_token = ""
}

data "kalix_cli_version" "cli_version" {
}

data "kalix_projects" "kalix_projects" {
}

