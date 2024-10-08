package cmd

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/hwhang0917/countersign/internal/websocket"
	"github.com/hwhang0917/countersign/pkg/config"
	"github.com/hwhang0917/countersign/pkg/database"
)

func validateAPIKey(_ *fiber.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(config.GetAPIKey()))
	hashedKey := sha256.Sum256([]byte(key))
	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func StartServer() {
	config.LoadConfig()
	database.InitDB()
	database.MigrateModels(database.GetDB())
	database.PopulateWords(database.GetDB())

	app := fiber.New(fiber.Config{
		JSONEncoder:    sonic.Marshal,
		JSONDecoder:    sonic.Unmarshal,
		ProxyHeader:    "X-Forwarded-For",
		TrustedProxies: []string{"127.0.0.1"},
	})
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(healthcheck.New())
	// app.Use(keyauth.New(keyauth.Config{
	// 	KeyLookup: "header:X-API-Key",
	// 	Validator: validateAPIKey,
	// }))
	app.Use("/ws", websocket.UpgradeWebSocket)
	app.Get("/ws/otp", websocket.GetOTPHandler())
	app.Get("/api/interval", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"interval": config.GetInterval(),
		})
	})
	app.Static("/", "./public")

	log.Fatal(app.Listen(":" + config.GetPort()))

	defer database.CloseDB()
	defer app.Shutdown()
}
