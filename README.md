## Terratest to test the terraform  azure subscription module using Go Lang

### Basically terratest is use to test the terraform module using GOlang. You can use the terratest for kubernetes or many other but in our case we are using this for terraform module. This template can use to test the Azure subcription module using Go Lang. You need to export the credentials of your Azure account or az login with the Master service principal so that you have access to create the Azure subcription and then run the terratest against the Azure subscription module.

---
You can follow the below instructions to run the terratest code and test your subscription module.

Step 1:- You must have the terraform module to run the terratest because you need to pass the module path to the terratest so that terratest can run the terraform code then perform the test on it.

Step 2:- Now you need to configure the azure credentials in your system with export process or az login credentials but with you master service pricipal credentials.

For Windows:-
 

        $Env:ARM_CLIENT_ID=""

        $Env:ARM_CLIENT_SECRET=""

        $Env:ARM_TENANT_ID=""

        $Env:ARM_SUBSCRIPTION_ID=""

For Linux:- 

        export AZURE_TENANT_ID=""

        export AZURE_CLIENT_ID=""

        export AZURE_CLIENT_SECRET=""

        export AZURE_SUBSCRIPTION_ID=""

Step 3:- or you can go with the az login.

    az login --service-principal -u < > -p < > --tenant < >

Step 4:-  Run the below command to run your test case.

    go mod init <name your module>

    go mod tidy

    go test -v  

