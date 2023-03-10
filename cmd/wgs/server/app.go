package server

import (
	"net/http"
	"time"

	"wgs"
	"wgs/internal/config"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/neo-f/soda"
)

var ROUTES = []func(app *soda.Soda){
	RegisterBase,
}

// startHttpServer starts configures and starts an HTTP server on the given URL.
// It shuts down the server if any error is received in the error channel.
func InitApp() *soda.Soda {
	app := soda.New("Wireguard Server", wgs.Version,
		soda.EnableValidateRequest(),
		soda.WithOpenAPISpec(config.URL_OPENAPI),
		soda.WithStoplightElements(config.URL_API_DOC),
		soda.WithFiberConfig(
			fiber.Config{
				ReadTimeout:       time.Second * 10,
				EnablePrintRoutes: true,
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					if err == nil {
						return c.Next()
					}
					status := http.StatusInternalServerError //default error status
					if e, ok := err.(*fiber.Error); ok {     // it's a custom error, so use the status in the error
						status = e.Code
					}
					msg := map[string]interface{}{"code": status, "message": err.Error()}
					return c.Status(status).JSON(msg)
				},
			},
		),
	)
	app.OpenAPI().Info.Contact = &openapi3.Contact{Name: "NEO", Email: "tmpgfw@gmail.com"}
	app.OpenAPI().Info.Description = `wireguard server`
	for _, route := range ROUTES {
		route(app)
	}
	return app
}
