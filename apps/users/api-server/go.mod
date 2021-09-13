module github.com/SunSince90/ship-krew/users/api-server

go 1.17

require (
	github.com/SunSince90/ship-krew/users/api v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v6 v6.7.1
	github.com/gofiber/fiber/v2 v2.18.0
	github.com/rs/zerolog v1.25.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/pflag v1.0.5
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.29.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/SunSince90/ship-krew/users/api => ../api
