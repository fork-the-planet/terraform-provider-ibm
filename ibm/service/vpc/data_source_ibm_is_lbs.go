// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	loadBalancers      = "load_balancers"
	CRN                = "crn"
	CreatedAt          = "created_at"
	isLbProfile        = "profile"
	ProvisioningStatus = "provisioning_status"
	ID                 = "id"
)

func DataSourceIBMISLBS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISLBSRead,
		Schema: map[string]*schema.Schema{
			loadBalancers: {
				Type:        schema.TypeList,
				Description: "Collection of load balancers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBAccessMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access mode of this load balancer",
						},
						isAttachedLoadBalancerPoolMembers: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The load balancer pool members attached to this load balancer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this load balancer pool member.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool member.",
									},
								},
							},
						},
						ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						CRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancer's CRN",
						},
						CreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this pool was created.",
						},
						"failsafe_policy_actions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The supported `failsafe_policy.action` values for this load balancer's pools.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"dns": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DNS configuration for this load balancer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this DNS instancer",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier of the DNS zone.",
									},
								},
							},
						},
						ProvisioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this load balancer",
						},
						isLBName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer name",
						},
						isLBAvailability: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability of this load balancer",
						},
						isLBInstanceGroupsSupported: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this load balancer supports instance groups.",
						},
						isLBSourceIPPersistenceSupported: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this load balancer supports source IP session persistence.",
						},
						isLBUdpSupported: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this load balancer supports UDP.",
						},
						isLBRouteMode: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Load Balancer route mode",
						},

						isLBType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer type",
						},

						isLBStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer status",
						},

						isLBOperatingStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer operating status",
						},

						isLBPublicIPs: {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Load Balancer Public IPs",
						},

						isLBPrivateIPs: {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Load Balancer private IPs",
						},
						isLBPrivateIPDetail: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The private IP addresses assigned to this load balancer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isLBPrivateIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isLBPrivateIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isLBPrivateIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isLBPrivateIpId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isLBPrivateIpResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isLBSubnets: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer subnets list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subnet's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer subnet",
									},
									name: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this load balancer subnet",
									},
									CRN: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet",
									},
								},
							},
						},

						isLBTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "Tags associated to Load Balancer",
						},

						isLBAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},

						isLBResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer Resource group",
						},

						isLBHostName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load Balancer Host Name",
						},

						isLBListeners: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer Listeners list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The listener's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer listener",
									},
								},
							},
						},
						isLbProfile: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The profile to use for this load balancer",
						},

						isLBPools: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer Pools list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The pool's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool",
									},
									name: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this load balancer pool",
									},
								},
							},
						},

						flex.ResourceControllerURL: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
						},

						flex.ResourceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the resource",
						},

						flex.ResourceGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group name in which resource is provisioned",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISLBSRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	err := getLbs(context, d, meta)
	if err != nil {
		return err
	}
	return nil
}

