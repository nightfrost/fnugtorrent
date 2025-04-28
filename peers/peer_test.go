package peers

import "testing"

func TestGeneratePeerID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Generate Peer ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePeerID()
			if len(got) != 20 {
				t.Errorf("GeneratePeerID() length = %v, want %v", len(got), 20)
			}
		})
	}
}
