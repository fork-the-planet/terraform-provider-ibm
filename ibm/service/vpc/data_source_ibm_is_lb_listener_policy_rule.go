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

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISLBListenerPolicyRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenerPolicyRuleRead,

		Schema: map[string]*schema.Schema{
			isLBListenerPolicyRuleLBID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			isLBListenerPolicyRuleListenerID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The listener identifier.",
			},
			isLBListenerPolicyRulePolicyID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy identifier.",
			},
			isLBListenerPolicyRuleid: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule identifier.",
			},
			"condition": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The condition of the rule.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this rule was created.",
			},
			"field": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The field. This is applicable to `header`, `query`, and `body` rule types.If the rule type is `header`, this property is required.If the rule type is `query`, this is optional. If specified and the rule condition is not`matches_regex`, the value must be percent-encoded.If the rule type is `body`, this is optional.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule's canonical URL.",
			},
			"provisioning_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The provisioning status of this rule.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the rule.Body rules are applied to form-encoded request bodies using the `UTF-8` character set.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value to be matched for rule condition.If the rule type is `query` and the rule condition is not `matches_regex`, the value must be percent-encoded.",
			},
		},
	}
}

func dataSourceIBMIsLbListenerPolicyRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_lb_listener_policy_rule", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getLoadBalancerListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{}

	getLoadBalancerListenerPolicyRuleOptions.SetLoadBalancerID(d.Get(isLBListenerPolicyRuleLBID).(string))
	getLoadBalancerListenerPolicyRuleOptions.SetListenerID(d.Get(isLBListenerPolicyRuleListenerID).(string))
	getLoadBalancerListenerPolicyRuleOptions.SetPolicyID(d.Get(isLBListenerPolicyRulePolicyID).(string))
	getLoadBalancerListenerPolicyRuleOptions.SetID(d.Get(isLBListenerPolicyRuleid).(string))

	loadBalancerListenerPolicyRule, _, err := vpcClient.GetLoadBalancerListenerPolicyRuleWithContext(context, getLoadBalancerListenerPolicyRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetLoadBalancerListenerPolicyRuleWithContext failed: %s", err.Error()), "(Data) ibm_is_lb_listener_policy_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*loadBalancerListenerPolicyRule.ID)
	if err = d.Set("condition", loadBalancerListenerPolicyRule.Condition); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting condition: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-condition").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(loadBalancerListenerPolicyRule.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-created_at").GetDiag()
	}
	if !core.IsNil(loadBalancerListenerPolicyRule.Field) {
		if err = d.Set("field", loadBalancerListenerPolicyRule.Field); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting field: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-field").GetDiag()
		}
	}
	if err = d.Set("href", loadBalancerListenerPolicyRule.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-href").GetDiag()
	}
	if err = d.Set("provisioning_status", loadBalancerListenerPolicyRule.ProvisioningStatus); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting provisioning_status: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-provisioning_status").GetDiag()
	}
	if err = d.Set("type", loadBalancerListenerPolicyRule.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-type").GetDiag()
	}
	if err = d.Set("value", loadBalancerListenerPolicyRule.Value); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting value: %s", err), "(Data) ibm_is_lb_listener_policy_rule", "read", "set-value").GetDiag()
	}

	return nil
}
