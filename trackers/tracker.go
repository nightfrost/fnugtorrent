package main

import (
	"bytes"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type TrackerResponse struct {
	Interval      int        `bencode:"interval"`
	Peers         []PeerInfo `bencode:"peers"`
	FailureReason string     `bencode:"failure reason"`
}

type PeerInfo struct {
	PeerID string `bencode:"peer id"`
	IP     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

func processTrackerRequest(trackerUrl string, infoHash string, peerID string, port int, uploaded int, downloaded int, left int, event string) (string, error) {
	request, err := url.Parse(trackerUrl)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("info_hash", string(infoHash))
	params.Set("peer_id", peerID)
	params.Set("port", strconv.Itoa(port))
	params.Set("uploaded", strconv.Itoa(uploaded))
	params.Set("downloaded", strconv.Itoa(downloaded))
	params.Set("left", strconv.Itoa(left))
	if event != "" {
		params.Set("event", event)
	}

	request.RawQuery = params.Encode()
	return request.String(), nil
}

func processTrackerResponse(responseBody []byte) (TrackerResponse, error) {
	var result TrackerResponse
	err := bencode.Unmarshal(bytes.NewReader(responseBody), &result)
	if err != nil {
		return TrackerResponse{}, err
	}
	return result, nil
}
