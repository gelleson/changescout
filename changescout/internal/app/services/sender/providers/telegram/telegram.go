package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"io"
	"net/http"
)

//go:generate mockery --name Doer
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Telegram struct {
	doer Doer
}

func New() *Telegram {
	return &Telegram{
		doer: http.DefaultClient,
	}
}

type message struct {
	ChatID    *string `json:"chat_id"`
	Text      string  `json:"text"`
	ParseMode string  `json:"parse_mode"`
}

func encode[T any](v T) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func (t *Telegram) Send(notification string, conf domain.Notification) error {
	tUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *conf.Token)

	req, _ := http.NewRequest("POST", tUrl, encode(message{
		ChatID:    conf.Destination,
		Text:      notification,
		ParseMode: "Markdown",
	}))

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.doer.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}

	return nil
}
