// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rum_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchrum"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccRUMAppMonitorDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var monitor cloudwatchrum.MetricDestinationSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_rum_app_monitor.test"
	resourceName := "aws_rum_app_monitor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudwatchrum.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchrum.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricsDestinationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsDestinationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsDestinationExists(ctx, resourceName, &monitor),
					resource.TestCheckResourceAttrPair(dataSourceName, "", resourceName, ""),
					resource.TestCheckResourceAttr(dataSourceName, "", ""),
				),
			},
		},
	})
}
