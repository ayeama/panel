terraform {
  required_providers {
    panel = {
      source = "github.com/ayeama/panel"
    }
  }
}

provider "panel" {
  endpoint = "http://localhost:8000"
}

resource "panel_server" "minecraft" {
  count = 10

  image = "localhost/ayeama/panel/server/minecraft:0.0.1-jre21"
}
