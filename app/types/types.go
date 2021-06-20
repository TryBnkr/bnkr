package types

import "github.com/MohammedAl-Mahdawi/bnkr/utils/forms"

// MsgResponse defined the message payload
type MsgResponse struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Theme           string
	IsAuthenticated int
	UserId          uint
	UserName        string
	Form            *forms.Form
}

type MailData struct {
	To       []string
	From     string
	Subject  string
	Content  string
	Template string
}
