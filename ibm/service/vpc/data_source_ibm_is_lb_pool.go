// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISLBPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbPoolRead,

		Schema: map[string]*schema.Schema{
			"lb": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The pool identifier.",
			},
			"failsafe_policy": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A load balancer failsafe policy action:- `forward`: Forwards requests to the `target` pool.- `fail`: Rejects requests with an HTTP `503` status code.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"healthy_member_threshold_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The healthy member count at which the failsafe policy action will be triggered. At present, this is always `0`, but may be modifiable in the future.",
						},
						"target": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If `action` is `forward`, the target pool to forward to.If `action` is `fail`, this property will be absent.The targets supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this load balancer pool.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this load balancer pool. The name is unique across all pools for the load balancer.",
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The user-defined name for this load balancer pool.",
			},
			"algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load balancing algorithm.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this pool was created.",
			},
			"health_monitor": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The health monitor of this pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check interval in seconds. Interval must be greater than timeout value.",
						},
						"max_retries": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check max retries.",
						},
						"port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check port number. If specified, this overrides the ports specified in the server member resources.",
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The health check timeout in seconds.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol type of this load balancer pool health monitor.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the health monitor on which the unexpected property value was encountered.",
						},
						"url_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health check URL path. Applicable only if the health monitor `type` is `http` or`https`. This value must be in the format of an [origin-form request target](https://tools.ietf.org/html/rfc7230#section-5.3.1).",
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pool's canonical URL.",
			},
			"instance_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The instance group that is managing this pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this instance group.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this instance group.",
						},
					},
				},
			},
			"members": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The backend server members of the pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this load balancer pool member.",
						},
					},
				},
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used for this load balancer pool.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the pool on which the unexpected property value was encountered.",
			},
			"provisioning_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The provisioning status of this pool.",
			},
			"proxy_protocol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PROXY protocol setting for this pool:- `v1`: Enabled with version 1 (human-readable header format)- `v2`: Enabled with version 2 (binary header format)- `disabled`: DisabledSupported by load balancers in the `application` family (otherwise always `disabled`).",
			},
			"session_persistence": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The session persistence of this pool.The enumerated values for this property are expected to expand in the future. Whenprocessing this property, check for and log unknown values. Optionally haltprocessing and surface the error, or bypass the pool on which the unexpectedproperty value was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session persistence cookie name. Applicable only for type `app_cookie`. Names starting with `IBM` are not allowed.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session persistence type. The `http_cookie` and `app_cookie` types are applicable only to the `http` and `https` protocols.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsLbPoolRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_lb_pool", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	var loadBalancerPool *vpcv1.LoadBalancerPool

	if v, ok := d.GetOk("identifier"); ok {
		getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{}

		getLoadBalancerPoolOptions.SetLoadBalancerID(d.Get("lb").(string))
		getLoadBalancerPoolOptions.SetID(v.(string))

		loadBalancerPoolInfo, _, err := sess.GetLoadBalancerPoolWithContext(context, getLoadBalancerPoolOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetLoadBalancerPoolWithContext failed: %s", err.Error()), "(Data) ibm_is_lb_pool", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		loadBalancerPool = loadBalancerPoolInfo

	} else if v, ok := d.GetOk("name"); ok {
		listLoadBalancerPoolsOptions := &vpcv1.ListLoadBalancerPoolsOptions{}

		listLoadBalancerPoolsOptions.SetLoadBalancerID(d.Get("lb").(string))

		loadBalancerPoolCollection, _, err := sess.ListLoadBalancerPoolsWithContext(context, listLoadBalancerPoolsOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListLoadBalancerPoolsWithContext failed: %s", err.Error()), "(Data) ibm_is_lb_pool", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		name := v.(string)
		for _, data := range loadBalancerPoolCollection.Pools {
			if *data.Name == name {
				loadBalancerPool = &data
				break
			}
		}
		if loadBalancerPool == nil {
			log.Printf("[DEBUG] No LoadBalancerPool found with name (%s)", name)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("No LoadBalancerPool found with name: %s", name), "(Data) ibm_is_lb_pool", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	}

	d.SetId(*loadBalancerPool.ID)
	if err = d.Set("algorithm", loadBalancerPool.Algorithm); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting algorithm: %s", err), "(Data) ibm_is_lb_pool", "read", "set-algorithm").GetDiag()
	}
	failsafePolicy := []map[string]interface{}{}
	if loadBalancerPool.FailsafePolicy != nil {
		modelMap, err := dataSourceIBMIsLbPoolLoadBalancerPoolFailsafePolicyToMap(loadBalancerPool.FailsafePolicy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_lb_pool", "read", "failsafe_policy-to-map").GetDiag()
		}
		failsafePolicy = append(failsafePolicy, modelMap)
	}
	if err = d.Set("failsafe_policy", failsafePolicy); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting failsafe_policy: %s", err), "(Data) ibm_is_lb_pool", "read", "set-failsafe_policy").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(loadBalancerPool.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_lb_pool", "read", "set-created_at").GetDiag()
	}

	if loadBalancerPool.HealthMonitor != nil {
		poolHealthMonitor := loadBalancerPool.HealthMonitor.(*vpcv1.LoadBalancerPoolHealthMonitor)
		err = d.Set("health_monitor", dataSourceLoadBalancerPoolFlattenHealthMonitor(*poolHealthMonitor))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting health_monitor: %s", err), "(Data) ibm_is_lb_pool", "read", "set-health_monitor").GetDiag()
		}
	}
	if err = d.Set("href", loadBalancerPool.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_lb_pool", "read", "set-href").GetDiag()
	}

	if loadBalancerPool.InstanceGroup != nil {
		err = d.Set("instance_group", dataSourceLoadBalancerPoolFlattenInstanceGroup(*loadBalancerPool.InstanceGroup))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting instance_group: %s", err), "(Data) ibm_is_lb_pool", "read", "set-instance_group").GetDiag()
		}
	}

	if loadBalancerPool.Members != nil {
		err = d.Set("members", dataSourceLoadBalancerPoolFlattenMembers(loadBalancerPool.Members))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting members: %s", err), "(Data) ibm_is_lb_pool", "read", "set-members").GetDiag()
		}
	}

	if err = d.Set("identifier", loadBalancerPool.ID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting identifier: %s", err), "(Data) ibm_is_lb_pool", "read", "set-identifier").GetDiag()
	}

	if err = d.Set("name", loadBalancerPool.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_lb_pool", "read", "set-name").GetDiag()
	}
	if err = d.Set("protocol", loadBalancerPool.Protocol); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting protocol: %s", err), "(Data) ibm_is_lb_pool", "read", "set-protocol").GetDiag()
	}
	if err = d.Set("provisioning_status", loadBalancerPool.ProvisioningStatus); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting provisioning_status: %s", err), "(Data) ibm_is_lb_pool", "read", "set-provisioning_status").GetDiag()
	}

	if err = d.Set("proxy_protocol", loadBalancerPool.ProxyProtocol); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting proxy_protocol: %s", err), "(Data) ibm_is_lb_pool", "read", "set-proxy_protocol").GetDiag()
	}

	if loadBalancerPool.SessionPersistence != nil {
		err = d.Set("session_persistence", dataSourceLoadBalancerPoolFlattenSessionPersistence(*loadBalancerPool.SessionPersistence))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting session_persistence: %s", err), "(Data) ibm_is_lb_pool", "read", "set-session_persistence").GetDiag()
		}
	}

	return nil
}

