// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISEndpointGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISEndpointGatewayRead,
		Importer:    &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			isVirtualEndpointGatewayName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint gateway name",
			},
			isVirtualEndpointGatewayResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway resource type",
			},
			isVirtualEndpointGatewayCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this Endpoint gateway",
			},
			isVirtualEndpointGatewayResourceGroupID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group id",
			},
			isVirtualEndpointGatewayCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway created date and time",
			},
			isVirtualEndpointGatewayServiceEndpoints: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The fully qualified domain names for the target service. A fully qualified domain name for the target service",
			},
			isVirtualEndpointGatewayHealthState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway health state",
			},
			isVirtualEndpointGatewayLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway lifecycle state",
			},
			isVirtualEndpointGatewayLifecycleReasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			isVirtualEndpointGatewaySecurityGroups: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Endpoint gateway securitygroups list",
			},
			isVirtualEndpointGatewayAllowDnsResolutionBinding: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to allow this endpoint gateway to participate in DNS resolution bindings with a VPC that has dns.enable_hub set to true.",
			},
			isVirtualEndpointGatewayIPs: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Endpoint gateway IPs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayIPsID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPs id",
						},
						isVirtualEndpointGatewayIPsName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPs name",
						},
						isVirtualEndpointGatewayIPsResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP resource type",
						},
						isVirtualEndpointGatewayIPsAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP Address",
						},
					},
				},
			},
			isVirtualEndpointGatewayTarget: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Endpoint gateway target",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayTargetName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target name",
						},
						isVirtualEndpointGatewayTargetCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target crn",
						},
						isVirtualEndpointGatewayTargetResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target resource type",
						},
					},
				},
			},
			isVirtualEndpointGatewayVpcID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC id",
			},
			isVirtualEndpointGatewayTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags for VPE",
			},
			isVirtualEndpointGatewayAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func dataSourceIBMISEndpointGatewayRead(
	context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_virtual_endpoint_gateway", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	name := d.Get(isVirtualEndpointGatewayName).(string)

	options := sess.NewListEndpointGatewaysOptions()
	options.Name = &name

	results, response, err := sess.ListEndpointGateways(options)
	if err != nil {
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListEndpointGateways failed: %s\n%s", err.Error(), response), "(Data) ibm_is_virtual_endpoint_gateway", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	allrecs := results.EndpointGateways

	if len(allrecs) == 0 {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListEndpointGateways failed: %s", err.Error()), "(Data) ibm_is_virtual_endpoint_gateway", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	result := allrecs[0]
	d.SetId(*result.ID)
	d.Set(isVirtualEndpointGatewayName, result.Name)
	d.Set(isVirtualEndpointGatewayAllowDnsResolutionBinding, result.AllowDnsResolutionBinding)
	d.Set(isVirtualEndpointGatewayCRN, result.CRN)
	d.Set(isVirtualEndpointGatewayHealthState, result.HealthState)
	d.Set(isVirtualEndpointGatewayCreatedAt, result.CreatedAt.String())
	d.Set(isVirtualEndpointGatewayLifecycleState, result.LifecycleState)
	if err := d.Set(isVirtualEndpointGatewayLifecycleReasons, resourceEGWFlattenLifecycleReasons(result.LifecycleReasons)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_virtual_endpoint_gateway", "read", "lifecycle_reasons-to-map").GetDiag()
	}
	d.Set(isVirtualEndpointGatewayResourceType, result.ResourceType)
	d.Set(isVirtualEndpointGatewayIPs, flattenIPs(result.Ips))
	d.Set(isVirtualEndpointGatewayResourceGroupID, result.ResourceGroup.ID)
	d.Set(isVirtualEndpointGatewayTarget, flattenEndpointGatewayTarget(
		result.Target.(*vpcv1.EndpointGatewayTarget)))
	d.Set(isVirtualEndpointGatewayVpcID, result.VPC.ID)
	if len(result.ServiceEndpoints) > 0 {
		d.Set(isVirtualEndpointGatewayServiceEndpoints, result.ServiceEndpoints)
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of VPE (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVirtualEndpointGatewayAccessTags, accesstags)
	if result.SecurityGroups != nil {
		d.Set(isVirtualEndpointGatewaySecurityGroups, flattenDataSourceSecurityGroups(result.SecurityGroups))
	}
	return nil
}
