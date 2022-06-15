data "unleash_users" "admin" {
  query = "admin"
}

output "admin_users" {
  value = data.unleash_users.admin.users
}