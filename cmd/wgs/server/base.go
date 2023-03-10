package server

import (
	"runtime"
	"time"
	"wgs"
	"wgs/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/neo-f/soda"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

func RegisterBase(app *soda.Soda) {
	app.Use(
		cors.New(),
		compress.New(),
		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Ctx(c.UserContext()).Error().Interface("err", e).Msg(string(buf))
			},
		}),
		requestid.New(
			requestid.Config{
				Generator:  func() string { return xid.New().String() },
				ContextKey: "trace_id",
			},
		),
		logger.New(logger.Config{
			TimeFormat: time.RFC3339,
			Format:     "${time} [${locals:trace_id}] ${status} - ${latency} ${method} ${path}\n",
		}),
	)

	app.Get(config.URL_VERSION, version).
		SetSummary("System Version").
		AddJSONResponse(200, SystemVersion{}).
		AddTags("System").OK()
}

type SystemVersion struct {
	Version   string `json:"version"    oai:"description=git version"`
	BuildTime string `json:"build_time" oai:"description=build time"`
	GoVersion string `json:"go_version" oai:"description=go version"`
}

func version(c *fiber.Ctx) error {
	return c.JSON(SystemVersion{
		Version:   wgs.Version,
		BuildTime: wgs.BuildTime,
		GoVersion: wgs.GoVersion,
	})
}
