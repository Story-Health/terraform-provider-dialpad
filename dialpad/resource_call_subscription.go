package dialpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type callSubscriptionRequest struct {
	CallStates []string `json:"call_states"`
	EndpointId string   `json:"endpoint_id"`
}

type callSubscriptionResponse struct {
	Id         string          `json:"id"`
	CallStates []string        `json:"call_states"`
	Webhook    WebhookResponse `json:"webhook"`
}

func resourceCallSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCallSubscriptionCreate,
		ReadContext:   resourceCallSubscriptionRead,
		UpdateContext: resourceCallSubscriptionUpdate,
		DeleteContext: resourceCallSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"call_states": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCallSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subId := d.Id()
	client := m.(*client)

	req, err := client.NewRequest("GET", fmt.Sprintf("/subscriptions/call/%s", subId), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	sub := callSubscriptionResponse{}
	err = json.Unmarshal(res, &sub)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("call_states", sub.CallStates)
	d.Set("endpoint_id", sub.Webhook.Id)

	return diags
}

func resourceCallSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client)

	statesSet := d.Get("call_states").(*schema.Set).List()
	states := make([]string, len(statesSet))
	for i, raw := range statesSet {
		states[i] = raw.(string)
	}

	subRequest := callSubscriptionRequest{
		CallStates: states,
		EndpointId: d.Get("endpoint_id").(string),
	}

	body, err := json.Marshal(subRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	req, err := client.NewRequest("POST", "/subscriptions/call", bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	sub := callSubscriptionResponse{}
	err = json.Unmarshal(res, &sub)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sub.Id)

	return resourceCallSubscriptionRead(ctx, d, m)
}

func resourceCallSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client)
	id := d.Id()

	statesSet := d.Get("call_states").(*schema.Set).List()
	states := make([]string, len(statesSet))
	for i, raw := range statesSet {
		states[i] = raw.(string)
	}

	subRequest := callSubscriptionRequest{
		CallStates: states,
		EndpointId: d.Get("endpoint_id").(string),
	}

	body, err := json.Marshal(subRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := client.NewRequest("PATCH", fmt.Sprintf("/subscriptions/call/%s", id), bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCallSubscriptionRead(ctx, d, m)
}

func resourceCallSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client)

	id := d.Id()

	req, err := client.NewRequest("DELETE", fmt.Sprintf("/subscriptions/call/%s", id), nil)
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
