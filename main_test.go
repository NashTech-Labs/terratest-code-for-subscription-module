package main

// Import key modules for this code ...

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// define the variables
var (
	globalBackendConf = make(map[string]interface{})

	globalEnvVars = make(map[string]string)

	subscription_name = "test-subscription"
)

// define the structure for function

type Subscription struct {
	ID string `json:"id"`

	DisplayName string `json:"displayName"`

	State string `json:"state"`

	TenantID string `json:"tenantId"`
}

// define the const variable
const (
	apiVersion = "2020-09-01"
)

// fucntion to setup the credentials

func setTerraformVariables() (map[string]string, error) {

	// Getting enVars from environment variables

	CLIENT_ID := os.Getenv("AZURE_CLIENT_ID")

	CLIENT_SECRET := os.Getenv("AZURE_CLIENT_SECRET")

	TENANT_ID := os.Getenv("AZURE_TENANT_ID")

	SUBSCRIPTION_ID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Creating globalEnVars for terraform call through Terratest

	if CLIENT_ID != "" {

		globalEnvVars["CLIENT_ID"] = CLIENT_ID

		globalEnvVars["CLIENT_SECRET"] = CLIENT_SECRET

		globalEnvVars["SUBSCRIPTION_ID"] = SUBSCRIPTION_ID

		globalEnvVars["TENANT_ID"] = TENANT_ID

	}

	return globalEnvVars, nil

}

// function to  run the test cases

func TestTerraform_azure_subscription(t *testing.T) {

	t.Parallel()

	setTerraformVariables()

	// Use Terratest to deploy the infrastructure

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{

		// Set the path to the Terraform code that will be tested.

		TerraformDir: "../module",

		// Variables to pass to our Terraform code using -var options

		Vars: map[string]interface{}{

			"subscription_name": subscription_name,
		},

		// globalvariables for user account

		EnvVars: globalEnvVars,

		// Backend values to set when initialziing Terraform

		BackendConfig: globalBackendConf,

		// Disable colors in Terraform commands so its easier to parse stdout/stderr

		NoColor: true,

		// Reconfigure is required if module deployment and go test pipelines are running in one stage

		Reconfigure: true,
	})

	// terraform destroy will run in the last
	defer terraform.Destroy(t, terraformOptions)

	// this will run the terraform init and terraform apply
	terraform.InitAndApply(t, terraformOptions)

	// this will run the terraform init and terraform plan

	// terraform.InitAndPlan(t, terraformOptions)

	expectedsubscription_id := terraform.Output(t, terraformOptions, "subscription_id")

	expectedsubscription_name := terraform.Output(t, terraformOptions, "subscription_name")

	expectedazuerm_subscription_tenant_id := terraform.Output(t, terraformOptions, "azuerm_subscription_tenant_id")

	fmt.Println("PRINTING THE RESOURCE PROPERTIES FROM OUTPUT FILE......................................")

	fmt.Printf("subscription_id : %s\n", expectedsubscription_id)

	fmt.Printf("subscription_name : %s\n", expectedsubscription_name)

	fmt.Printf("azuerm_subscription_tenant_id : %s\n", expectedazuerm_subscription_tenant_id)

	fmt.Println("PRINTING THE RESOURCE PROPERTIES FROM OUTPUT FILE HAS BEEN ENDED........................")

	accessToken, err := getAccessToken()

	if err != nil {

		fmt.Printf("Failed to get access token: %s\n", err.Error())

		return

	}

	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s?api-version=%s", expectedsubscription_id, apiVersion)

	fmt.Println("Here the URL for the subscription Module Rest API:", url)

	JSON, err := SubscriptionInfo(url, accessToken)

	if err != nil {

		fmt.Printf("Failed to get subscription details: %s\n", err.Error())

		return

	}

	actual_data, err := fetchSubscriptionInfo(JSON)

	if err != nil {

		log.Fatalf("failed to obtain a tha values from the RESTAPI: %v", err)

	}

	fmt.Println("\nSubscription Details from REST API - opened :-----------------")

	fmt.Printf("Display Name: %s\n", actual_data.DisplayName)
	fmt.Printf("Tenant ID: %s\n", actual_data.TenantID)

	fmt.Println("\nSubscription Details from REST API - closed :-----------------")

	// Test cases

	fmt.Println("Test cases are  running........")

	t.Run("Subscription_Name has been matched..", func(t *testing.T) {

		assert.Equal(t, expectedsubscription_name, actual_data.DisplayName)

	})

	t.Run("Tenant_Id has been matched..", func(t *testing.T) {

		assert.Equal(t, expectedazuerm_subscription_tenant_id, actual_data.TenantID)

	})

}

func getAccessToken() (string, error) {

	cmd := exec.Command("az", "account", "get-access-token", "--query", "accessToken", "--output", "tsv")

	output, err := cmd.Output()

	if err != nil {

		return "", err

	}

	return strings.TrimSpace(string(output)), nil

}

func SubscriptionInfo(url, accessToken string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {

		return nil, err

	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)

	if err != nil {

		return nil, err

	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return nil, err

	}

	return body, nil

}

func fetchSubscriptionInfo(subscriptionJSON []byte) (Subscription, error) {

	fmt.Println("#############################---Complete JSON Response from the RESTAPI---#######################################")

	fmt.Println(string(subscriptionJSON))

	fmt.Println("#############################---Complete JSON Response from the RESTAPI---#######################################")

	var subscription Subscription

	// unmarshal is use to convert the json ot xml data into go data structure

	err := json.Unmarshal(subscriptionJSON, &subscription)

	if err != nil {

		fmt.Printf("Failed to unmarshal JSON response: %s\n", err.Error())

		return Subscription{}, err

	}

	return subscription, err

}
