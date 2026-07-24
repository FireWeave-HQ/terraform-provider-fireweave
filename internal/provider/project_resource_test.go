package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("set TF_ACC=1 to run acceptance tests")
	}
	if os.Getenv("FIREWEAVE_API_KEY") == "" {
		t.Skip("FIREWEAVE_API_KEY required for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "fireweave_project" "test" {
  name = "tf-acc-project"
  slug = "tf-acc-project"
  description = "acceptance test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fireweave_project.test", "name", "tf-acc-project"),
					resource.TestCheckResourceAttr("fireweave_project.test", "slug", "tf-acc-project"),
					resource.TestCheckResourceAttrSet("fireweave_project.test", "id"),
				),
			},
			{
				ResourceName:      "fireweave_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
