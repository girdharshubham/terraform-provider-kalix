terraform {
  required_providers {
    kalix = {
      source = "girdharshubham/kalix"
    }
  }
}
provider "kalix" {
  path = "/usr/local/bin/kalix"
}

data "kalix_cli_version" "cli_version" {}