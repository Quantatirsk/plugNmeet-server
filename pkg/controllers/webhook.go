package controllers

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/livekit/protocol/livekit"
	"github.com/mynaparrot/plugnmeet-server/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

func HandleWebhook(c *fiber.Ctx) error {
	data := c.Body()
	body := make([]byte, len(data))
	copy(body, data)

	token := c.Request().Header.Peek("Authorization")
	authToken := make([]byte, len(token))
	copy(authToken, token)

	if len(authToken) == 0 {
		return c.SendStatus(fiber.StatusForbidden)
	}

	m := models.NewAuthTokenModel()
	// here request is coming from livekit
	// so, we'll use livekit secret to validate
	ourHash, err := m.ValidateLivekitWebhookToken(string(authToken))
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	sha := sha256.Sum256(body)
	hash := base64.StdEncoding.EncodeToString(sha[:])

	if subtle.ConstantTimeCompare([]byte(ourHash), []byte(hash)) != 1 {
		return c.SendStatus(fiber.StatusForbidden)
	}

	op := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	event := new(livekit.WebhookEvent)
	if err = op.Unmarshal(body, event); err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	models.NewWebhookModel(event)

	return c.SendStatus(fiber.StatusOK)
}
