
# Allgemeine Variablen


# Public Variablen

variable "userdata" {
    description = "Cloud-init Script"
    default = "base.yaml"
}

# Scripts

data "template_file" "userdata" {
  template = file(var.userdata)
}
