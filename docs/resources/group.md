# Resource: ldap_group

Creates an LDAP group. 

## Example Usage

### Microsoft Active Directory
```hcl
resource "ldap_group" "sales_managers" {
  cn          = "Sales Managers"
  path        = "OU=Groups,OU=Example,DC=corp,DC=example,DC=com"
  description = "Example Corp Sales Managers"
  members     = [ldap_user.jsmith.id]
}
```

### OpenLDAP
```hcl
resource "ldap_group" "sales_managers" {
  object_class = ["top", "posixGroup"]
  cn           = "Sales Managers"
  path         = "OU=Groups,OU=Example,DC=corp,DC=example,DC=com"
  description  = "Example Corp Sales Managers"
  gid_number   = 10001
  member_uids  = [ldap_user.jsmith.uid]
}
```

## Argument Reference

The following arguments are supported:

* `cn` - (Required) The name that represents the object. Used to perform searches.

* `description` - (Optional) Specifies a description of the object.

* `display_name` - (Optional) The display name for an object.

* `gid_number` - (Optional) Contains an integer value that uniquely identifies a group in an administrative domain.

* `group_category` - (Optional) Specifies the category of the group. The acceptable values for this parameter are ``"Distribution"`` and ``"Security."``

* `group_scope` - (Optional)Specifies the scope of the group. The acceptable values for this parameter are "Global," "DomainLocal" and "Universal."

* `homepage` - (Optional) Specifies the URL of the home page of the object.

* `members` - (Optional) Specifies an array of user, group, and computer objects to add to the group.

* `member_uids` - (Optional) Contains the login names of the members of a group.

* `name` - (Optional) Specifies the name of the object.

* `object_class` - (Optional) The list of classes from which this object is derived. Defaults to ``["top", "group"]``

* `path` - (Required) Specifies the X.500 path of the OU or container where the new object is created.

* `sam_account_name` - (Optional) Specifies the Security Account Manager (SAM) account name of the group.

* `sam_account_type` - (Optional) Specifies the Security Account Manager (SAM) account type of the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The distinguished name of the LDAP group
