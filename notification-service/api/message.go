package api

type Message struct {
	TargetEmail  string `json:"target_email"`
	Subject      string `json:"subject"`
	MessagePlain string `json:"message_plain"`
	MessageHtml  string `json:"message_html"`
}
