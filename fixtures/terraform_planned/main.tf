variable "contents" {
  type = "string"
}

resource "local_file" "file" {
  filename = "foo"
  content  = "bar"
}
