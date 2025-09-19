az network public-ip create -g evanbaker-hyperscale -n bastion --sku Standard --location eastus2
az network bastion create -g evanbaker-hyperscale -n bastion --public-ip-address bastion --vnet-name k3a-vnet --location eastus2 --sku Standard
az network public-ip create -g evanbaker-hyperscale -n k3a-natgw --sku Standard --location eastus2;
az network nat gateway create -g evanbaker-hyperscale -n k3a-natgw --location eastus2 --public-ip-addresses k3a-natgw;
az network vnet subnet update -g evanbaker-hyperscale --vnet-name k3a-vnet -n default --nat-gateway k3a-natgw;
az network nic create \
  --resource-group evanbaker-hyperscale \
  --name etcd-nic \
  --location eastus2 \
  --subnet /subscriptions/37deca37-c375-4a14-b90a-043849bd2bf1/resourceGroups/evanbaker-hyperscale/providers/Microsoft.Network/virtualNetworks/k3a-vnet/subnets/default \
  --accelerated-networking true \
  --network-security-group k3a-nsg

az vm create \
  --resource-group evanbaker-hyperscale \
  --name etcd \
  --image AzureLinux \
  --size Standard_E96as_v6 \
  --admin-username azureuser \
  --authentication-type ssh \
  --ssh-key-values ~/.ssh/id_ed25519.pub \
  --nics etcd-nic \
  --public-ip-address "" \
  --os-disk-size-gb 128

