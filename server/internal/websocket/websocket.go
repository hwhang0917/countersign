package websocket

import (
	"log"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	fiberWs "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/hwhang0917/countersign/internal/hash"
	"github.com/hwhang0917/countersign/pkg/config"
)

type Request struct {
	AskText string `json:"ask_text"`
}
type Response struct {
	Success  bool   `json:"success"`
	AskText  string `json:"ask_text,omitempty"`
	OTP      string `json:"otp,omitempty"`
	TimeLeft string `json:"time_left,omitempty"`
}

func getRealIP(c *fiberWs.Conn) string {
	if ip := c.Headers("X-Forwarded-For"); ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	if ip := c.Headers("X-Real-IP"); ip != "" {
		return ip
	}
	return c.IP()
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
		ip := getRealIP(c)
		log.Printf("Client connected [%s]", ip)

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

			if req.AskText == "" {
				res.Success = false
				res.AskText = ""
				res.OTP = ""
				res.TimeLeft = ""
			} else {
				res.Success = true
				res.AskText = req.AskText
				res.OTP = hash.GetOTP(req.AskText)
				res.TimeLeft = strconv.FormatInt(hash.GetTimeLeft(int64(config.GetInterval())), 10)
			}

			if err := c.WriteJSON(res); err != nil {
				log.Fatalf("Error writing message: %v", err)
			}
		}
	})
}
