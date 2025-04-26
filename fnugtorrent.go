package main

import (
	"fmt"

	"nightfrost.com/fnugtorrent/bencode"
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
}
