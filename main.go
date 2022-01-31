package main

import (
  "context"
  "flag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
  "github.com/terraform-providers/terraform-provider-postgresql/lazy_provider_wrapper"
  "github.com/terraform-providers/terraform-provider-postgresql/postgresql"
  "log"
)

func main() {
  var debugMode bool
  flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
  flag.Parse()

  opts := &plugin.ServeOpts{
    ProviderFunc: lazy_provider_wrapper.Wrap(postgresql.Provider, postgresql.InlineProviderConfigure)}

  if debugMode {
    err := plugin.Debug(context.Background(), "rr0hit/postgresql", opts)
    if err != nil {
      log.Fatal(err.Error())
    }
    return
  }

  plugin.Serve(opts)
}
