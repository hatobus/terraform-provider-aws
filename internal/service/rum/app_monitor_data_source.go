// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rum

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchrum"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_rum_app_monitor")
func DataSourceAppMonitor() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAppMonitorRead,

		Schema: map[string]*schema.Schema{
			"allow_cookies": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_xray": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"excluded_pages": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"guest_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"included_pages": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"session_sample_rate": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"telemetries": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

const (
	DSAppMonitor             = "CloudWatch RUM Data Source"
	ListAppMonitorMaxResults = 100
)

func dataSourceAppMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RUMConn(ctx)
	name := d.Get("app_monitor_name").(string)

	monitor, err := findAppMonitorByName(ctx, conn, name)
	if err != nil {
		return append(diags, create.DiagError(names.RUM, create.ErrActionReading, DSAppMonitor, name, err)...)
	}

	d.SetId(aws.StringValue(monitor.Id))

	d.Set("allow_cookies", monitor.AppMonitorConfiguration.AllowCookies)
	d.Set("enable_xray", monitor.AppMonitorConfiguration.EnableXRay)
	d.Set("guest_role_arn", monitor.AppMonitorConfiguration.GuestRoleArn)
	d.Set("identity_pool_id", monitor.AppMonitorConfiguration.IdentityPoolId)
	d.Set("session_sample_rate", monitor.AppMonitorConfiguration.SessionSampleRate)

	d.Set("domain", monitor.Domain)

	setTagsOut(ctx, monitor.Tags)

	if err := d.Set("excluded_pages", flex.FlattenStringList(monitor.AppMonitorConfiguration.ExcludedPages)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting excluded_pages error: %s", err)
	}

	if err := d.Set("included_pages", flex.FlattenStringList(monitor.AppMonitorConfiguration.IncludedPages)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting included_pages error: %s", err)
	}

	if err := d.Set("telemetries", flex.FlattenStringList(monitor.AppMonitorConfiguration.Telemetries)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting telemetries error: %s", err)
	}

	return diags
}

func findAppMonitorByName(ctx context.Context, conn *cloudwatchrum.CloudWatchRUM, name string) (*cloudwatchrum.AppMonitor, error) {
	var monitorName string
	input := &cloudwatchrum.ListAppMonitorsInput{
		MaxResults: aws.Int64(ListAppMonitorMaxResults),
	}

	if err := conn.ListAppMonitorsPagesWithContext(ctx, input, func(page *cloudwatchrum.ListAppMonitorsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, r := range page.AppMonitorSummaries {
			if r == nil {
				continue
			}
			if aws.StringValue(r.Name) == name {
				monitorName = aws.StringValue(r.Name)
				return false
			}
		}
		return !lastPage
	}); err != nil {
		return nil, err
	}

	monitor, err := conn.GetAppMonitorWithContext(ctx, &cloudwatchrum.GetAppMonitorInput{
		Name: aws.String(monitorName),
	})
	if err != nil {
		return nil, err
	}

	if monitor.AppMonitor == nil || monitorName == "" {
		return nil, fmt.Errorf("no app monitor found with name %q", monitorName)
	}

	return monitor.AppMonitor, nil
}
