package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	MsgChoke         = 0
	MsgUnchoke       = 1
	MsgInterested    = 2
	MsgNotInterested = 3
	MsgHave          = 4
	MsgBitfield      = 5
	MsgRequest       = 6
	MsgPiece         = 7
	MsgCancel        = 8
)

type Message struct {
	ID      byte
	Payload []byte
}

func buildMessage(id byte, payload []byte) []byte {
	length := uint32(len(payload) + 1) // +1 for the ID
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, length)
	buf.WriteByte(id)
	buf.Write(payload)
	return buf.Bytes()
}

func parseMessage(conn net.Conn) (messageID byte, payload []byte, err error) {
	//Message length (4 bytes, big-endian)
	lengthBytes := make([]byte, 4)
	_, err = io.ReadFull(conn, lengthBytes)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read message length: %w", err)
	}
	messageLength := binary.BigEndian.Uint32(lengthBytes)

	if messageLength == 0 {
		// Keep-alive message - specification says to ignore.
		return 0, nil, nil
	}

	//Message ID (1 byte)
	idBytes := make([]byte, 1)
	_, err = io.ReadFull(conn, idBytes)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read message ID: %w", err)
	}
	messageID = idBytes[0]

	//payload (messageLength - 1 bytes)
	payloadLength := int(messageLength - 1)
	payload = make([]byte, payloadLength)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return messageID, payload, nil
}
