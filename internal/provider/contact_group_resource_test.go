// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSafeContactGroupResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a contact group
			{
				Config: providerConfig + readTestFile(t, "../../testdata/contact_group/minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_contact_group.test", "display_name", "Contact group with display name only"),
					resource.TestCheckResourceAttr("dt_contact_group.test", "contact_count", "0"),
				),
			},
			// update with description
			{
				Config: providerConfig + readTestFile(t, "../../testdata/contact_group/full.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dt_contact_group.test", "display_name", "Contact group with display name and description"),
					resource.TestCheckResourceAttr("dt_contact_group.test", "description", "This is a full contact group with all attributes."),
					resource.TestCheckResourceAttr("dt_contact_group.test", "contact_count", "0"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_contact_group.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_contact_group.test"].Primary.Attributes["name"], nil
				},
			},
		},
	})
}
