package models

type TorrentFile struct {
	Announce string         `bencode:"announce"`
	Info     InfoDictionary `bencode:"info"`
}

type InfoDictionary struct {
	Name        string           `bencode:"name"`
	PieceLength int              `bencode:"piece length"`
	Pieces      string           `bencode:"pieces"`
	Length      int              `bencode:"length"`
	Files       []FileDictionary `bencode:"files"`
}

type FileDictionary struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}
