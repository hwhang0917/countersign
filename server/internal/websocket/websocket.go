package websocket

import (
	"log"
	"strconv"

	"github.com/bytedance/sonic"
	fiberWs "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/hwhang0917/countersign/internal/hash"
	"github.com/hwhang0917/countersign/pkg/config"
)

type Request struct {
	Command string `json:"command"`
	AskText string `json:"ask_text"`
}
type Response struct {
	Command string `json:"command"`
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func UpgradeWebSocket(c *fiber.Ctx) error {
	if fiberWs.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func GetOTPHandler() fiber.Handler {
	return fiberWs.New(func(c *fiberWs.Conn) {
		ip := c.IP()
		log.Printf("Client connected [%s]", ip)

		askKey := c.Params("askKey")
		var req Request
		var res Response

		for {
			var (
				rawRequest []byte
				err        error
			)
			if _, rawRequest, err = c.ReadMessage(); err != nil {
				log.Printf("Client disconnected [%s]", ip)
				break
			}
			if err := sonic.Unmarshal(rawRequest, &req); err != nil {
				log.Fatalf("Error unmarshaling: %v", err)
				continue
			}

			switch req.Command {
			case "ask":
				otp := hash.GetOTP(askKey)
				res.Success = true
				res.Command = "ask"
				res.Data = otp
				break
			case "time":
				interval := config.GetInterval()
				timeLeft := strconv.FormatInt(hash.GetTimeLeft(int64(interval)), 10)
				res.Success = true
				res.Command = "time"
				res.Data = timeLeft
				break
			default:
				res.Success = false
				res.Command = req.Command
				res.Data = "Invalid command"
				break
			}
			if err := c.WriteJSON(res); err != nil {
				log.Fatalf("Error writing message: %v", err)
			}
		}
	})
}
