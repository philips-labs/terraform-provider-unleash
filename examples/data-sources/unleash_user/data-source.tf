data "unleash_user" "user" {
  id = 1
}

output "user_details" {
  value = data.unleash_user.user
}