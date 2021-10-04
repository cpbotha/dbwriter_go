Minimal database writer API with gin, gorm and swagger put together by a go
newbie in 2021 to test the developer experience.

See the blog post [Developer experience setting up a minimal API in Go, C# and
Python at
vxlabs.com](https://vxlabs.com/2021/10/03/dx-minimal-api-go-csharp-python/) for
more information.

See you,
https://charlbotha.com/

## quickstart

```shell
go run
```

... and then browse to `localhost:8080/swagger/index.html`

## Update swagger docs

```shell
# install the swag command which you'll need to update swagger docs
go install github.com/swaggo/swag/cmd/swag@latest
# whenever you update code / annotated comments
~/go/bin/swag init
# in theory this should install deps listed in go.sum if necessary
go run
```

## References

Examples of adding swagger docs to simple gin apps with bare controller functions:

- https://levelup.gitconnected.com/tutorial-generate-swagger-specification-and-swaggerui-for-gin-go-web-framework-9f0c038483b5
- https://medium.com/@mayur.das4/generate-restful-api-documentation-by-integration-swagger-baa52aefd2dd

## How I got started

Reminders because I'm a newbie:

```shell
go init
go get -u github.com/jinzhu/gorm
go get github.com/jinzhu/gorm/dialects/sqlite
go get github.com/gin-gonic/gin
go install github.com/swaggo/swag/cmd/swag@latest
~/go/bin/swag init # setup docs/ dir; you have to run this every time you update annotations
go get -v github.com/swaggo/gin-swagger
```
