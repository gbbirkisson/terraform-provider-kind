provider "kind" {

}

resource "kind" "cluster" {
  name = "kind"
  config = abspath("${path.module}/config.yaml")
}
