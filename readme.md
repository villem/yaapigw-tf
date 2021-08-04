# yaapigw-tf a Terraform provider for YAAPIGW

Yaapigw-tf implements the following providers:

- *yaapigw_run_template* a generic run any kind of template resource
- *yaapigw_single_dynamic_firewall* a specialized resource that only runs YAAPIGW single dynamic firewall template. 


## Build and local install

Execute command `make install`

## Development

Yaapigw-tf depends on yaapigw-go-client (YGC). Adding new templates do not require changes to YGC.

In case that yaapigw-go-client development is need too, the best is to clone YGC repository and use it locally.
First git clone the yaapigw-tf and then YGC.
Eg. 

`
cd ..
git clone github.com/villem/yaapigw-go-client
cd yaapigw-tf-dir
go mod edit -replace github.com/villem/yaapigw-go-client=../yaapigw-go-client
`

 

