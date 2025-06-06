// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isInstanceNICReservedIPLimit  = "limit"
	isInstanceNICReservedIPSort   = "sort"
	isInstanceNICReservedIPs      = "reserved_ips"
	isInstanceNICReservedIPsCount = "total_count"
)

func DataSourceIBMISInstanceNICReservedIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISInstanceNICReservedIPsRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
			*/
			isInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance identifier.",
			},
			isInstanceNICID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance network interface identifier.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isInstanceNICReservedIPs: {
				Type:        schema.TypeList,
				Description: "Collection of reserved IPs in this subnet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceNICReservedIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						isInstanceNICReservedIPAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If reserved ip shall be deleted automatically",
						},
						isInstanceNICReservedIPCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the reserved IP was created.",
						},
						isInstanceNICReservedIPhref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						isInstanceNICReservedIPID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP",
						},
						isInstanceNICReservedIPName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this reserved IP.",
						},
						isInstanceNICReservedIPOwner: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
						},
						isInstanceNICReservedIPType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						isInstanceNICReservedIPTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserved IP target id",
						},
					},
				},
			},
			isInstanceNICReservedIPsCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages",
			},
		},
	}
}

func dataSourceIBMISInstanceNICReservedIPsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_network_interface_reserved_ips", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceID := d.Get(isInstanceID).(string)
	nicID := d.Get(isInstanceNICID).(string)

	// Flatten all the reserved IPs
	start := ""
	allrecs := []vpcv1.ReservedIP{}
	for {
		options := &vpcv1.ListInstanceNetworkInterfaceIpsOptions{
			InstanceID:         &instanceID,
			NetworkInterfaceID: &nicID,
		}
		if start != "" {
			options.Start = &start
		}

		result, response, err := sess.ListInstanceNetworkInterfaceIpsWithContext(context, options)
		if err != nil || response == nil || result == nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListInstanceNetworkInterfaceIpsWithContext failed: %s", err.Error()), "(Data) ibm_is_instance_network_interface_reserved_ips", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		start = flex.GetNext(result.Next)
		allrecs = append(allrecs, result.Ips...)
		if start == "" {
			break
		}
	}
	// Now store all the reserved IP info with their response tags
	reservedIPs := []map[string]interface{}{}
	for _, data := range allrecs {
		ipsOutput := map[string]interface{}{}
		ipsOutput[isInstanceNICReservedIPAddress] = *data.Address
		ipsOutput[isInstanceNICReservedIPAutoDelete] = *data.AutoDelete
		ipsOutput[isInstanceNICReservedIPCreatedAt] = (*data.CreatedAt).String()
		ipsOutput[isInstanceNICReservedIPhref] = *data.Href
		ipsOutput[isInstanceNICReservedIPID] = *data.ID
		ipsOutput[isInstanceNICReservedIPName] = *data.Name
		ipsOutput[isInstanceNICReservedIPOwner] = *data.Owner
		ipsOutput[isInstanceNICReservedIPType] = *data.ResourceType
		target, ok := data.Target.(*vpcv1.ReservedIPTarget)
		if ok {
			ipsOutput[isReservedIPTarget] = target.ID
		}
		reservedIPs = append(reservedIPs, ipsOutput)
	}

	d.SetId(time.Now().UTC().String()) // This is not any reserved ip or instance id but state id
	if err = d.Set(isInstanceNICReservedIPs, reservedIPs); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting reserved_ips: %s", err), "(Data) ibm_is_instance_network_interface_reserved_ips", "read", "set-reserved_ips").GetDiag()
	}
	if err = d.Set(isInstanceNICReservedIPsCount, len(reservedIPs)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_is_instance_network_interface_reserved_ips", "read", "set-total_count").GetDiag()
	}
	if err = d.Set(isInstanceID, instanceID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting instance: %s", err), "(Data) ibm_is_instance_network_interface_reserved_ips", "read", "set-instance").GetDiag()
	}
	if err = d.Set(isInstanceNICID, nicID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting network_interface: %s", err), "(Data) ibm_is_instance_network_interface_reserved_ips", "read", "set-network_interface").GetDiag()
	}
	return nil
}
