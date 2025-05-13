package trackers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackpal/bencode-go"
	"nightfrost.com/fnugtorrent/models"
)

func BuildInitialTrackerRequest(trackerUrl string, infoHash string, peerID string, port int, event string) (string, error) {
	request, err := url.Parse(trackerUrl)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("info_hash", url.QueryEscape(infoHash))
	params.Set("peer_id", url.QueryEscape(peerID))
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
	var trackerUrlAsHttps = strings.Replace(trackerRequestUrl, "http://", "https://", 1)
	response, err := http.Get(trackerUrlAsHttps)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	fmt.Println("raw response:", response)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func ProcessTrackerResponse(responseBody []byte) (models.TrackerResponse, error) {
	fmt.Printf("Raw Tracker Response:\n%s\n", string(responseBody))
	var result models.TrackerResponse
	err := bencode.Unmarshal(bytes.NewReader(responseBody), &result)
	if err != nil {
		return models.TrackerResponse{}, err
	}
	return result, nil
}
