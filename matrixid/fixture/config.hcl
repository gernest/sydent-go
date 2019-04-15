mode = "test"

server {
  name             = "matrixid"
  port             = "9891"
  ssl_cert         = "fixture/certs/localhost.pem"
  ssl_key          = "fixture/certs/localhost-key.pem"
  client_http_base = "https://localhost:9891"

  crypto {
    algorithm   = "ed25519"
    version     = "0"
    signing_key = "Je3daK6jsLslp585i/jt/2ioP563Sn3py5VTuyrEYfXjf2GI&#43;oBs0c&#43;rochHDZewd30ws8UfHzzw0DB4trbuYQ"
    verify_key  = "439hiPqAbNHPq6HIRw2XsHd9MLPFHx888NAweLa27mE"
  }
}

db {
  driver = "postgres"
  conn   = "user=postgres dbname=killua sslmode=disable"
}

email {
  provider "smtp" {
    state = "disabled"

    settings {
      host     = "smtp.gmail.com"
      password = ""
      port     = "587"
      username = ""
    }
  }

  provider "sendgrid" {
    state = "disabled"

    settings {
      api_key = ""
    }
  }

  invite {
    from     = ""
    template = "invite"
  }

  verification {
    from                   = ""
    Template               = "verification"
    response_page_template = "verify_response"
  }
}
