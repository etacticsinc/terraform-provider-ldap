# Resource: ldap_user

Creates an LDAP user account. 

## Example Usage

### Microsoft Active Directory
```hcl
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

### OpenLDAP
```hcl
resource "ldap_user" "jsmith" {
  object_class   = ["posixAccount", "shadowAccount", "inetOrgPerson"]
  cn             = "John C Smith"
  path           = "OU=Users,OU=Example,DC=corp,DC=example,DC=com"
  home_directory = "/home/jsmith"
  surname        = "smith"
  uid            = "jsmith"
  uid_number     = 20001
  gid_number     = 100
}
```

## Argument Reference

* `city` - (Optional) Specifies the town or city.

* `cn` - (Required) The name that represents the object. Used to perform searches.

* `country` - (Optional) Specifies the country or region code.

* `display_name` - (Optional) The display name for an object.

* `email_address` - (Optional) Specifies the user's e-mail address.

* `gid_number` - (Optional) Contains an integer value that uniquely identifies a group in an administrative domain.

* `given_name` - (Optional) Contains the given name (first name) of the user.

* `home_directory` - (Optional) The home directory for the account.

* `name` - (Optional) Specifies the name of the object.

* `object_class` - (Optional) The list of classes from which this object is derived. Defaults to 
```hcl
["top", "person","organizationalPerson","user"]
```

* `path` - (Required) Specifies the X.500 path of the OU or container where the new object is created.

* `postal_code` - (Optional) Specifies the postal code or zip code.

* `sam_account_name` - (Optional) Specifies the Security Account Manager (SAM) account name of the user.

* `sam_account_type` - (Optional) Specifies the Security Account Manager (SAM) account type of the user.

* `street_address` - (Optional) Specifies a street address.

* `state` - (Optional) Specifies a state or province.

* `surname` - (Optional) Specifies the user's last name or surname.

* `uid` - (Optional) A user ID.

* `uid_number` - (Optional) Contains a number that uniquely identifies a user in an administrative domain.

* `user_principal_name` - (Optional) Specifies a user principal name (UPN) in the format <USER>@<DNS-domain-name>.

## Attribute Reference

* `id` - The distinguished name of the LDAP user
