package utils

import (
	"fmt"
	"os"
	"reflect"

	"nightfrost.com/fnugtorrent/models"
)

func BytesToBitfield(bitfield []byte) []byte {
	return append([]byte{}, bitfield...)
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

func PrintModel(model any) {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	fmt.Printf("--- %s ---\n", t.Name())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag.Get("bencode")

		fmt.Printf("%s (%s): ", fieldName, tag)

		switch field.Kind() {
		case reflect.String:
			fmt.Println(field.String())
		case reflect.Int:
			fmt.Println(field.Int())
		case reflect.Slice:
			fmt.Println("[")
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Kind() == reflect.Struct {
					PrintModel(field.Index(j).Interface())
				} else {
					fmt.Printf("  - %v\n", field.Index(j).Interface())
				}
			}
			fmt.Println("]")
		case reflect.Struct:
			fmt.Println("{")
			PrintModel(field.Interface())
			fmt.Println("}")
		default:
			fmt.Printf("Unsupported type: %v\n", field.Kind())
		}
	}
	fmt.Println("--- End ---")
}
