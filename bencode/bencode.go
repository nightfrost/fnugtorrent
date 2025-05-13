package bencode

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
	"nightfrost.com/fnugtorrent/models"
)

func DecodeTorrentFile(filePath string) (*models.TorrentFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var torrentFile models.TorrentFile

	err = bencode.Unmarshal(file, &torrentFile)
	if err != nil {
		return nil, err
	}
	return &torrentFile, nil
}

func CalculateInfoHash(info models.InfoDictionary) (string, error) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, info)

	if err != nil {
		return "", fmt.Errorf("bencode marshal error: %w", err)
	}

	hash := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(hash[:]), nil
}
