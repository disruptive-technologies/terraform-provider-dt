// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSafeContactResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a contact
			{
				Config: providerConfig + readTestFile(t, "../../testdata/contact/minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_contact.test", "email", "some.one@example.com"),
				),
			},
			// update with description
			{
				Config: providerConfig + readTestFile(t, "../../testdata/contact/full.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dt_contact.test", "display_name", "Some One Else"),
					resource.TestCheckResourceAttr("dt_contact.test", "email", "some.one.else@example.com"),
					resource.TestCheckResourceAttr("dt_contact.test", "phone_number", "+1234567890"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_contact.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_contact.test"].Primary.Attributes["name"], nil
				},
			},
		},
	})
}
