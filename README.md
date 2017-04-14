# terraform-ci

Continuous integration service for [Terraform](https://terraform.io). 

## Planned Features
* Securely stores template variables
* Generates plans when changes committed to source
* Review and apply plans

## Developing
* Backend written in go
* Frontend vue.js app

## Testing
* `make test`

## Running
* Backend: `go run main.go`
* Frontend: `cd site; npm run dev`