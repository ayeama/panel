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

  endpoint = "https://panel.ayeama.com/api"
  verify   = false
  username = "" // TODO
  password = "" // TODO
}
