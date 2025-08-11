# Copyright (c) HashiCorp, Inc.

resource "dt_project" "test" {
  organization = "organizations/cvinmt9aq9sc738g6eog"
  display_name = "Contact Test Project"
  location     = {}
}

resource "dt_contact_group" "test" {
  organization = "organizations/cvinmt9aq9sc738g6eog"
  display_name = "Store Employees"
}

resource "dt_contact" "test" {
  contact_group = dt_contact_group.test.name
  project       = dt_project.test.name
  email         = "some.one.else@example.com"
  display_name  = "Some One Else"
  phone_number  = "+1234567890"
}
