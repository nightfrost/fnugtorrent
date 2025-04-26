package bencode

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/jackpal/bencode-go"
)

func DecodeTorrentFile(filePath string) (map[string]any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]any)

	err = bencode.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CalculateInfoHash(infoDict any) (string, error) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, infoDict)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(hash[:]), nil
}
