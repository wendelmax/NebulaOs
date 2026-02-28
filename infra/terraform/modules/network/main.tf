variable "environment" {}

resource "null_resource" "network_placeholder" {
  triggers = {
    env = var.environment
  }
}
