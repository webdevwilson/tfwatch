# terraform-ci

Continuous integration service for [Terraform](https://terraform.io). 

## Current State
* Requires terraform repositories to be checked out to a directory (CHECKOUT_DIR) and to have a terraform.tfplan

## Planned Features
* Securely stores template variables
* Generates plans when changes committed to source
* Review and apply plans

## Configuration
* PORT - (default: 3000) The port the HTTP server will bind to
* CHECKOUT_DIR - (default: .state/projects) The directory that contains Terraform repositories. These must have a terraform.tfplan in them.
* STATE_PATH - (default: .state) The location where state is stored on disk
* LOG_LEVEL

## Developing

You can download the latest release from the [Releases](https://github.com/webdevwilson/terraform-ci/releases) page.

* Backend written in go
* Frontend vue.js app

### API Endpoints

#### Projects

* `/status` - (GET) Get service status
* `/api/projects` - (GET, PUT) List all projects, create project
* `/api/projects/{guid}` - (POST, DELETE) Update or delete projects
* `/api/plan/{guid}` - (GET) Return the current plan associated with the project guid

## Testing
* `make test`

## Running
* Backend: `go run main.go`
* Frontend: `cd site; npm run dev`