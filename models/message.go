package models

import "github.com/zebresel-com/mongodm"

type Message struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Sender   interface{} `json:"sender"  bson:"sender" model:"User" relation:"11" required:"true"`
	Receiver interface{} `json:"receiver"  bson:"receiver"  model:"User" relation:"1n" required:"true"`
	Text     string      `json:"text"  bson:"text" minLen:"1" maxLen:"500" required:"true"`
}
