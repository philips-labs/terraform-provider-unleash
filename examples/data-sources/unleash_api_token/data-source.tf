data "unleash_api_token" "filter" {
  token_name = "bobjoe"
  projects   = ["*"]
}

output "token" {
  value = data.unleash_api_token.filter.token
}
