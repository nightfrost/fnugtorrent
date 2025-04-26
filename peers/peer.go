package main

import "net"

func buildHandshake(infoHash string, peerID string) []byte {
	handshake := make([]byte, 68)
	handshake[0] = 19
	copy(handshake[1:20], "BitTorrent protocol")
	copy(handshake[20:28], infoHash)
	copy(handshake[28:48], peerID)
	return handshake
}

func readHandshake(conn net.Conn) (string, string, error) {
	handshake := make([]byte, 68)
	_, err := conn.Read(handshake)
	if err != nil {
		return "", "", err
	}
	if string(handshake[:19]) != "BitTorrent protocol" {
		return "", "", err
	}
	infoHash := string(handshake[28:48])
	if infoHash == "" {
		return "", "", err
	}
	peerID := string(handshake[48:68])
	if peerID == "" {
		return "", "", err
	}
	return infoHash, peerID, nil
}
