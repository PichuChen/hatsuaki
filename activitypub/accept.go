package activitypub

import (
	"log/slog"

	"github.com/pichuchen/hatsuaki/datastore/actor"
)

func SendAccept(senderActor *actor.Actor, recevierActorID string, followActiveObjectID string) {
	slog.Info("SendAccept", "sender", senderActor.GetUsername(), "receiver", recevierActorID, "object", followActiveObjectID)

	acceptActive := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Accept",
		"actor":    senderActor.GetFullID(),
		"object":   followActiveObjectID,
	}

	// the object is transient, in which case the id MAY be omitted
	// acceptActive["id"] = senderActor.GetFullID() + "/outbox/" + time.Now().String()

	// append to the sender's outbox
	// senderActor.AppendOutboxObject(acceptActive["id"].(string))

	SendActivity(senderActor.GetUsername(), recevierActorID, acceptActive)

}
