# Copyright (c) HashiCorp, Inc.

resource "dt_contact_group" "myContactGroup" {
  organization = "organizations/cvinmt9aq9sc738g6eog"
  display_name = "Assistant to the Regional Manager"
  description  = "Not the Assistant Regional Manager, but the Assistant to the Regional Manager."
}
