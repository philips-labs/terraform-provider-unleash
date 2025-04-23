data "unleash_api_tokens" "bobjoe_tokens" {
  token_name = "bobjoe"
}

output "tokens" {
  value = data.unleash_api_tokens.bobjoe_tokens.tokens
}
