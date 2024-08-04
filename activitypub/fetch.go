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

func FetchObject(id string, actorUsername string, sign bool) (map[string]interface{}, error) {

	reqURL := id
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
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
		a, err := actor.FindActorByUsername(actorUsername)
		if err != nil {
			slog.Error("FetchObject", "error", err)
			return nil, err
		}
		keyID := fmt.Sprintf("%s#main-key", a.GetFullID())
		signature.Signature(a.GetPrivateKey(), keyID, req)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("FetchObject", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Read response body failed", "Error", err)
		return nil, err
	}

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	if err != nil {
		slog.Error("Unmarshal response body failed", "Error", err, "body", string(respByte))
		return nil, err
	}

	slog.Info("fetch object success", "object", id, "Response", respMap)

	errorStr, ok := respMap["error"]
	if ok {
		if strings.Contains(errorStr.(string), "Request not signed") {
			slog.Info("request not signed, retry with signature")
			return FetchObject(id, actorUsername, true)
		}
		slog.Error("fetch object failed", "object", id, "Error", errorStr)
		return nil, errors.New("get actor failed")
	}

	return respMap, nil
}
