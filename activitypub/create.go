package activitypub

import (
	"log"
	"log/slog"

	"github.com/pichuchen/hatsuaki/datastore/actor"
	"github.com/pichuchen/hatsuaki/datastore/object"
)

func SendCreate(senderActor *actor.Actor, object *object.Object) {
	slog.Info("SendCreate", "sender", senderActor.GetUsername(), "object", object.GetFullID())

	createActivity := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Create",
		"actor":    senderActor.GetFullID(),
		"object":   object,
	}
	receivers := map[string]bool{}

	// 如果有任何 to, bto, cc, bcc, audience 的話，都要加進去
	if to := object.GetTo(); len(to) > 0 {
		createActivity["to"] = to
		for _, to := range to {
			receivers[to] = true
		}
	}

	if bto := object.GetBto(); len(bto) > 0 {
		createActivity["bto"] = bto
		for _, bto := range bto {
			receivers[bto] = true
		}
	}

	if cc := object.GetCC(); len(cc) > 0 {
		createActivity["cc"] = cc
		for _, cc := range cc {
			receivers[cc] = true
		}
	}

	if bcc := object.GetBCC(); len(bcc) > 0 {
		createActivity["bcc"] = bcc
		for _, bcc := range bcc {
			receivers[bcc] = true
		}
	}

	if audience := object.GetAudience(); len(audience) > 0 {
		createActivity["audience"] = audience
		for _, audience := range audience {
			receivers[audience] = true
		}
	}

	createActivity["id"] = object.GetFullID() + "/activity"
	createActivity["published"] = object.GetPublished()

	followers := senderActor.GetFollowerIDs()
	followerMap := map[string]bool{}
	for _, followerID := range followers {
		followerMap[followerID] = true
	}
	sendCnt := 0
	for recevierActorID := range receivers {
		if recevierActorID == senderActor.GetFullID() {
			continue
		}
		if recevierActorID == object.GetAttributedTo() {
			continue
		}
		if _, ok := followerMap[recevierActorID]; ok {
			continue
		}
		if recevierActorID == "https://www.w3.org/ns/activitystreams#Public" {
			continue
		}
		if recevierActorID == object.GetAttributedTo()+"/followers" {
			// 轉傳給所有的 followers
			for recevierActorID, _ := range followerMap {
				go SendActivity(senderActor.GetUsername(), recevierActorID, createActivity)
				sendCnt++
			}
			continue
		}
		go SendActivity(senderActor.GetUsername(), recevierActorID, createActivity)
		sendCnt++
	}
	log.Println("SendCreate", "receivers", receivers, "sendCnt", sendCnt)
}
