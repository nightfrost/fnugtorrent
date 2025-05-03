package bencode

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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

func CalculateInfoHash(infoDict any) (string, error) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, infoDict)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(hash[:]), nil
}
