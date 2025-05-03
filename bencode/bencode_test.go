package bencode

import (
	"reflect"
	"testing"
)

func TestDecodeTorrentFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name:    "Decode Torrent File",
			args:    args{"C:\\Users\\Nightfrost\\source\\repos\\FnugTorrent\\testdata\\ubuntu-22.04.5-live-server-amd64.iso.torrent"},
			want:    map[string]any{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeTorrentFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeTorrentFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeTorrentFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
