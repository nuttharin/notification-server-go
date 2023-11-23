package Linenotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"bytes"
	"net/http"
	NotificationDTO "notification-server/dto"
	Utils "notification-server/utils"

	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
	// "os"
)

var loc, _ = time.LoadLocation("Asia/Bangkok")

func SendLineNotify(message NotificationDTO.MessageDetection) int {
	log.Println("[INFO] : SendLineNotify")
	// message 1000 characters max
	// Return as response StatusCode
	// 200: Success
	// 400: Bad request
	// 401: Invalid access token
	// 500: Failure due to server error
	// Other: Processed over time or stopped

	// godotenv package

	// LINE_NOTIFY_TOKEN := ""
	// LINE_NOTIFY_TOKEN := Utils.GetEnv("LINE_NOTIFY_TOKEN")
	LINE_NOTIFY_TOKEN := message.Line

	client := resty.New()
	var result map[string]string

	// Example send Sticker
	// json.Unmarshal([]byte(`{
	//           "message": "smart energy",
	//           "stickerId": "125",
	//           "stickerPackageId": "1"
	//    }`), &result)

	body := fmt.Sprintf(`{
			"message": "\n Source-name : %v \n Roi-Name : %v \n Rule-base : %v"
         }`, message.ConfigName, message.RoiName, message.Description)

	json.Unmarshal([]byte(body), &result)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+LINE_NOTIFY_TOKEN).
		SetFormData(result).Post("https://notify-api.line.me/api/notify")
	if err != nil {
		log.Fatalf("ERROR LINE Notify API: %s", err)
	}

	log.Println("Line Notify sent successfully ", resp)

	return resp.StatusCode()

}

func SendLineNotifyMsgAndImg(message NotificationDTO.MessageDetection) int {
	log.Println("[INFO] : SendLineNotify")
	// Your LINE Notify token
	LINE_NOTIFY_TOKEN := message.Line
	token := LINE_NOTIFY_TOKEN

	// Your base64 encoded image string
	base64Image := string(message.ImageUrl)

	// Decode the base64 string to []byte
	imageData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64Image))
	if err != nil {
		log.Println("Error decoding base64 data: ", err)
	}

	// resize
	// Write down resolution image
	// Calculate the thumbnail size
	img, _, err := image.Decode(strings.NewReader(string(imageData)))
	if err != nil {
		log.Println("error Decode image to jpeg: ", err)
		return 0
	}

	thumbnailWidth := img.Bounds().Size().X / 3
	thumbnailHeight := img.Bounds().Size().Y / 3
	thumbnail := resize.Resize(uint(thumbnailWidth), uint(thumbnailHeight), img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumbnail, nil)
	if err != nil {
		log.Println("error encoding image to jpeg: ", err)
		return 0
	}

	// Encode the jpeg to base64
	resizedImageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	imageData, err = base64.StdEncoding.DecodeString(strings.TrimSpace(resizedImageBase64))
	if err != nil {
		log.Println("Error decoding resize base64 data: ", err)
		return 0
	}

	// Prepare multipart message
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add message part
	TimeStampInt, _ := strconv.ParseInt(message.TimeStamp, 10, 64)
	// convert milliseconds to seconds
	TimeStampIntSecond := TimeStampInt / 1000
	timeStampStr := time.Unix(TimeStampIntSecond, (TimeStampInt%1000)*int64(time.Millisecond)).In(loc).Format("2006-01-02 15:04:05")
	msgStr := fmt.Sprintf("\nTime : %s\nSource-name : %s\nRoi-Name : %s\nRule-base : %s", timeStampStr, message.ConfigName, message.RoiName, message.Description)

	if err := writer.WriteField("message", msgStr); err != nil {
		fmt.Println("error writing message field: %w", err)
		return 0
	}

	// Add image part
	part, err := writer.CreateFormFile("imageFile", "image.jpg")
	if err != nil {
		log.Println("Error creating form file 'imageFile': ", err)
		return 0

	}
	_, err = part.Write(imageData)
	if err != nil {
		log.Println("Error writing to form file 'imageFile': ", err)
		return 0
	}

	// Close the writer to finalize the multipart message
	err = writer.Close()
	if err != nil {
		log.Println("Error closing writer: ", err)
		return 0
	}

	// Create request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetRequestURI("https://notify-api.line.me/api/notify")
	req.Header.SetContentType(writer.FormDataContentType())
	req.SetBody(body.Bytes())

	// Send request
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = fasthttp.Do(req, resp)
	if err != nil {
		log.Println("Error sending request to LINE Notify: ", err)
		return 0

	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("Error response from LINE Notify: ", resp.StatusCode())
		return 0

	}
	return resp.StatusCode()
}

func SendNotify(msg string) int {
	// Set your Line Notify access token here
	accessToken := Utils.GetEnv("LINE_NOTIFY_TOKEN")

	// URL for Line Notify API
	url := "https://notify-api.line.me/api/notify"

	// Message to send
	message := "Hello, Line Notify!"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a request with the access token and message
	req, err := http.NewRequest("POST", url, bytes.NewBufferString("message="+message))
	if err != nil {
		panic(err)
	}

	// Set the request headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		panic("Failed to send Line Notify message")
	}

	// Notification sent successfully
	println("Line Notify message sent successfully!")
	return 1
}
