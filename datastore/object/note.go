package object

import (
	"strings"
	"time"

	"github.com/pichuchen/hatsuaki/datastore/config"
)

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-note
func NewNote() *Object {
	id := GenerateUUIDv7()
	note := Object{
		"id":        id,
		"type":      "Note",
		"published": time.Now().Format(time.RFC3339),
	}
	datastore.Store(id, &note)
	return &note
}

func (o *Object) GetFullID() string {
	id := o.GetID()
	if strings.HasPrefix(id, "https://") {
		return id
	}
	return "https://" + config.GetDomain() + "/.activitypub/object/" + id
}

func (o *Object) GetID() string {
	return (*o)["id"].(string)
}

func (o *Object) GetPublished() string {
	return (*o)["published"].(string)
}

func (o *Object) GetContent() string {
	return (*o)["content"].(string)
}

func (o *Object) SetContent(content string) {
	(*o)["content"] = content
}

func (o *Object) GetAttributedTo() string {
	return (*o)["attributedTo"].(string)
}

func (o *Object) SetAttributedTo(actor string) {
	(*o)["attributedTo"] = actor
}
