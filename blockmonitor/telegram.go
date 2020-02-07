package blockmonitor

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
)

var (
	botAPITokenFlag = cli.StringFlag{
		Name:  "apiToken",
		Usage: "The API token",
		Value: "935265043:AAG02pbFdpOhT31ALs1tzlLTl_NVnxYuTF4",
	}
	chatIDFlag = cli.Int64Flag{
		Name:  "chatId",
		Usage: "The ID of group/chanel",
		Value: -375817595,
	}
)

// Telegram struct
type Telegram struct {
	ChatId  int64
	Bot     *tgbotapi.BotAPI
	IsDebug bool
}

// NewTeleClientFlag returns flags for telegram
func NewTeleClientFlag() []cli.Flag {
	return []cli.Flag{botAPITokenFlag, chatIDFlag}
}

// NewTeleClientFromFlag returns new telegram client
func NewTeleClientFromFlag(ctx *cli.Context) (*Telegram, error) {
	var (
		botAPIToken = ctx.String(botAPITokenFlag.Name)
		chatID      = ctx.Int64(chatIDFlag.Name)
	)

	telegram := &Telegram{
		ChatId:  chatID,
		IsDebug: false,
	}
	bot, err := tgbotapi.NewBotAPI(botAPIToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = telegram.IsDebug

	telegram.Bot = bot
	return telegram, nil
}

func (t *Telegram) SendMessage(content string, caption string) error {
	text := fmt.Sprintf("<b>%s</b>: %s", caption, content)
	msg := tgbotapi.NewMessage(t.ChatId, text)
	msg.ParseMode = "html"
	_, err := t.Bot.Send(msg)

	if err != nil {
		return err
	}
	return nil
}
