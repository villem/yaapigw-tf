module example.com/yaapigw_template

go 1.16

replace example.com/yaapigw_client => ../yaapigw_client

require (
	example.com/yaapigw_client v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/hashicorp/terraform-provider-hashicups v0.3.1 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
)
