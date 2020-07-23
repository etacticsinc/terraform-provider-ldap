# Resource: ldap_organizational_unit

Creates an LDAP organizational unit. 

## Example Usage

```hcl
resource "ldap_organizational_unit" "servers" {
  ou          = "Servers"
  path        = "OU=Example,DC=corp,DC=example,DC=com"
}
```

## Argument Reference

The following arguments are supported:



## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The distinguished name of the LDAP organizational unit
