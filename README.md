<a href="https://etactics.com">
    <img src="https://github.com/etacticsinc/terraform-provider-ldap/blob/master/etactics-logo.png" alt="Etactics logo" title="Etactics" align="right" height="50" />
</a>

<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

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

## Build from Source

If you're building the provider, follow the [plugin installation instructions](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

Requirements:
- [Terraform](https://www.terraform.io/downloads.html) 0.12
- [Go](https://golang.org/doc/install) 1.13+

Clone repository to: `$GOPATH/src/github.com/etacticsinc/terraform-provider-ldap`

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/etacticsinc/terraform-provider-ldap
go build
```

Install the provider:
```sh
cd $GOPATH/src/github.com/etacticsinc/terraform-provider-ldap
go install
```
