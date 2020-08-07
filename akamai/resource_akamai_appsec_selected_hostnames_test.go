package akamai

import (
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAkamaiSelectedHostnames_basic(t *testing.T) {

	dataSourceName := "akamai_appsec_selected_hostnames.appsecselectedhostnames"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiSelectedHostnamesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSelectedHostnamesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiSelectedHostnamesExists,
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSelectedHostnamesConfig() string {
	return `
provider "akamai" {
  appsec_section = "global"
}

data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.configid
}

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    configid = data.akamai_appsec_configuration.appsecconfigedge.configid
    version = data.akamai_appsec_configuration.appsecconfigedge.latestversion 
    hostnames = ["*.example.net","example.com","m.example.com"]  
}


`
}

func testAccCheckAkamaiSelectedHostnamesExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_selected_hostnames" {
			continue
		}
		//rname := rs.Primary.ID
		configid, _ := strconv.Atoi(rs.Primary.Attributes["configid"])
		version, _ := strconv.Atoi(rs.Primary.Attributes["version"])
		ccresp := appsec.NewSelectedHostnamesResponse()
		err = ccresp.GetSelectedHostnames(configid, version, "TEST")

		if err != nil {
			return err
		}
	}
	return nil
}