variable "environment" {
  description = "Execution environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "region" {
  description = "Target cloud/physical region"
  type        = string
  default     = "local"
}

variable "organization_name" {
  description = "Name of the sovereign organization"
  type        = string
  default     = "NebulaOrg"
}
