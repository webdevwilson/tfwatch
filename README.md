# terraform-ci

Continuous integration service for [Terraform](https://terraform.io). 

## Current State
* Requires terraform repositories to be checked out to a directory (CHECKOUT_DIR) and to have a terraform.tfplan

## Planned Features
* Securely stores template variables
* Generates plans when changes committed to source
* Review and apply plans

## Configure with environment variables
* **PORT** - The port the HTTP server will bind to. Default is `3000`.
* **CHECKOUT_DIR** - The directory that contains Terraform repositories. These must have a terraform.tfplan in them. Default is `.state/projects`.
* **STATE_PATH** - The location where state is stored on disk. Default is `.state/projects`.
* **LOG_LEVEL** - Valid values are: `DEBUG`, `INFO`, `WARN`, `ERROR`. Default is `INFO`.

## Developing

You can download the latest release from the [Releases](https://github.com/webdevwilson/terraform-ci/releases) page.

* Backend written in `go 1.8.1`
* Frontend vue.js app under NodeJS `v6.10.2`

### API Endpoints

* **/status** - `GET` Get service status
* **/api/projects** - `GET`,`PUT` List all projects, create project
* **/api/projects/{guid}** - `POST`,`DELETE` Update or delete projects
* **/api/projects/{guid}/tfplan** - `GET` Return the current plan associated with the project guid

## Testing
* `make test`

## Running
* Backend: `go run main.go`
* Frontend: `cd site; npm install; npm run dev`
