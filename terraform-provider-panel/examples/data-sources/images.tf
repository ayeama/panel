terraform {
  required_providers {
    panel = {
      source = "github.com/ayeama/panel"
    }
  }
}

provider "panel" {
  alias = "local"

  endpoint = "http://localhost:8000"
}

provider "panel" {
  alias = "remote"

  endpoint = "" // TODO
  verify   = false
  username = "" // TODO
  password = "" // TODO
}

data "panel_images" "local_images" {
    provider = panel.local
}

data "panel_images" "remote_images" {
    provider = panel.remote
}
