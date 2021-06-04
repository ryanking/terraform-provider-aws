package aws

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/atest"
	awsprovider "github.com/terraform-providers/terraform-provider-aws/provider"
)

func testAccPreCheckAWSSSOAdminInstances(t *testing.T) {
	conn := atest.Provider.Meta().(*awsprovider.AWSClient).SSOAdminConn

	var instances []*ssoadmin.InstanceMetadata
	err := conn.ListInstancesPages(&ssoadmin.ListInstancesInput{}, func(page *ssoadmin.ListInstancesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		instances = append(instances, page.Instances...)

		return !lastPage
	})

	if atest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if len(instances) == 0 {
		t.Skip("skipping acceptance testing: No SSO Instance found.")
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func TestAccDataSourceAWSSSOAdminInstances_basic(t *testing.T) {
	dataSourceName := "data.aws_ssoadmin_instances.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { atest.PreCheck(t); testAccPreCheckAWSSSOAdminInstances(t) },
		ErrorCheck: atest.ErrorCheck(t, ssoadmin.EndpointsID),
		Providers:  atest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSSSOAdminInstancesConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "arns.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "identity_store_ids.#", "1"),
					atest.MatchAttrGlobalARNNoAccount(dataSourceName, "arns.0", "sso", regexp.MustCompile("instance/(sso)?ins-[a-zA-Z0-9-.]{16}")),
					resource.TestMatchResourceAttr(dataSourceName, "identity_store_ids.0", regexp.MustCompile("^[a-zA-Z0-9-]*")),
				),
			},
		},
	})
}

const testAccDataSourceAWSSSOAdminInstancesConfigBasic = `data "aws_ssoadmin_instances" "test" {}`
