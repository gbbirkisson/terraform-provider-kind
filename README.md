Terraform Provider Kind
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

Guðmundur Björn Birkisson

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Kind](https://github.com/kubernetes-sigs/kind) v0.5.1
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Usage
---------------------

```
provider "kind" {

}

# For example, use custom config
resource "kind" "cluster" {
  name = "kind"
  config = abspath("${path.module}/config.yaml") // Optional
}
```

Building and testing the Provider
---------------------

You will need go 1.11+ installed to use go modules.


```bash
# Build
make install
# Test
cd example
terraform plan
```