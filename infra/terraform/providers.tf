terraform {
  required_providers {
    # Using local providers for Phase 2 foundation
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2.0"
    }
  }
}

provider "local" {}
provider "null" {}
