module github.com/villem/yaapigw-tf

go 1.16

replace example.com/yaapigw_client => ../yaapigw_client

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/rs/zerolog v1.23.0
	github.com/villem/yaapigw-go-client v0.0.0-20210730030013-99bfab02169e
)
