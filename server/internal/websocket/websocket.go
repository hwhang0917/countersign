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
		askKey := c.Params("askKey")
		var req Request
		var res Response
		for {
			_, rawRequest, err := c.ReadMessage()
			if err != nil {
				log.Fatalf("Error reading message: %v", err)
				break
			}
			if err := sonic.Unmarshal(rawRequest, &req); err != nil {
				log.Fatalf("Error unmarshaling: %v", err)
				break
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
			case "disconnect":
				res.Success = true
				res.Command = "disconnect"
				res.Data = "Disconnecting..."
				c.Close()
				break
			default:
				res.Success = false
				res.Command = req.Command
				res.Data = "Invalid command"
				break
			}
			c.WriteJSON(res)
		}
	})
}
