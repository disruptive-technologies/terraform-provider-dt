# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "All escalation types"
  project_id   = data.dt_project.test.id

  trigger = {
    field = "temperature"
    range = {
      lower = 0
      type  = "OUTSIDE"
    }
  }
  escalation_levels = [
    {
      display_name   = "corrigo"
      escalate_after = "3600s"
      actions = [
        {
          // without studio dashboard url and description
          type = "CORRIGO"
          corrigo_config = {
            asset_id        = "asset-id-1"
            client_id       = "client-id"
            client_secret   = "super-secret"
            company_name    = "company-name"
            contact_address = "contact-address"
            contact_name    = "contact-name"
            customer_id     = "customer-id"
            sub_type_id     = "sub-type-id"
            task_id         = "task-id"
          }
        },
        {
          // full corrigo config
          type = "CORRIGO"
          corrigo_config = {
            asset_id               = "asset-id-2"
            client_id              = "client-id"
            client_secret          = "super-secret"
            company_name           = "company-name"
            contact_address        = "contact-address"
            contact_name           = "contact-name"
            customer_id            = "customer-id"
            sub_type_id            = "sub-type-id"
            task_id                = "task-id"
            work_order_description = "Temperature $celsius is over the limit"
          }
        }
      ]
    },
    {
      display_name   = "email"
      escalate_after = "3600s"
      actions = [{
        type = "EMAIL"
        email_config = {
          body    = "Temperature $celsius is over the limit"
          subject = "Temperature Alert"
          recipients = [
            "someone@example.com"
          ]
        }
      }]
    },
    {
      display_name   = "phone call"
      escalate_after = "3600s"
      actions = [{
        type = "PHONE_CALL"
        phone_call_config = {
          introduction = "This is an automated call from Disruptive Technologies"
          message      = "Temperature $celsius is over the limit for device $name"
          recipients = [
            "+4798765432"
          ]
        }
      }]
    },
    {
      display_name   = "service channel"
      escalate_after = "3600s"
      actions = [{
        type = "SERVICE_CHANNEL"
        service_channel_config = {
          description = "Temperature $celsius is over the limit"
          store_id    = "store-id"
          trade       = "REFRIGERATION"
        }
      }]
    },
    {
      display_name   = "SMS"
      escalate_after = "3600s"
      actions = [
        {
          type = "SMS"
          sms_config = {
            body = "Temperature $celsius is over the limit"
            recipients = [
              "+4798765432"
            ]
          }
        }
      ]
    },
    {
      display_name = "webhook"
      actions = [{
        type = "WEBHOOK"
        webhook_config = {
          url = "https://example.com/webhook"
          headers = {
            "Content-Type" : "application/json"
          }
          signature_secret = "super-secret"
        }
      }]
    }
  ]
}
