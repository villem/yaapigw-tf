terraform {
  required_providers {
    yaapigw = {
      version = "~> 0.1"
      source  = "github.com/villem/yaapigw-tf"
    }
  }
}

provider "yaapigw" {
  username = "education"
  password = "test123"
}

#resource "yaapigw_run_template" "test1" {
#  name = "template_single_fw_for_public_cloud"
#  inputs = {
#    LOG_SERVER_NAME = "LogServer 10.211.55.7"
#    SINGLE_FW__NAME = "YAAPIGW TF  NGFW-3"
#    foo = "bar"
#  }
#}

resource "yaapigw_single_dynamic_firewall" "test2" {
  #inputs  {
  #  log_server_name = "foo"
  #  single_fw_name = "bar"
  #}
  inputs {
    log_server_name = "LogServer 10.211.55.7"
    single_fw__name = "YAAPIGW_single_dynamic_firewall TF NGFW-3"
    extend_override = {
      FOO_BAR = "foo bar"
    }
  }
}

output "test2_initialconfiguration" {
  value = yaapigw_single_dynamic_firewall.test2.outputs["initial_configuration"]
}

#output "test1_init_cfg" {
#  value = yaapigw_run_template.test1.outputs["initial_configuration"]
#}