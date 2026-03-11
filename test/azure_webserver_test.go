package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "f7c39fda-34b4-47ff-82e7-efe69c03d3db"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "ruda0008",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// TEST 2  - Confirm NIC exists
	vmNics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, vmNics, nicName,
		"Expected NIC '%s' to be connected to VM '%s'", nicName, vmName)

	// TEST 3  - Confirm the VM is running
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, "Canonical", vmImage.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer)
	assert.Equal(t, "22_04-lts-gen2", vmImage.SKU)
}
