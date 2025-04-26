package peers

import (
	"fmt"
	"io"
	"net"
)

func buildHandshake(infoHash string, peerID string) []byte {
	handshake := make([]byte, 68)
	handshake[0] = 19
	copy(handshake[1:20], "BitTorrent protocol")
	copy(handshake[20:28], infoHash)
	copy(handshake[28:48], peerID)
	return handshake
}

func readHandshake(conn net.Conn) (string, string, error) {
	//Length prefix (1)
	lengthPrefix := make([]byte, 1)
	_, err := io.ReadFull(conn, lengthPrefix)
	if err != nil {
		return "", "", fmt.Errorf("failed to length prefix: %w", err)
	}
	if lengthPrefix[0] != 19 {
		return "", "", fmt.Errorf("invalid protocol identifier length: %d", lengthPrefix[0])
	}

	//Protocol Identifier (19)
	protocolIdentifier := make([]byte, 19)
	_, err = io.ReadFull(conn, protocolIdentifier)
	if err != nil {
		return "", "", fmt.Errorf("failed to read protocol identifier: %w", err)
	}
	if string(protocolIdentifier) != "BitTorrent protocol" {
		return "", "", fmt.Errorf("invalid protocol identifier: %s", string(protocolIdentifier))
	}

	//Reserved (8)
	reserved := make([]byte, 8)
	_, err = io.ReadFull(conn, reserved)
	if err != nil {
		return "", "", fmt.Errorf("failed to read reserved bytes: %w", err)
	}

	//Info Hash (20)
	infoHashBytes := make([]byte, 20)
	_, err = io.ReadFull(conn, infoHashBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to read info hash: %w", err)
	}
	infoHash := string(infoHashBytes)

	//Peer ID (20)
	peerIDBytes := make([]byte, 20)
	_, err = io.ReadFull(conn, infoHashBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to read peer ID: %w", err)
	}
	peerID := string(peerIDBytes)

	return infoHash, peerID, nil
}
