variable "environment" {}
variable "project_name" {}

resource "null_resource" "compute_placeholder" {
  triggers = {
    project = var.project_name
  }
}
