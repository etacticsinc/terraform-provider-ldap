# LDAP Provider

The Lightweight Directory Access Protocol (LDAP) provider is used to configure and interact with remote directory services such as Microsoft Active Directory or OpenLDAP. This provider must be configured with the appropriate host address and administrative bind credentials before use.

## Example Usage

```hcl
variable ldap_password {}

# Configure the LDAP Provider
provider "ldap" {
  server        = "ldap://corp.example.com"
  bind_dn       = "CN=Admin,OU=Users,OU=Example,DC=corp,DC=example,DC=com"
  bind_password = var.ldap_password
}

# Create a user account
resource "ldap_user" "jsmith" {
  cn                  = "jsmith"
  path                = "OU=Users,OU=Example,DC=corp,DC=example,DC=com"
  given_name          = "John"
  surname             = "Smith"
  display_name        = "John C Smith"
  email_address       = "jsmith@example.com"
  user_principal_name = "jsmith@corp.example.com"
}
```

## Argument Reference

The following arguments are supported in the LDAP ``provider`` block:

* ``server`` - The LDAP server managed by this provider.
* ``bind_dn`` - The distinguished name of the administrative user account used to access the directory.
* ``bind_password`` - The password used for authentication.
