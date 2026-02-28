# NebulaOS Global Infrastructure Foundation

terraform {
  required_version = ">= 1.5.0"
}

# Local variables for resource naming and tagging
locals {
  common_name = "nebula-${var.environment}"
  tags = {
    Project     = "NebulaOS"
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# Example: Compute Module (Placeholder)
module "compute" {
  source = "./modules/compute"

  environment = var.environment
  project_name = local.common_name
}

# Example: Network Module (Placeholder)
module "network" {
  source = "./modules/network"

  environment = var.environment
}
