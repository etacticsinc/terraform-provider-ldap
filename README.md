# Terraform Provider for LDAP

* LDAP Reference: [ldap.com](https://ldap.com/)
* Terraform Tutorials: [learn.hashicorp.com](https://learn.hashicorp.com/terraform?track=getting-started#getting-started)

The Terraform LDAP provider is a plugin that allows for automated management of an existing directory server and its resources. 

This provider is maintained and used internally by the Etactics, Inc. development team.

## Getting Started

- [Using the Provider](docs/index.md)
- [Creating a User Account](docs/resources/user.md)
- [Creating a Group](docs/resources/group.md)
- [Creating an Organizational Unit](docs/resources/organizational_unit.md)

## Installation

See [Third-party Plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) in the Terraform documentation for installation instructions.

Pre-compiled binaries are available from the [Releases](https://github.com/etacticsinc/terraform-provider-ldap/releases) page.

## Building from Source

System Requirements:
- [Terraform](https://www.terraform.io/downloads.html) 0.12
- [Go](https://golang.org/doc/install) 1.13+

Clone the repository and build the provider:

```sh
git clone https://github.com/etacticsinc/terraform-provider-ldap.git
cd terraform-provider-ldap
go build -o terraform-provider-ldap
```

After building, follow the [plugin installation instructions](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) and initialize the provider to begin use.

