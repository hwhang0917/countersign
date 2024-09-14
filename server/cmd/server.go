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

type InputMessage struct {
	Command string `json:"command"`
	AskText string `json:"ask_text"`
}

type ResponseMessage struct {
	Command      string `json:"command"`
	ResponseText string `json:"response_text"`
	SecondsLeft  int    `json:"seconds_left"`
}

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

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:X-API-Key",
		Validator: validateAPIKey,
	}))
	app.Use("/ws", websocket.UpgradeWebSocket)
	app.Get("/ws/otp", websocket.GetOTPHandler())

	log.Fatal(app.Listen(":" + config.GetPort()))

	defer database.CloseDB()
	defer app.Shutdown()
}
