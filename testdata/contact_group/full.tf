# Copyright (c) HashiCorp, Inc.

resource "dt_contact_group" "test" {
  organization = "organizations/cvinmt9aq9sc738g6eog"
  display_name = "Contact group with display name and description"
  description  = "This is a full contact group with all attributes."
}
