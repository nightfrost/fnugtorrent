package trackers

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
	"nightfrost.com/fnugtorrent/models"
)

func BuildInitialTrackerRequest(trackerUrl string, infoHash string, peerID string, port int, event string) (string, error) {
	request, err := url.Parse(trackerUrl)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("info_hash", string(infoHash))
	params.Set("peer_id", peerID)
	params.Set("port", strconv.Itoa(port))
	params.Set("uploaded", strconv.Itoa(0))
	params.Set("downloaded", strconv.Itoa(0))
	params.Set("left", strconv.Itoa(0))
	if event != "" {
		params.Set("event", event)
	}

	request.RawQuery = params.Encode()
	return request.String(), nil
}

func ProcessTrackerRequest(trackerRequestUrl string) ([]byte, error) {
	response, err := http.Get(trackerRequestUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func ProcessTrackerResponse(responseBody []byte) (models.TrackerResponse, error) {
	var result models.TrackerResponse
	err := bencode.Unmarshal(bytes.NewReader(responseBody), &result)
	if err != nil {
		return models.TrackerResponse{}, err
	}
	return result, nil
}
