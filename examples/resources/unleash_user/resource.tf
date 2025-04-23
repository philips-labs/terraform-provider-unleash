resource "unleash_user" "my_user" {
  name       = "Bob Joe"
  email      = "bob.joe@gmail.com"
  root_role  = "Editor"
  send_email = false
}