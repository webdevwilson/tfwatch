package model

import "testing"

var model = []Project{
	Project{"", "foo", "git@github.com:webdevwilson/terraform-ci.git", "/foo", map[string]string{}},
	Project{"", "bar", "git@github.com:webdevwilson/terraform-ci.git", "/bar", map[string]string{}},
}

func TestList(t *testing.T) {

}
