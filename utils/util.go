package utils

import "os"

func BytesToBitfield(bitfield []byte) []byte {
	return append([]byte{}, bitfield...)
}

func GetTotalLength(infoDict map[string]any) int {
	if length, ok := infoDict["length"].(int64); ok {
		return int(length)
	}

	totalLength := 0
	files := infoDict["files"].([]interface{})
	for _, file := range files {
		fileDict := file.(map[string]any)
		totalLength += int(fileDict["length"].(int64))
	}
	return totalLength
}

func CreateOutPutFile(infoDict map[string]any) (*os.File, error) {
	fileName := infoDict["name"].(string)

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
