##Overview
The primary purpose is to implement support for some services which the official Azure 
client does not support. Currently it adds support for:

- [RateCard](https://msdn.microsoft.com/en-us/library/azure/mt219005.aspx)


##Examples
####Pricing
This example lists all Linux Pay-as-you-go VM Pricing and Sizes in USD for all VMs in the US. Example output:

```
{Standard_F16 1.025 northcentralus 16 32 256}.
{Basic_A0 0.018 westus 1 0.75 20}.
{Basic_A1 0.044 westus2 1 1.75 40}.
{Standard_F16 0.891 southcentralus 16 32 256}.
{Basic_A3 0.176 eastus 4 7 120}.
{Standard_D2_v2 0.146 northcentralus 2 7 100}.
...
```
Fields are: Size, Pricing (per hour), Region, CPU, RAM, Disk.