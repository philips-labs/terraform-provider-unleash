data "unleash_api_token" "filter" {
  username = "bobjoe"
  projects = ["*"]
}

output "token" {
  value = data.unleash_api_token.filter.token
}
