package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIP_basic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2FloatingIPExists("opentelekomcloud_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_timeout(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIP_timeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2FloatingIPExists("opentelekomcloud_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("FloatingIP still exists")
		}
	}

	return nil
}

const testAccNetworkingV2FloatingIP_basic = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}
`

var testAccNetworkingV2FloatingIP_fixedip_bind = fmt.Sprintf(`
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "%s"
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = opentelekomcloud_networking_subnet_v2.subnet_1.network_id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.20"
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  pool = "%s"
  port_id = opentelekomcloud_networking_port_v2.port_1.id
  fixed_ip = opentelekomcloud_networking_port_v2.port_1.fixed_ip.1.ip_address
}
`, env.OS_EXTGW_ID, env.OS_POOL_NAME)

const testAccNetworkingV2FloatingIP_timeout = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
