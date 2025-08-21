// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSafeProjectResourceExamples(t *testing.T) {
	t.Parallel()
	t.Log("TestAccProjectResourceExamples")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: providerConfig + readTestFile(t, "../../examples/resources/dt_project/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.my_project", "display_name", "Terraform created project"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.latitude", "63.44539"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.longitude", "10.910202"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.time_location", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_project.my_project", "inventory", "false"),
				),
			},
		},
	})
}

func TestAccSafeProjectResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and read testing
				Config: providerConfig + readTestFile(t, "../../testdata/project/with_location.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test", "display_name", "Acceptance Test Project"),
					resource.TestCheckResourceAttr("dt_project.test", "location.latitude", "63.44539"),
					resource.TestCheckResourceAttr("dt_project.test", "location.longitude", "10.910202"),
					resource.TestCheckResourceAttr("dt_project.test", "location.time_location", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_project.test", "inventory", "false"),
					resource.TestCheckResourceAttr("dt_project.test", "labels.%", "0"),
				),
			},
			{
				// Import testing
				ResourceName:                         "dt_project.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_project.test"].Primary.Attributes["name"], nil
				},
			},
			{
				// Update testing
				Config: providerConfig + readTestFile(t, "../../testdata/project/with_updated_location.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test", "display_name", "Acceptance Test Project"),
					resource.TestCheckResourceAttr("dt_project.test", "location.latitude", "59.910953"),
					resource.TestCheckResourceAttr("dt_project.test", "location.longitude", "10.63904"),
					resource.TestCheckResourceAttr("dt_project.test", "location.time_location", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_project.test", "inventory", "false"),
					resource.TestCheckResourceAttr("dt_project.test", "labels.%", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + readTestFile(t, "../../testdata/project/empty_location.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test", "display_name", "Empty Location Project"),
					resource.TestCheckResourceAttr("dt_project.test", "location.time_location", "UTC"),
					resource.TestCheckResourceAttr("dt_project.test", "inventory", "false"),
					resource.TestCheckResourceAttr("dt_project.test", "labels.%", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "dt_project" "test_labels" {
    display_name = "Project with labels"
    organization = "organizations/cvinmt9aq9sc738g6eog"
    location = {
        time_location = "Europe/Oslo"
    }
    labels = {
        "foo" = "bar"
    }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test_labels", "display_name", "Project with labels"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.%", "1"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.foo", "bar"),
				),
			},
			{
				Config: providerConfig + `
resource "dt_project" "test_labels" {
    display_name = "Project with labels"
    organization = "organizations/cvinmt9aq9sc738g6eog"
    location = {
        time_location = "Europe/Oslo"
    }
    labels = {
        "foo" = "baz"
        "new" = "label"
    }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test_labels", "display_name", "Project with labels"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.%", "2"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.foo", "baz"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.new", "label"),
				),
			},
			{
				Config: providerConfig + `
resource "dt_project" "test_labels" {
    display_name = "Project with labels"
    organization = "organizations/cvinmt9aq9sc738g6eog"
    location = {
        time_location = "Europe/Oslo"
    }
    labels = {}
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test_labels", "display_name", "Project with labels"),
					resource.TestCheckResourceAttr("dt_project.test_labels", "labels.%", "0"),
				),
			},
		},
	})
}
