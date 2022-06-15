resource "unleash_user" "my_user" {
  name       = "Bob Joe"
  email      = "bob.joe@gmail.com"
  username   = "bobjoe"
  root_role  = "Editor"
  send_email = false
}