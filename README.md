# terraform-ci 

Continuous integration service for [Terraform](https://terraform.io). 

* Build Status - [![CircleCI](https://circleci.com/gh/webdevwilson/terraform-ci.svg?style=svg)](https://circleci.com/gh/webdevwilson/terraform-ci)

## Current State

* Requires terraform repositories to be checked out to a directory (CHECKOUT_DIR) and to have a terraform.tfplan

### Quickstart

Get it up and running quickly in docker:

`docker run -v $CHECKOUT_DIR:/var/lib/terraform-ci --expose=3000 webdevwilson/terraform-ci /var/lib/terraform-ci`

## Planned Features

* Securely stores template variables
* Generates plans when changes committed to source
* Review and apply plans

## Configuration

### Command Line Flags

To see the command line flags for configuration.

`terraform-ci -h`

### Environment Variables

* **CHECKOUT_DIR** - The directory that contains Terraform repositories. These must have a terraform.tfplan in them. Default is `/var/lib/terraform-ci`.
* **CLEAR_STATE** - Clear the state when this variable is set. Default is `false`.
* **LOG_LEVEL** - Valid values are: `DEBUG`, `INFO`, `WARN`, `ERROR`. Default is `INFO`.
* **PLAN_INTERVAL** - The number of minutes between plan refreshes. Default is `5`.
* **PORT** - The port the HTTP server will bind to. Default is `3000`.
* **STATE_DIR** - The location where state is stored on disk. Default is `.terraform-ci-data/projects`.
* **SITE_DIR** - Directory containing static site resources. Default is `site/dist`.

## Developing

You can download the latest release from the [Releases](https://github.com/webdevwilson/terraform-ci/releases) page.

### Backend

The backend is written using `go 1.8.1`. [govendor](https://github.com/kardianos/govendor) is used for vendoring. To install the tools
you will need to develop on the backend, run `make tools`.

To run the backend with live reload run `fresh`. 

**You must have `$GOPATH/bin` on your `$PATH` for fresh to work**

### Frontend

Frontend vue.js app under NodeJS `v6.10.2`

To run the front end, run `npm run dev` from the `/site` directory. You will also need the backend running `go run main.go` to respond to API requests.

**NOTE: To run V2 of the site, run `npm run dev` from the `siteV2` directory**

### API Endpoints

* **/status** - `GET` Get service status
* **/api/projects** - `GET`,`PUT` List all projects, create project
* **/api/projects/{guid}** - `POST`,`DELETE` Update or delete projects
* **/api/projects/{guid}/tfplan** - `GET` Return the current plan associated with the project guid

### 

## Testing
* `make test`

## Running
* Backend: `go run main.go`
* Frontend: `cd site; npm install; npm run dev`
