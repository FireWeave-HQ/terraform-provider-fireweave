package provider_test

import (
	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"fireweave": providerserver.NewProtocol6WithError(provider.New("test")()),
}
