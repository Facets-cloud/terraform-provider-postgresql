package lazy_provider_wrapper

import (
  "fmt"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Wrap(a func() *schema.Provider, inlineProviderConfigure schema.ConfigureFunc) func() *schema.Provider {
	return func() *schema.Provider {
		provider := a()
		originalSchema := provider.Schema
		provider.Schema = map[string]*schema.Schema{}
		provider.ConfigureFunc = nil

		for key, element := range provider.ResourcesMap {
		 element.Schema["inline_provider"] = &schema.Schema{
		   Type:        schema.TypeList,
		   Optional:    true,
		   Description: "Inline provider.",
		   Elem: &schema.Resource{
		     Schema: originalSchema,
		   },
		   MaxItems: 1,
		   ForceNew: true,
		 }

      if (element.Create != nil) {
        element.Create = wrapResourceMethod(inlineProviderConfigure, element.Create)
      }
      if (element.Read != nil) {
        element.Read = wrapResourceMethod(inlineProviderConfigure, element.Read)
      }
      if (element.Delete != nil) {
        element.Delete = wrapResourceMethod(inlineProviderConfigure, element.Delete)
      }
      if (element.Update != nil) {
        element.Update = wrapResourceMethod(inlineProviderConfigure, element.Update)
      }
      if (element.Exists != nil) {
        element.Exists = wrapResourceExistsMethod(inlineProviderConfigure, element.Exists)
      }

      fmt.Println("Key:", key, "=>", "Element:", element.Schema)
    }

    return provider
  }
}

func wrapResourceMethod(inlineProviderConfigure schema.ConfigureFunc, originalFn func(*schema.ResourceData, interface{}) error) func(data *schema.ResourceData, i interface{}) error {
  return func(data *schema.ResourceData, i interface{}) error {
    c, _ := inlineProviderConfigure(data)
    return originalFn(data, c)
  }
}

func wrapResourceExistsMethod(inlineProviderConfigure schema.ConfigureFunc, originalFn func(*schema.ResourceData, interface{}) (bool, error)) func(data *schema.ResourceData, i interface{}) (bool, error) {
  return func(data *schema.ResourceData, i interface{}) (bool, error) {
    c, _ := inlineProviderConfigure(data)
    return originalFn(data, c)
  }
}


