package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS, but should be able to set connection_limit to a number, but can't yet.
func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccLBV2ListenerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("huaweicloud_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"huaweicloud_lb_listener_v2.listener_1", "connection_limit", "-1"),
				),
			},
			resource.TestStep{
				Config: TestAccLBV2ListenerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"huaweicloud_lb_listener_v2.listener_1", "connection_limit", "100"),
				),
			},
		},
	})
}

func testAccCheckLBV2ListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_listener_v2" {
			continue
		}

		_, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2ListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

const TestAccLBV2ListenerConfig_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${huaweicloud_lb_loadbalancer_v2.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`

const TestAccLBV2ListenerConfig_update = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  admin_state_up = "true"
  loadbalancer_id = "${huaweicloud_lb_loadbalancer_v2.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`
