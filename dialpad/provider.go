package dialpad

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("DIALPAD_API_KEY", nil),
				Required:    true,
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dialpad_webhook":           resourceWebhook(),
			"dialpad_call_subscription": resourceCallSubscription(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiKey := d.Get("api_key").(string)
	client := NewClient(apiKey)

	return client, diags
}
