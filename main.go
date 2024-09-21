package main

import (
	"github.com/SamuelDBines/terraform-provider-tekton/tekton"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tekton.Provider,
	})
}
