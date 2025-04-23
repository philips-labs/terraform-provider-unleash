resource "unleash_api_token" "my_token" {
  username    = "bobjoe"
  type        = "client"
  expires_at  = "2024-10-19"
  environment = "development"
  projects    = ["*"]
}