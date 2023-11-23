package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	// LineBot "notification-server/notification/line_bot"
	// Email "notification-server/notification/email"
	NotificationDTO "notification-server/dto"
	LineNotify "notification-server/notification/line-notify"
	// /Users/tharintantayothin/Desktop/arv/deepzoom/last/deepzoom/v2/notification-server/notification/line_notify/line_notify.go
	// Notification "notification-server/dto"
)

// Message represent the structure of a message sent by a "sender" client
type MessageDetection struct {
	JobId     string      `json:"job_id"`
	ConfigId  string      `json:"object_id"`
	DeviceId  string      `json:"device_id"`
	TimeStamp string      `json:"timestamp"`
	ImageUrl  []byte      `json:"presigned_url"`
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Response  interface{} `json:"responseBody"`
	// Custom payload
	Description string `json:"description"` // "Found people > 100 in ConfigName,DeviceId"
	Line        string `json:"line"`

	Email string `json:"email"`
	// SendType    []string `json:"sendType"` // ["line","email","sms","etc"]
	SendType string `json:"sendType"` // "line","email"

	ConfigName string `json:"configName"`
	RoiName    string `json:"roiName"`
}

// WS
type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

var senderUpgrader = websocket.Upgrader{
	ReadBufferSize:  102400,
	WriteBufferSize: 102400, CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var x Notification.NotificationReq

// Concurrency handling - sending live stream messages
func (c *Connection) Send(message NotificationDTO.MessageDetection) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Socket.WriteJSON(message)
}

func BroadcastToReceiver(msg NotificationDTO.MessageDetection) {

	// title := "Title"
	// message := "ConfigId : " + msg.ConfigId
	// // send to LineBot
	// LineBot.SendLineMessage(title + " : " + message)

	// // send to LineNotify
	// LineNotify.SendLineNotify(title + " : " + message)

	// // send to Email
	// sender := "sender@mail.com"
	// sendTo := []string{
	// 	"to@mail.com",
	// }
	// Email.SendEmail(title, message, sender, sendTo)
	log.Println("[NOTI] MessageDetection")
	if strings.ToLower(msg.SendType) == "line" {
		go LineNotify.SendLineNotifyMsgAndImg(msg)
	} else if strings.ToLower(msg.SendType) == "email" {
		log.Println("email")
		//go SendEmail("Notification", "ConfigId : "+msg.ConfigId, "from@example.com", sendTo)
	} else if strings.ToLower(msg.SendType) == "sms" {
		log.Println("sms")
	}

}

func senderHandler(conn *websocket.Conn, r *http.Request) {

	var jobsCount = 0

	for {
		// Read message from sender
		var msg NotificationDTO.MessageDetection
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			conn.Close()
			break
			// break out of loop
		}

		go BroadcastToReceiver(msg)

		jobsCount = jobsCount + 1
		log.Println("[NOTI] send job message:", msg.JobId, msg.DeviceId, msg.TimeStamp, jobsCount)

	}
}

func handleSender(w http.ResponseWriter, r *http.Request) {
	conn, err := senderUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("[NOTI] -- -- -- -- -- -- ")
	go senderHandler(conn, r)
}

func main() {

	// test Send
	// SendLineNotify("DeepZoom")

	http.HandleFunc("/send", handleSender)

	log.Println("[NOTI] Starting server on port 8006")

	err := http.ListenAndServe(":8006", nil)
	if err != nil {
		log.Println(err)
		return
	}

}
