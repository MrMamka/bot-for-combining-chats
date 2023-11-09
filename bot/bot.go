package bot

import (
	"errors"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const welcomeText string = "Здесь должно быть приветственное сообщение"

var NotFoundEnvError = errors.New("env not found")

type TelegramBot struct {
	token string
	bot   *tgbotapi.BotAPI
}

func NewTelegramBot() *TelegramBot {
	return new(TelegramBot)
}

func (tb *TelegramBot) SetTokenFromEnv(env string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	tb.token = os.Getenv(env)
	if tb.token == "" {
		return NotFoundEnvError
	}
	return nil
}

func (tb *TelegramBot) Start(isDebug bool) error {
	var err error
	tb.bot, err = tgbotapi.NewBotAPI(tb.token)
	if err != nil {
		return err
	}

	tb.bot.Debug = isDebug

	log.Printf("Authorized on account %s", tb.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	updates := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msgReq := update.Message
			log.Printf("[%s] %s", msgReq.From.UserName, msgReq.Text)

			switch msgReq.Text { // заменить на мапу
			case "\\start":
				err = tb.startMsgHandler(msgReq)
			default:
				err = tb.defaultMsgHandler(msgReq)
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tb *TelegramBot) startMsgHandler(msgReq *tgbotapi.Message) error {
	msgResp := tgbotapi.NewMessage(msgReq.Chat.ID, welcomeText)

	_, err := tb.bot.Send(msgResp)
	if err != nil {
		return err
	}
	return nil
}

func (tb *TelegramBot) defaultMsgHandler(msgReq *tgbotapi.Message) error {
	msgResp := tgbotapi.NewMessage(msgReq.Chat.ID, msgReq.Text)
	msgResp.ReplyToMessageID = msgReq.MessageID

	_, err := tb.bot.Send(msgResp)
	if err != nil {
		return err
	}
	return nil
}
