package main

import (
	"fmt"

	"nightfrost.com/fnugtorrent/bencode"
	"nightfrost.com/fnugtorrent/peers"
	"nightfrost.com/fnugtorrent/trackers"
)

func main() {
	fmt.Print("Enter .torrent file path: ")
	var filePath string
	fmt.Scanln(&filePath)

	torrentData, err := bencode.DecodeTorrentFile(filePath)
	if err != nil {
		fmt.Println("Error decoding torrent file:", err)
		return
	}

	//calculate info-hash from dict
	infoDict, ok := torrentData["info"]
	if !ok {
		fmt.Println("No 'info' section found in the torrent file")
		return
	}
	infoHash, err := bencode.CalculateInfoHash(infoDict)
	if err != nil {
		fmt.Println("Error calculating info hash:", err)
		return
	}

	//setup tracker request
	trackerURL, ok := torrentData["announce"].(string)
	if !ok {
		fmt.Println("Error: 'announce' not found in torrent file or is not a string")
		return
	}
	peerID := peers.GeneratePeerID()
	port := 6881
	trackerRequestURL, err := trackers.BuildInitialTrackerRequest(trackerURL, infoHash, peerID, port, "started")
	if err != nil {
		fmt.Println("Error processing tracker request:", err)
		return
	}

	//process tracker request & response
	trackerResponseBytes, err := trackers.ProcessTrackerRequest(trackerRequestURL)
	if err != nil {
		fmt.Println("Error processing tracker response:", err)
		return
	}
	trackerResponse, err := trackers.ProcessTrackerResponse(trackerResponseBytes)
	if err != nil {
		fmt.Println("Error processing tracker response:", err)
	}

	go peers.HandlePeers(trackerResponse.Peers, infoHash, peerID, torrentData)

	select {}
}
