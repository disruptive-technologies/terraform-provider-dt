# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name         = "Organization Level Alert"
  parent_resource_name = "organizations/cvinmt9aq9sc738g6eog"


  trigger = {
    field = "temperature"
    range = {
      lower = 0
      upper = 30
    }
  }
  escalation_levels = [
    {
      display_name = "Email General Manager"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            contact_groups = [
              "organizations/cvinmt9aq9sc738g6eog/contactGroups/d2dkclv9a2cc7390cis0"
            ]
            body    = "Temperature $celsius°C is out of range for organization"
            subject = "$projectDisplayName Temperature Alert"
          }
        }
      ]
      escalate_after = "7200s"
    },
    {
      display_name = "SMS General Manager"
      actions = [
        {
          type = "SMS"
          sms_config = {
            contact_groups = [
              "organizations/cvinmt9aq9sc738g6eog/contactGroups/d2dkclv9a2cc7390cis0"
            ]
            body = "Temperature $celsius°C is out of range for organization"
          }
        }
      ]
      escalate_after = "7200s"
    },
    {
      display_name = "Call General Manager"
      actions = [
        {
          type = "PHONE_CALL"
          phone_call_config = {
            contact_groups = [
              "organizations/cvinmt9aq9sc738g6eog/contactGroups/d2dkclv9a2cc7390cis0"
            ]
            introduction = "This is an automated call from Disruptive Technologies"
            message      = "Temperature $celsius°C is out of range for organization"
          }
        }
      ]
    }
  ]
}
