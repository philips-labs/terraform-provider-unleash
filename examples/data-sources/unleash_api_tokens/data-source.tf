data "unleash_api_tokens" "bobjoe_tokens" {
  username = "bobjoe"
}

output "tokens" {
  value = data.unleash_api_tokens.bobjoe_tokens.tokens
}
