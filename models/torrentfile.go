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
	Private     int              `bencode:"private"`
	Files       []FileDictionary `bencode:"files"`
}

type FileDictionary struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

type PeerInfo struct {
	PeerID string `bencode:"peer id"`
	IP     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

type TrackerResponse struct {
	Interval      int        `bencode:"interval"`
	Peers         []PeerInfo `bencode:"peers"`
	FailureReason string     `bencode:"failure reason"`
}

type Message struct {
	ID      byte
	Payload []byte
}
