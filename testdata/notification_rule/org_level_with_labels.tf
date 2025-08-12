# Copyright (c) HashiCorp, Inc.
resource "dt_project" "test" {
  display_name = "Test Project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location     = {}
}

resource "dt_notification_rule" "test" {
  display_name         = "Organization Level Alert with Labels"
  parent_resource_name = "organizations/cvinmt9aq9sc738g6eog"
  project_labels = {
    "StoreID" = ""
  }



  trigger = {
    field = "temperature"
    range = {
      lower = 0
      upper = 30
    }
  }
  escalation_levels = [
    {
      display_name = "Email Someone"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            recipients = [
              "someone@example.com"
            ]
            body    = "Temperature $celsius°C is out of range for organization"
            subject = "$projectDisplayName Temperature Alert"
          }
        }
      ]
    }
  ]
}
