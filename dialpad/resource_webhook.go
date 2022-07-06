package dialpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Signature struct {
	Algorithm string `json:"algo"`
	Secret    string `json:"secret"`
	Type      string `json:"type"`
}

type WebhookRequest struct {
	HookUrl string `json:"hook_url"`
	Secret  string `json:"secret"`
}

type WebhookResponse struct {
	Id        string    `json:"id"`
	HookUrl   string    `json:"hook_url"`
	Signature Signature `json:"signature"`
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

func resourceWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hookId := data.Id()
	client := meta.(*Client)

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

	data.Set("hook_url", hook.HookUrl)
	data.Set("secret", hook.Signature.Secret)

	return diags
}

func resourceWebhookCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	hookRequest := WebhookRequest{
		HookUrl: data.Get("hook_url").(string),
	}

	secret := data.Get("secret")
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

	data.SetId(hook.Id)

	return resourceWebhookRead(ctx, data, meta)
}

func resourceWebhookUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	id := data.Id()

	hookRequest := WebhookRequest{
		HookUrl: data.Get("hook_url").(string),
	}

	secret := data.Get("secret")
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

	return resourceWebhookRead(ctx, data, meta)
}

func resourceWebhookDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Client)

	id := data.Id()

	req, err := client.NewRequest("DELETE", fmt.Sprintf("/webhooks/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
