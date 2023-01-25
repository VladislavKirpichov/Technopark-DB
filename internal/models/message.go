package models

import "strings"

type Message struct {
	ErrorCode   int    `json:"-"`
	Msg         string `json:"message"`
	TextDetails string `json:"message_details,omitempty"`
}

func (msg *Message) Error() string {
	return strings.Join([]string{msg.Msg, msg.TextDetails}, " ")
}

func (msg *Message) Code() int {
	return msg.ErrorCode
}

func (msg *Message) SetTextDetails(text string) *Message {
	msg.TextDetails = text
	return msg
}
