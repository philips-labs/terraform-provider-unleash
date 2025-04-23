resource "unleash_api_token" "my_token" {
  token_name  = "bobjoe"
  type        = "client"
  expires_at  = "2050-04-15T14:30:45Z"
  environment = "development"
  projects    = ["*"]
}