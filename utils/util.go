package utils

import (
	"os"

	"nightfrost.com/fnugtorrent/models"
)

func BytesToBitfield(bitfield []byte) []byte {
	return append([]byte{}, bitfield...)
}

func GetTotalLength(infoDict models.InfoDictionary) int {
	if infoDict.Length != 0 {
		return infoDict.Length
	}

	totalLength := 0
	files := infoDict.Files
	for _, file := range files {
		totalLength += file.Length
	}
	return totalLength
}

func CreateOutPutFile(infoDict models.InfoDictionary) (*os.File, error) {
	fileName := infoDict.Name

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ChoosePiece(availablePieces []byte, downloadedPieces []bool) int {
	for i := range downloadedPieces {
		if !downloadedPieces[i] && isPieceAvailable(availablePieces, i) {
			return i
		}
	}
	return -1 // No piece available
}
func isPieceAvailable(availablePieces []byte, index int) bool {
	byteIndex := index / 8
	bitOffset := 7 - (index % 8)
	return (availablePieces[byteIndex]>>bitOffset)&1 == 1
}
