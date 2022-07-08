package dialpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type signature struct {
	Algorithm string `json:"algo"`
	Secret    string `json:"secret"`
	Type      string `json:"type"`
}

type webhookRequest struct {
	HookUrl string `json:"hook_url"`
	Secret  string `json:"secret"`
}

type WebhookResponse struct {
	Id        string    `json:"id"`
	HookUrl   string    `json:"hook_url"`
	Signature signature `json:"signature"`
}

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"hook_url": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  false,
			},
		},
	}
}

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hookId := d.Id()
	client := m.(*client)

	req, err := client.NewRequest("GET", fmt.Sprintf("/webhooks/%s", hookId), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	hook := WebhookResponse{}
	err = json.Unmarshal(res, &hook)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("hook_url", hook.HookUrl)
	d.Set("secret", hook.Signature.Secret)

	return diags
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client)

	hookRequest := webhookRequest{
		HookUrl: d.Get("hook_url").(string),
	}

	secret := d.Get("secret")
	if secret != "" {
		hookRequest.Secret = secret.(string)
	}

	body, err := json.Marshal(hookRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	req, err := client.NewRequest("POST", "/webhooks", bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	hook := WebhookResponse{}
	err = json.Unmarshal(res, &hook)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hook.Id)

	return resourceWebhookRead(ctx, d, m)
}

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client)
	id := d.Id()

	hookRequest := webhookRequest{
		HookUrl: d.Get("hook_url").(string),
	}

	secret := d.Get("secret")
	if secret != "" {
		hookRequest.Secret = secret.(string)
	}

	body, err := json.Marshal(hookRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := client.NewRequest("PATCH", fmt.Sprintf("/webhooks/%s", id), bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWebhookRead(ctx, d, m)
}

func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client)

	id := d.Id()

	req, err := client.NewRequest("DELETE", fmt.Sprintf("/webhooks/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
