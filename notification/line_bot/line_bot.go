package LineBot

import (
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	// Utils "notification-server/utils"
	Utils "notification-server/utils"
)

func SendLineMessage(newMessage string) {

	// userLineID is Creator of line bot will receive text message
	userLineID := Utils.GetEnv("USER_LINE_ID")

	channelSecret := Utils.GetEnv("LINE_BOT_CHANNEL_SECRET")
	channelAccessToken := Utils.GetEnv("LINE_BOT_CHANNEL_TOKEN")

	message := linebot.NewTextMessage(newMessage)

	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	var messages []linebot.SendingMessage

	// append some message to messages
	messages = append(messages, message)

	_, err = bot.PushMessage(userLineID, messages...).Do()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Line Bot sent successfully")
	}

}
