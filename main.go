package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/rh01/terraform-provider-kubevirt-yd/kubevirt"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubevirt.Provider})
}