func dataSourceLoadBalancerPoolFlattenHealthMonitor(result vpcv1.LoadBalancerPoolHealthMonitor) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerPoolHealthMonitorToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerPoolHealthMonitorToMap(healthMonitorItem vpcv1.LoadBalancerPoolHealthMonitor) (healthMonitorMap map[string]interface{}) {
	healthMonitorMap = map[string]interface{}{}

	if healthMonitorItem.Delay != nil {
		healthMonitorMap["delay"] = healthMonitorItem.Delay
	}
	if healthMonitorItem.MaxRetries != nil {
		healthMonitorMap["max_retries"] = healthMonitorItem.MaxRetries
	}
	if healthMonitorItem.Port != nil {
		healthMonitorMap["port"] = healthMonitorItem.Port
	}
	if healthMonitorItem.Timeout != nil {
		healthMonitorMap["timeout"] = healthMonitorItem.Timeout
	}
	if healthMonitorItem.Type != nil {
		healthMonitorMap["type"] = healthMonitorItem.Type
	}
	if healthMonitorItem.URLPath != nil {
		healthMonitorMap["url_path"] = healthMonitorItem.URLPath
	}

	return healthMonitorMap
}

func dataSourceLoadBalancerPoolFlattenInstanceGroup(result vpcv1.InstanceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerPoolInstanceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerPoolInstanceGroupToMap(instanceGroupItem vpcv1.InstanceGroupReference) (instanceGroupMap map[string]interface{}) {
	instanceGroupMap = map[string]interface{}{}

	if instanceGroupItem.CRN != nil {
		instanceGroupMap["crn"] = instanceGroupItem.CRN
	}
	if instanceGroupItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerPoolInstanceGroupDeletedToMap(*instanceGroupItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instanceGroupMap["deleted"] = deletedList
	}
	if instanceGroupItem.Href != nil {
		instanceGroupMap["href"] = instanceGroupItem.Href
	}
	if instanceGroupItem.ID != nil {
		instanceGroupMap["id"] = instanceGroupItem.ID
	}
	if instanceGroupItem.Name != nil {
		instanceGroupMap["name"] = instanceGroupItem.Name
	}

	return instanceGroupMap
}

func dataSourceLoadBalancerPoolInstanceGroupDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerPoolFlattenMembers(result []vpcv1.LoadBalancerPoolMemberReference) (members []map[string]interface{}) {
	for _, membersItem := range result {
		members = append(members, dataSourceLoadBalancerPoolMembersToMap(membersItem))
	}

	return members
}

func dataSourceLoadBalancerPoolMembersToMap(membersItem vpcv1.LoadBalancerPoolMemberReference) (membersMap map[string]interface{}) {
	membersMap = map[string]interface{}{}

	if membersItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerPoolMembersDeletedToMap(*membersItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		membersMap["deleted"] = deletedList
	}
	if membersItem.Href != nil {
		membersMap["href"] = membersItem.Href
	}
	if membersItem.ID != nil {
		membersMap["id"] = membersItem.ID
	}

	return membersMap
}

func dataSourceLoadBalancerPoolMembersDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerPoolFlattenSessionPersistence(result vpcv1.LoadBalancerPoolSessionPersistence) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerPoolSessionPersistenceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerPoolSessionPersistenceToMap(sessionPersistenceItem vpcv1.LoadBalancerPoolSessionPersistence) (sessionPersistenceMap map[string]interface{}) {
	sessionPersistenceMap = map[string]interface{}{}

	if sessionPersistenceItem.CookieName != nil {
		sessionPersistenceMap["cookie_name"] = sessionPersistenceItem.CookieName
	}
	if sessionPersistenceItem.Type != nil {
		sessionPersistenceMap["type"] = sessionPersistenceItem.Type
	}

	return sessionPersistenceMap
}

func dataSourceIBMIsLbPoolLoadBalancerPoolFailsafePolicyToMap(model *vpcv1.LoadBalancerPoolFailsafePolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["action"] = *model.Action
	modelMap["healthy_member_threshold_count"] = flex.IntValue(model.HealthyMemberThresholdCount)
	if model.Target != nil {
		targetMap, err := dataSourceIBMIsLbPoolLoadBalancerPoolReferenceToMap(model.Target)
		if err != nil {
			return modelMap, err
		}
		modelMap["target"] = []map[string]interface{}{targetMap}
	}
	return modelMap, nil
}

func dataSourceIBMIsLbPoolLoadBalancerPoolReferenceToMap(model *vpcv1.LoadBalancerPoolReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsLbPoolDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func dataSourceIBMIsLbPoolDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}
