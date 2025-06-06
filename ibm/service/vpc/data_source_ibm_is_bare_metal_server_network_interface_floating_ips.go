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

const (
	isBareMetalServerNetworkInterface = "network_interface"
	floatingIPId                      = "id"
)

func DataSourceIBMIsBareMetalServerNetworkInterfaceFloatingIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPsRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},
			isBareMetalServerNetworkInterface: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network interface identifier of bare metal server",
			},

			//floating ip properties
			isBareMetalServerNicFloatingIPs: {
				Type:        schema.TypeList,
				Description: "The floating IPs associated with this network interface.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						floatingIPName: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the floating IP",
						},
						floatingIPId: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the floating IP",
						},

						floatingIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating IP address",
						},

						floatingIPStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating IP status",
						},

						floatingIPZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone name",
						},

						floatingIPTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target info",
						},

						floatingIPCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating IP crn",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	nicID := d.Get(isBareMetalServerNetworkInterface).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_network_interface_floating_ips", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	allFloatingIPs := []vpcv1.FloatingIP{}
	options := &vpcv1.ListBareMetalServerNetworkInterfaceFloatingIpsOptions{
		BareMetalServerID:  &bareMetalServerID,
		NetworkInterfaceID: &nicID,
	}

	fips, _, err := sess.ListBareMetalServerNetworkInterfaceFloatingIpsWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListBareMetalServerNetworkInterfaceFloatingIpsWithContext failed: %s", err.Error()), "(Data) ibm_is_bare_metal_server_network_interface_floating_ips", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	allFloatingIPs = append(allFloatingIPs, fips.FloatingIps...)
	fipInfo := make([]map[string]interface{}, 0)
	for _, ip := range allFloatingIPs {
		l := map[string]interface{}{}

		l[floatingIPName] = *ip.Name
		l[floatingIPAddress] = *ip.Address
		l[floatingIPStatus] = *ip.Status
		l[floatingIPZone] = *ip.Zone.Name

		l[floatingIPCRN] = *ip.CRN

		target, ok := ip.Target.(*vpcv1.FloatingIPTarget)
		if ok {
			l[floatingIPTarget] = target.ID
		}

		l[floatingIPId] = *ip.ID

		fipInfo = append(fipInfo, l)
	}
	d.SetId(dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPsID(d))

	if err = d.Set(isBareMetalServerNicFloatingIPs, fipInfo); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting floating_ips: %s", err), "(Data) ibm_is_bare_metal_server_network_interface_floating_ips", "read", "set-floating_ips").GetDiag()
	}
	return nil
}

// dataSourceIBMISBMSProfilesID returns a reasonable ID for a BMS network interface floating ip list.
func dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
