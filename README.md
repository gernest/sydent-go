<p align="center">
  <img src="logo.png"/>
</p>

 # sydent-go

sysent-go is a port of [sydent]() which is a reference implementation of the identity service specification defined my matrix.org. sydent-go is written in the Go programming language and comes with additional features that makes it ideal for production deployments.

## Features

- Fully ports [sydent]() meaning you can use this as drop in replacement.
- High performance, this is designed to scale well leveraging concurrency primitives provided by Go
- Well tested, making it easy to catch bugs before they roll to production(There is no tests in the reference implementation)
- Use postgresql ( this will make almost 80% of your worries go away)
- heavily instrumented (exports lots of metrics about the running service)
- i18n aware , you can help translate the email messages
- Single binary etc


## Installation

### install from source

You need `go1.12+`

```
go install github.com/gernest/sydent-go
```

### running

```
sydent-go server /path/to/config/file
```

### Configuration

There are two wat to configure the service, you can use environment variables
or you can use a configuration file.

The configuration file uses hcl(hashicorp configuration language). The following is a sample configuration

```hcl
mode = "prod"

server {
  name             = "sydent-go"
  port             = "9891"
  client_http_base = "https://localhost:9891"

  crypto {
    algorithm   = "ed25519"
    version     = "0"
    signing_key = "${SYDENT_PRIVATE_KEY}"
    verify_key  = "${SYDENT_PUBLIC_KEY}"
  }
}

db {
  driver = "postgres"
  conn   = "${SYDENT_DB_CONN}"
}

email {
  provider "smtp" {
    state = "enabled"

    settings {
      host     = "$SYDENT_SMTP_HOST"
      password = "$SYDENT_SMTP_PASSWORD"
      port     = "$SYDENT_SMTP_PORT"
      username = "$SYDENT_SMTP_USERNAME"
    }
  }

  provider "sendgrid" {
    state = "disabled"

    settings {
      api_key = "$SYDENT_SENDGRID_APIKEY"
    }
  }

  invite {
    from     = "$SYDENT_SMTP_USERNAME"
    template = "invite"
  }

  verification {
    from                   = "$SYDENT_SMTP_USERNAME"
    Template               = "verification"
    response_page_template = "verify_response"
  }
}
```

#### `mode`

This defines how the service is run. The values are 

- `prod` if the service is running in production 
- `dev` if the service is running in development.

The value is used to determine the amount of logging and certain features, for
instance it doesnt make sense to instrument a development service. So instrumentation will be off when in `dev` mode.

#### `db`

Setup database connection. We are using postgres by default because that is the
only database we support for now.

- `driver` - the database driver to use , defaults to `postgres`
- `conn` - connection string to use to connect to the database


#### `server`

Configures the webserver.


##### `name`

This is the name that identifies this service. Note that this is very important because it is used to sign the mappings. 

##### `name`

The port number to bind the service to


##### `client_http_base`

This is the base url `http[s]://host:port` that can be used by clients to reach this service. It is used mainle in the emails sent with links for verification/validation that needs to point back to this service. It should be resolvable to the host running this service.