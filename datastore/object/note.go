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
		"id":        "https://" + config.GetDomain() + "/.activitypub/object/" + id,
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

func (o *Object) GetType() string {
	return (*o)["type"].(string)
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

func (o *Object) GetInReplyTo() string {
	s, ok := (*o)["inReplyTo"].(string)
	if !ok {
		return ""
	}
	return s
}

func (o *Object) GetTag() []map[string]string {
	l, ok := (*o)["tag"].([]map[string]string)
	if !ok {
		return []map[string]string{}
	}
	return l
}

func (o *Object) AddTag(name, url string) {
	list, ok := (*o)["tag"].([]map[string]string)
	if !ok {
		list = []map[string]string{}
	}
	for _, t := range list {
		if t["name"] == name {
			return
		}
	}
	(*o)["tag"] = append(list, map[string]string{"name": name, "href": url, "type": "Mention"})
}

func (o *Object) SetInReplyTo(inReplyTo string) {
	(*o)["inReplyTo"] = inReplyTo
}

func (o *Object) GetTo() []string {
	l, ok := (*o)["to"].([]string)
	if !ok {
		return []string{}
	}
	return l
}

func (o *Object) AddTo(to string) {
	list, ok := (*o)["to"].([]string)
	if !ok {
		list = []string{}
	}
	for _, t := range list {
		if t == to {
			return
		}
	}
	(*o)["to"] = append(list, to)
}

func (o *Object) GetBto() []string {
	l, ok := (*o)["bto"].([]string)
	if !ok {
		return []string{}
	}
	return l
}

func (o *Object) AddBto(bto string) {
	list, ok := (*o)["bto"].([]string)
	if !ok {
		list = []string{}
	}
	for _, b := range list {
		if b == bto {
			return
		}
	}
	(*o)["bto"] = append(list, bto)
}

func (o *Object) GetCC() []string {
	l, ok := (*o)["cc"].([]string)
	if !ok {
		return []string{}
	}
	return l
}

func (o *Object) AddCC(cc string) {
	list, ok := (*o)["cc"].([]string)
	if !ok {
		list = []string{}
	}
	for _, c := range list {
		if c == cc {
			return
		}
	}
	(*o)["cc"] = append(list, cc)
}

func (o *Object) GetBCC() []string {
	l, ok := (*o)["bcc"].([]string)
	if !ok {
		return []string{}
	}
	return l
}

func (o *Object) AddBCC(bcc string) {
	list, ok := (*o)["bcc"].([]string)
	if !ok {
		list = []string{}
	}
	for _, b := range list {
		if b == bcc {
			return
		}
	}
	(*o)["bcc"] = append(list, bcc)
}

func (o *Object) GetAudience() []string {
	l, ok := (*o)["audience"].([]string)
	if !ok {
		return []string{}
	}
	return l
}

func (o *Object) AddAudience(audience string) {
	list, ok := (*o)["audience"].([]string)
	if !ok {
		list = []string{}
	}
	for _, a := range list {
		if a == audience {
			return
		}
	}
	(*o)["audience"] = append(list, audience)
}
