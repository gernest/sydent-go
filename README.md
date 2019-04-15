<p align="center">
  <img src="logo.png"/>
</p>

 # sydent-go

**short**
port of [sydent](https://github.com/matrix-org/sydent) to go using postgres as
the data store plus extra stuff for easy deployment and management.

**long**
sysent-go is a port of [sydent]() which is a reference implementation of the identity service specification defined my matrix.org. sydent-go is written in the Go programming language and comes with additional features that makes it ideal for production deployments.

## Features

- Fully ports [sydent]() meaning you can use this as drop in replacement.
- High performance, this is designed to scale well leveraging concurrency primitives provided by Go
- Well tested, making it easy to catch bugs before they roll to production(There is no tests in the reference implementation)
- Use postgresql ( this will make almost 80% of your worries go away)
- Well documented (ahead of you it is the longest well detailed and neat documentation of all what you need to painlessly manage this service, BE WARNED IT IS LONG SO USE TABLE OF CONTENTS TO NAVIGATE)
- heavily instrumented (exports lots of metrics about the running service)
- i18n aware , you can help translate the email messages
- Single binary etc
