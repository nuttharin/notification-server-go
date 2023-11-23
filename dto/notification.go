package NotificationDTO

type NotificationReq struct {
	Title   string `json:"title"  validate:"required"`
	Message string `json:"message"  validate:"required"`
}

type NotificationRes struct {
	Description string `json:"description" `
}

type NotificationEntity struct {
	Message         string `db:"message"`
	Status          string `db:"status"`
	CreatedBy       string `db:"created_by"`
	CreatedDateTime string `db:"created_datetime"`
	UpdatedBy       string `db:"updated_by" `
	UpdatedDateTime string `db:"updated_datetime"`
}

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
	Email       string `json:"email"`
	// SendType    []string `json:"sendType"` // ["line","email","sms","etc"]
	SendType string `json:"sendType"` // "line","email"

	ConfigName string `json:"configName"`
	RoiName    string `json:"roiName"`
}
