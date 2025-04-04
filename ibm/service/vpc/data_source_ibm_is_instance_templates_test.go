// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMISInstanceTemplates_dataBasic(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDConfig(vpcName, subnetName, sshKeyName, publicKey, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
				),
			},
		},
	})
}
func TestAccIBMISInstanceTemplatesDataSourceCluster(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDatasourceClusterConfig(vpcName, subnetName, sshKeyName, publicKey, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.boot_volume_attachment.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.cluster_network_attachments.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.cluster_network_attachments.0.cluster_network_interface.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.crn"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.enable_secure_boot"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.image"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.keys.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.network_attachments.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.primary_network_attachment.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.profile"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.resource_group"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.vpc"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.zone"),
				),
			},
		},
	})
}
func TestAccIBMISInstanceTemplates_dataconcom(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	ccmode := "sgx"
	esb := true
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDconcomConfig(vpcName, subnetName, sshKeyName, publicKey, templateName, ccmode, esb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.confidential_compute_mode"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.enable_secure_boot"),
				),
			},
		},
	})
}
func TestAccIBMISInstanceTemplates_dataVni(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDVniConfig(vpcName, subnetName, sshKeyName, publicKey, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
				),
			},
		},
	})
}
func TestAccIBMISInstanceTemplates_dataCatalog(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDCatalogConfig(vpcName, subnetName, sshKeyName, publicKey, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.catalog_offering.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.catalog_offering.0.version_crn"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.catalog_offering.0.plan_crn"),
					resource.TestCheckResourceAttr(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.image", ""),
				),
			},
		},
	})
}
func TestAccIBMISInstanceTemplates_dataReservedIp(t *testing.T) {
	randInt := acctest.RandIntRange(600, 700)
	publicKey := strings.TrimSpace(`
	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDQ+WiiUR1Jg3oGSmB/2//GJ3XnotriBiGN6t3iwGces6sUsvRkza1t0Mf05DKZxC/zp0WvDTvbit2gTkF9sD37OZSn5aCJk1F5URk/JNPmz25ZogkICFL4OUfhrE3mnyKio6Bk1JIEIypR5PRtGxY9vFDUfruADDLfRi+dGwHF6U9RpvrDRo3FNtI8T0GwvWwFE7bg63vLz65CjYY5XqH9z/YWz/asH6BKumkwiphLGhuGn03+DV6DkIZqr3Oh13UDjMnTdgv1y/Kou5UM3CK1dVsmLRXPEf2KUWUq1EwRfrJXkPOrBwn8to+Yydo57FgrRM9Qw8uzvKmnVxfKW6iG3oSGA0L6ROuCq1lq0MD8ySLd56+d1ftSDaUq+0/Yt9vK3olzVP0/iZobD7chbGqTLMCzL4/CaIUR/UmX08EA0Oh0DdyAdj3UUNETAj3W8gBrV6xLR7fZAJ8roX2BKb4K8Ed3YqzgiY0zgjqvpBYl9xZl0jgVX0qMFaEa6+CeGI8= root@ffd8363b1226
	`)
	vpcName := fmt.Sprintf("testvpc%d", randInt)
	subnetName := fmt.Sprintf("testsubnet%d", randInt)
	templateName := fmt.Sprintf("testtemplate%d", randInt)
	sshKeyName := fmt.Sprintf("testsshkey%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesReservedIpConfig(vpcName, subnetName, sshKeyName, publicKey, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_instance_templates.instance_templates_data", "templates.0.primary_network_interface.0.primary_ip.0.reserved_ip"),
				),
			},
		},
	})
}

func testAccCheckIBMISInstanceTemplatesDConfig(vpcName, subnetName, sshKeyName, publicKey, templateName string) string {
	return testAccCheckIBMISInstanceTemplateConfig(vpcName, subnetName, sshKeyName, publicKey, templateName) + fmt.Sprintf(`

		data "ibm_is_instance_templates" "instance_templates_data" {
			depends_on = [ibm_is_instance_template.instancetemplate1]
		}
	`)
}
func testAccCheckIBMISInstanceTemplatesDatasourceClusterConfig(vpcName, subnetName, sshKeyName, publicKey, templateName string) string {
	return fmt.Sprintf(`
		data "ibm_is_instance_templates" "instance_templates_data" {
		}
	`)
}
func testAccCheckIBMISInstanceTemplatesDconcomConfig(vpcName, subnetName, sshKeyName, publicKey, templateName, ccmode string, esb bool) string {
	return testAccCheckIBMISInstanceTemplateConComConfig(vpcName, subnetName, sshKeyName, publicKey, templateName, ccmode, esb) + fmt.Sprintf(`

		data "ibm_is_instance_templates" "instance_templates_data" {
			depends_on = [ibm_is_instance_template.instancetemplate1]
		}
	`)
}
func testAccCheckIBMISInstanceTemplatesDVniConfig(vpcName, subnetName, sshKeyName, publicKey, templateName string) string {
	return testAccCheckIBMISInstanceTemplateVniConfig(vpcName, subnetName, sshKeyName, publicKey, templateName) + fmt.Sprintf(`

		data "ibm_is_instance_templates" "instance_templates_data" {
			depends_on = [ibm_is_instance_template.instancetemplate1]
		}
	`)
}
func testAccCheckIBMISInstanceTemplatesDCatalogConfig(vpcName, subnetName, sshKeyName, publicKey, templateName string) string {
	return testAccCheckIBMISInstanceTemplateCatalogConfig(vpcName, subnetName, sshKeyName, publicKey, templateName) + fmt.Sprintf(`

		data "ibm_is_instance_templates" "instance_templates_data" {
		}
	`)
}
func testAccCheckIBMISInstanceTemplatesReservedIpConfig(vpcName, subnetName, sshKeyName, publicKey, templateName string) string {
	return testAccCheckIBMISInstanceTemplateRipConfig(vpcName, subnetName, sshKeyName, publicKey, templateName) + fmt.Sprintf(`

		data "ibm_is_instance_templates" "instance_templates_data" {
		}
	`)
}
