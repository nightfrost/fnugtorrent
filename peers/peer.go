package peers

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"

	"nightfrost.com/fnugtorrent/messages"
	"nightfrost.com/fnugtorrent/models"
	"nightfrost.com/fnugtorrent/utils"
)

const base32alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

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

// Borrowed from crypto/rand.Text() function. However, we only need 20 bytes.
func GeneratePeerID() string {
	src := make([]byte, 20)
	rand.Read(src)
	for i := range src {
		src[i] = base32alphabet[src[i]%32]
	}
	return string(src)
}

func HandlePeers(peers []models.PeerInfo, infoHash string, peerID string, torrentData map[string]any) {
	for _, peer := range peers {
		go func(peer models.PeerInfo) {
			addr := net.JoinHostPort(peer.IP, fmt.Sprintf("%d", peer.Port))
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println("Error connecting to peer:", addr, err)
				return
			}
			defer conn.Close()

			err = doHandshake(conn, infoHash, peerID)
			if err != nil {
				fmt.Println("Error performing handshake with peer:", addr, err)
				return
			}

			downloadFromPeer(conn, torrentData)
		}(peer)
	}
}

func doHandshake(conn net.Conn, infoHash string, peerID string) error {
	handshake := buildHandshake(infoHash, peerID)
	_, err := conn.Write(handshake)
	if err != nil {
		return err
	}
	_, _, err = readHandshake(conn)
	if err != nil {
		return err
	}

	return err
}

func downloadFromPeer(conn net.Conn, torrentData map[string]any) {
	defer conn.Close()

	interestedMsg := messages.BuildMessage(messages.MsgInterested, []byte{})
	_, err := conn.Write(interestedMsg)
	if err != nil {
		fmt.Println("Error sending interested message:", err)
		return
	}

	msgID, bitfield, err := messages.ParseMessage(conn)
	if err != nil {
		fmt.Println("Error receiving bitfield message:", err)
		return
	}

	if msgID != messages.MsgBitfield {
		fmt.Println("Expected bitfield message, got:", msgID)
		return
	}
	availablePieces := utils.BytesToBitfield(bitfield)

	infoDict := torrentData["info"].(map[string]any)
	totalLength := utils.GetTotalLength(infoDict)
}
