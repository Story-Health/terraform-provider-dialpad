package dialpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CallSubscriptionRequest struct {
	CallStates []string `json:"call_states"`
	EndpointId string `json:"endpoint_id"`
}

type CallSubscriptionResponse struct {
	Id        string    `json:"id"`
	CallStates []string `json:"call_states"`
	Webhook WebhookResponse `json:"webhook"`
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

func resourceCallSubscriptionRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subId := data.Id()
	client := meta.(*Client)

	req, err := client.NewRequest("GET", fmt.Sprintf("/subscriptions/call/%s", subId), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	sub := CallSubscriptionResponse{}
	err = json.Unmarshal(res, &sub)

	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("call_states", sub.CallStates)
	data.Set("endpoint_id", sub.Webhook.Id)

	return diags
}

func resourceCallSubscriptionCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	statesSet := data.Get("call_states").(*schema.Set).List()
	states := make([]string, len(statesSet))
	for i, raw := range statesSet {
		states[i] = raw.(string)
	}

	subRequest := CallSubscriptionRequest{
		CallStates: states,
		EndpointId: data.Get("endpoint_id").(string),
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

	sub := CallSubscriptionResponse{}
	err = json.Unmarshal(res, &sub)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(sub.Id)

	return resourceCallSubscriptionRead(ctx, data, meta)
}

func resourceCallSubscriptionUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	id := data.Id()

	statesSet := data.Get("call_states").(*schema.Set).List()
	states := make([]string, len(statesSet))
	for i, raw := range statesSet {
		states[i] = raw.(string)
	}

	subRequest := CallSubscriptionRequest{
		CallStates: states,
		EndpointId: data.Get("endpoint_id").(string),
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

	return resourceCallSubscriptionRead(ctx, data, meta)
}

func resourceCallSubscriptionDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Client)

	id := data.Id()

	req, err := client.NewRequest("DELETE", fmt.Sprintf("/subscriptions/call/%s", id), nil)
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
