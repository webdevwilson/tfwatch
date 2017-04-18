package model

import "testing"

var model = []Project{
	Project{"", "foo", "git@github.com:webdevwilson/terraform-ci.git", "/foo"},
	Project{"", "bar", "git@github.com:webdevwilson/terraform-ci.git", "/bar"},
}

func TestList(t *testing.T) {

}
