package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOTCBMSTagsV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_compute_bms_tags_v2.tags_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccBmsFlavorPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
