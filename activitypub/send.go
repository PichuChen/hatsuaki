package activitypub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/pichuchen/hatsuaki/activitypub/signature"
	"github.com/pichuchen/hatsuaki/datastore/actor"
)

func SendActivity(senderUsername string, recevierActorID string, activity map[string]interface{}) {
	// 這邊應該要把 activity 送到 recevierActorID 的 inbox

	// 首先要先取的對方的 inbox 位置
	inbox, err := GetInboxByActorID(recevierActorID, false)
	if err != nil {
		slog.Error("GetInboxByActorID failed", "error", err)
		return
	}

	// 然後送出 activity
	activityByte, err := json.Marshal(activity)
	if err != nil {
		slog.Error("Marshal activity failed", "error", err)
		return
	}

	slog.Info("SendActivity", "activity", string(activityByte), "receiver", recevierActorID)

	req, err := http.NewRequest("POST", inbox, strings.NewReader(string(activityByte)))
	if err != nil {
		slog.Error("Create request failed", "error", err)
		return
	}

	req.Header.Set("Content-Type", "application/activity+json")
	req.Header.Set("Accept", "application/activity+json, application/ld+json")

	// Add Date
	gmtTimeLoc := time.FixedZone("GMT", 0)
	s := time.Now().In(gmtTimeLoc).Format(http.TimeFormat)
	req.Header.Add("Date", s)

	// Add Host
	req.Header.Add("Host", req.URL.Host)

	// Add Signature
	senderActor, err := actor.FindActorByUsername(senderUsername)
	if err != nil {
		slog.Error("SendActivity", "error", err)
		return
	}
	keyID := fmt.Sprintf("%s#main-key", senderActor.GetFullID())
	signature.Signature(senderActor.GetPrivateKey(), keyID, req)

	slog.Info("SendActivity", "activity", string(activityByte))

	req.Body = io.NopCloser(strings.NewReader(string(activityByte)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("SendActivity failed", "error", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Read response body failed", "Error", err)
		return
	}

	slog.Info("SendActivity response", "body", string(respBody))

	if resp.StatusCode != http.StatusOK {
		slog.Error("SendActivity failed", "status code", resp.StatusCode)
		return
	}

	slog.Info("SendActivity success", "activity", activity, "receiver", recevierActorID)

}

// GetInboxByActorID 會回傳 actor 的 inbox 位置，
// 如果 sign 是 true 的話，則需要對回傳的位置進行簽章
func GetInboxByActorID(actorID string, sign bool) (string, error) {
	reqURL := actorID
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/activity+json, application/ld+json")
	// Add Date
	gmtTimeLoc := time.FixedZone("GMT", 0)
	s := time.Now().In(gmtTimeLoc).Format(http.TimeFormat)
	req.Header.Add("Date", s)

	// Add Host
	req.Header.Add("Host", req.URL.Host)

	// Add Signature
	if sign {
		instanceActor, err := actor.FindActorByUsername("instance.actor")
		if err != nil {
			slog.Error("GetInboxByActorID", "error", err)
			return "", err
		}
		keyID := fmt.Sprintf("%s#main-key", instanceActor.GetFullID())
		signature.Signature(instanceActor.GetPrivateKey(), keyID, req)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("GetInboxByActorID", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Read response body failed", "Error", err)
		return "", err
	}

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	if err != nil {
		slog.Error("Unmarshal response body failed", "Error", err, "body", string(respByte))
		return "", err
	}

	slog.Info("Get actor success", "actor", actorID, "Response", respMap)

	errorStr, ok := respMap["error"]
	if ok {
		if strings.Contains(errorStr.(string), "Request not signed") {
			slog.Info("request not signed, retry with signature")
			return GetInboxByActorID(actorID, true)
		}
		slog.Error("Get actor failed", "actor", actorID, "Error", errorStr)
		return "", errors.New("get actor failed")
	}

	// 如果有 sharedInbox 的話，就優先回傳 sharedInbox
	sharedInbox, ok := respMap["endpoints"].(map[string]interface{})["sharedInbox"]
	if ok {
		sharedInboxStr, ok := sharedInbox.(string)
		if ok {
			return sharedInboxStr, nil
		}
	}

	inbox, ok := respMap["inbox"]
	if !ok {
		slog.Error("No inbox in actor", "actor", actorID)
		return "", errors.New("no inbox in actor")
	}

	inboxStr, ok := inbox.(string)
	if !ok {
		slog.Error("Inbox is not string", "inbox", inbox)
		return "", errors.New("inbox is not string")
	}

	return inboxStr, nil

}