func getLbs(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_lbs", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	start := ""
	allrecs := []vpcv1.LoadBalancer{}
	for {
		listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listLoadBalancersOptions.Start = &start
		}
		lbs, _, err := sess.ListLoadBalancersWithContext(context, listLoadBalancersOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListLoadBalancersWithContext failed %s", err), "(Data) ibm_is_lbs", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		start = flex.GetNext(lbs.Next)
		allrecs = append(allrecs, lbs.LoadBalancers...)
		if start == "" {
			break
		}
	}

	lbList := make([]map[string]interface{}, 0)

	for _, lb := range allrecs {
		lbInfo := make(map[string]interface{})
		//	log.Printf("******* lb ******** : (%+v)", lb)
		lbInfo[ID] = *lb.ID
		if lb.Availability != nil {
			lbInfo[isLBAvailability] = *lb.Availability
		}
		if lb.AccessMode != nil {
			lbInfo[isLBAccessMode] = *lb.AccessMode
		}
		if lb.AttachedLoadBalancerPoolMembers != nil {
			lbInfo[isAttachedLoadBalancerPoolMembers] = dataSourceAttachedLoadBalancerPoolFlattenMembers(lb.AttachedLoadBalancerPoolMembers)
		}
		if lb.InstanceGroupsSupported != nil {
			lbInfo[isLBInstanceGroupsSupported] = *lb.InstanceGroupsSupported
		}
		if lb.SourceIPSessionPersistenceSupported != nil {
			lbInfo[isLBSourceIPPersistenceSupported] = *lb.SourceIPSessionPersistenceSupported
		}
		lbInfo[isLBName] = *lb.Name
		dnsList := make([]map[string]interface{}, 0)
		if lb.Dns != nil {
			dns := map[string]interface{}{}
			dns["instance_crn"] = lb.Dns.Instance.CRN
			dns["zone_id"] = lb.Dns.Zone.ID
			dnsList = append(dnsList, dns)
			lbInfo["dns"] = dnsList
		}
		if lb.RouteMode != nil {
			lbInfo[isLBRouteMode] = *lb.RouteMode
		}
		if lb.UDPSupported != nil {
			lbInfo[isLBUdpSupported] = *lb.UDPSupported
		}
		lbInfo[CRN] = *lb.CRN
		lbInfo[ProvisioningStatus] = *lb.ProvisioningStatus

		lbInfo[CreatedAt] = lb.CreatedAt.String()
		if lb.IsPublic != nil && *lb.IsPublic {
			lbInfo[isLBType] = "public"
		} else if lb.IsPrivatePath != nil && *lb.IsPrivatePath {
			lbInfo[isLBType] = "private_path"
		} else {
			lbInfo[isLBType] = "private"
		}
		lbInfo[isLBStatus] = *lb.ProvisioningStatus
		lbInfo[isLBOperatingStatus] = *lb.OperatingStatus
		publicIpList := make([]string, 0)
		if lb.PublicIps != nil {
			for _, ip := range lb.PublicIps {
				if ip.Address != nil {
					pubip := *ip.Address
					publicIpList = append(publicIpList, pubip)
				}
			}
		}

		lbInfo[isLBPublicIPs] = publicIpList
		privateIpList := make([]string, 0)
		privateIpDetailList := make([]map[string]interface{}, 0)
		if lb.PrivateIps != nil {
			for _, ip := range lb.PrivateIps {
				if ip.Address != nil {
					prip := *ip.Address
					privateIpList = append(privateIpList, prip)
				}
				currentPriIp := map[string]interface{}{}
				if ip.Address != nil {
					currentPriIp[isLBPrivateIpAddress] = ip.Address
				}
				if ip.Href != nil {
					currentPriIp[isLBPrivateIpHref] = ip.Href
				}
				if ip.Name != nil {
					currentPriIp[isLBPrivateIpName] = ip.Name
				}
				if ip.ID != nil {
					currentPriIp[isLBPrivateIpId] = ip.ID
				}
				if ip.ResourceType != nil {
					currentPriIp[isLBPrivateIpResourceType] = ip.ResourceType
				}
				privateIpDetailList = append(privateIpDetailList, currentPriIp)
			}
		}
		lbInfo[isLBPrivateIPs] = privateIpList
		lbInfo[isLBPrivateIPDetail] = privateIpDetailList
		//log.Printf("*******isLBPrivateIPs %+v", lbInfo[isLBPrivateIPs])

		if lb.Subnets != nil {
			subnetList := make([]map[string]interface{}, 0)
			for _, subnet := range lb.Subnets {
				//	log.Printf("*******subnet %+v", subnet)
				sub := make(map[string]interface{})
				sub[ID] = *subnet.ID
				sub[href] = *subnet.Href
				if subnet.CRN != nil {
					sub[CRN] = *subnet.CRN
				}
				if subnet.Name != nil && *subnet.Name != "" {
					sub[name] = *subnet.Name
				}
				subnetList = append(subnetList, sub)

			}
			lbInfo[isLBSubnets] = subnetList
			//	log.Printf("*******isLBSubnets %+v", lbInfo[isLBSubnets])

		}
		if lb.Listeners != nil {
			listenerList := make([]map[string]interface{}, 0)
			for _, listener := range lb.Listeners {
				lis := make(map[string]interface{})
				lis[ID] = *listener.ID
				lis[href] = *listener.Href
				listenerList = append(listenerList, lis)
			}
			lbInfo[isLBListeners] = listenerList
		}
		//log.Printf("*******isLBListeners %+v", lbInfo[isLBListeners])
		if lb.Pools != nil {
			poolList := make([]map[string]interface{}, 0)

			for _, p := range lb.Pools {
				pool := make(map[string]interface{})
				pool[name] = *p.Name
				pool[ID] = *p.ID
				pool[href] = *p.Href
				poolList = append(poolList, pool)

			}
			lbInfo[isLBPools] = poolList
		}
		lbInfo["failsafe_policy_actions"] = lb.FailsafePolicyActions
		if lb.Profile != nil {
			lbProfile := make(map[string]interface{})
			lbProfile[name] = *lb.Profile.Name
			lbProfile[href] = *lb.Profile.Href
			lbProfile[family] = *lb.Profile.Family
			lbInfo[isLbProfile] = lbProfile
		}
		lbInfo[isLBResourceGroup] = *lb.ResourceGroup.ID
		lbInfo[isLBHostName] = *lb.Hostname
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *lb.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
		}
		lbInfo[isLBTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *lb.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource Load Balancer (%s) access tags: %s", d.Id(), err)
		}
		lbInfo[isLBAccessTags] = accesstags

		controller, err := flex.GetBaseController(meta)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBaseController failed %s", err), "(Data) ibm_is_lbs", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		lbInfo[flex.ResourceControllerURL] = controller + "/vpc-ext/network/loadBalancers"
		lbInfo[flex.ResourceName] = *lb.Name
		//log.Printf("*******lbInfo %+v", lbInfo)

		if lb.ResourceGroup != nil {
			lbInfo[flex.ResourceGroupName] = *lb.ResourceGroup.ID
		}
		lbList = append(lbList, lbInfo)
		//	log.Printf("*******lbList %+v", lbList)

	}
	//log.Printf("*******lbList %+v", lbList)
	d.SetId(dataSourceIBMISLBsID(d))
	if err = d.Set("load_balancers", lbList); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting load_balancers %s", err), "(Data) ibm_is_lbs", "read", "load_balancers-set").GetDiag()
	}

	return nil
}

// dataSourceIBMISLBsID returns a reasonable ID for a transit gateways list.
func dataSourceIBMISLBsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
