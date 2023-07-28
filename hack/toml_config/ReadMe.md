
## Using the Toml config generator

Run the script generator using 
`./hack/toml_config/toml_config_apply.sh` command

This script will,
- Run the python script to generate Go structs using the /hack/toml_config/configs.csv file
- run `make manifests` to regenerate the CRD yaml files

If the toml_key column is empty, it means it's not a deployment.toml config, rather an operator-specific config being passed on. It's best to update such configs within the Go code itself, leaving the configs.csv to keep track of the WSO2 IS-specific deployment.toml configs.




