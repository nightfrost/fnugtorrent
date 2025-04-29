package peers

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"nightfrost.com/fnugtorrent/messages"
	"nightfrost.com/fnugtorrent/models"
	"nightfrost.com/fnugtorrent/utils"
)

const base32alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
const PieceLength = 16384 // 2^14, standard request size [cite: 93]

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
	pieceLength := infoDict["piece length"].(int)
	piecesHashes := infoDict["pieces"].(string)
	numPieces := len(piecesHashes) / 20

	outputFile, err := utils.CreateOutPutFile(infoDict)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	downloadedPieces := make([]bool, numPieces)
	var downloadedCount int
	var mu sync.Mutex

	for downloadedCount < numPieces {
		pieceIndex := utils.ChoosePiece(availablePieces, downloadedPieces)
		if pieceIndex == -1 {
			fmt.Println("No piece available - waiting and trying again.")
			time.Sleep(5 * time.Second)
			continue
		}

		err := downloadPiece(conn, pieceIndex, pieceLength, totalLength, piecesHashes, outputFile, &mu)
		if err != nil {
			fmt.Printf("Error downloading piece: %d : %d", pieceIndex, err)
			return
		}

		mu.Lock()
		downloadedPieces[pieceIndex] = true
		downloadedCount++
		mu.Unlock()

		fmt.Printf("Downloaded piece %d\n", pieceIndex)
	}

	fmt.Println("Download complete!")
}

func downloadPiece(conn net.Conn, pieceIndex int, pieceLength int, totalLength int, piecesHashes string, outputFile *os.File, mu *sync.Mutex) error {
	begin := 0
	pieceSize := pieceLength
	if pieceIndex == len(piecesHashes)/20-1 {
		pieceSize = totalLength - pieceLength*pieceIndex
	}

	for begin < pieceSize {
		blockSize := PieceLength
		if pieceSize-begin < PieceLength {
			blockSize = pieceSize - begin
		}

		request := make([]byte, 12)
		binary.BigEndian.PutUint32(request[0:4], uint32(pieceIndex))
		binary.BigEndian.PutUint32(request[4:8], uint32(begin))
		binary.BigEndian.PutUint32(request[8:12], uint32(blockSize))

		requestMsg := messages.BuildMessage(messages.MsgRequest, request)
		_, err := conn.Write(requestMsg)
		if err != nil {
			return fmt.Errorf("error sending request for piece %d, begin %d: %w", pieceIndex, begin, err)
		}

		begin += blockSize
	}

	begin = 0
	pieceData := make([]byte, pieceSize)
	for begin < pieceSize {
		_, payload, err := messages.ParseMessage(conn)
		if err != nil {
			return fmt.Errorf("error receiving piece data for piece %d, begin %d: %w", pieceIndex, begin, err)
		}

		if payload[0] != messages.MsgPiece {
			fmt.Println("Expected piece message, got:", payload[0])
			return fmt.Errorf("expected piece message, got: %d", payload[0])
		}

		pieceIndexReceived := int(binary.BigEndian.Uint32(payload[1:5]))
		blockBegin := int(binary.BigEndian.Uint32(payload[5:9]))
		blockData := payload[9:]

		if pieceIndex != pieceIndexReceived || begin != blockBegin {

			fmt.Printf("Piece index or begin offset mismatch. Expected: %d %d, received: %d %d", pieceIndex, begin, pieceIndexReceived, blockBegin)
			return fmt.Errorf("piece index or begin offset mismatch")
		}

		copy(pieceData[begin:], blockData)
		begin += len(blockData)
	}

	verifyPieceHash(pieceData, pieceIndex, piecesHashes)

	mu.Lock()
	_, err := outputFile.WriteAt(pieceData, int64(pieceIndex*pieceLength))
	mu.Unlock()
	if err != nil {
		return fmt.Errorf("error writing piece %d to file: %w", pieceIndex, err)
	}

	return nil
}

func verifyPieceHash(pieceData []byte, pieceIndex int, piecesHashes string) error {
	hash := sha1.Sum(pieceData)
	expectedHash := piecesHashes[pieceIndex*20 : (pieceIndex+1)*20]

	if hex.EncodeToString(hash[:]) != expectedHash {
		fmt.Println("Hash mismatch for piece", pieceIndex)
		return fmt.Errorf("hash mismatch")
	}
	return nil
}
