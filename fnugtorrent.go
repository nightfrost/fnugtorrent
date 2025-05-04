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
	if torrentData == nil {
		fmt.Println("Decoded .torrent file, but it seems to be empty.")
		return
	}
	fmt.Println("Decoded .torrent file successfully. Body:\n", torrentData)

	//calculate info-hash from dict
	infoHash, err := bencode.CalculateInfoHash(torrentData.Info)
	if err != nil {
		fmt.Println("Error calculating info hash:", err)
		return
	}
	fmt.Println("Info Hash calculated. Value:", infoHash)

	//setup tracker request
	peerID := peers.GeneratePeerID()
	port := 6881
	trackerRequestURL, err := trackers.BuildInitialTrackerRequest(torrentData.Announce, infoHash, peerID, port, "started")
	if err != nil {
		fmt.Println("Error processing tracker request:", err)
		return
	}
	fmt.Println("Tracker Request URL: ", trackerRequestURL)

	//process tracker request & response
	trackerResponseBytes, err := trackers.ProcessTrackerRequest(trackerRequestURL)
	if err != nil {
		fmt.Println("Error processing tracker Request:", err)
		return
	}
	trackerResponse, err := trackers.ProcessTrackerResponse(trackerResponseBytes)
	if err != nil {
		fmt.Println("Error processing tracker response:", err)
	}
	fmt.Println("Tracker Response:", trackerResponse)

	fmt.Print("Handle Peer")

	go peers.HandlePeers(trackerResponse.Peers, infoHash, peerID, *torrentData)

	select {}
}
