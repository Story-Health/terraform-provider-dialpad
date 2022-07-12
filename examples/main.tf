terraform {
  required_providers {
    dialpad = {
      version = "0.1.0"
      source  = "Story-Health/dialpad"
    }
  }
}

provider "dialpad" {
  api_key = "My-Api-Key"
}

resource "dialpad_webhook" "example" {
  hook_url = "https://example.com/hook"
  secret   = "my very secret secret"
}

resource "dialpad_call_subscription" "all_calls" {
  endpoint_id = dialpad_webhook.example.id
  call_states = ["all"]
}
