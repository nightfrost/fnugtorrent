package utils

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
