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

* `city` - (Optional) Specifies the town or city.

* `country` - (Optional) Specifies the country or region code.

* `description` - (Optional) Specifies a description of the object.

* `name` - (Optional) Specifies the name of the object.

* `object_class` - (Optional) The list of classes from which this object is derived. Defaults to ``["top","organizationalUnit"]``

* `ou` - (Required) The organizational unit name.

* `path` - (Required) Specifies the X.500 path of the OU or container where the new object is created.

* `postal_code` - (Optional) Specifies the postal code or zip code.

* `street_address` - (Optional) Specifies a street address.

* `state` - (Optional) Specifies a state or province.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The distinguished name of the LDAP organizational unit (e.g. ``OU=Servers,OU=Example,DC=corp,DC=example,DC=com``)
