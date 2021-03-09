package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgo_Process(t *testing.T) {
	a := NewArgo()
	tests := []struct {
		name    string
		payload string
		want    string
		wantErr bool
	}{
		{
			name:    "placeholder",
			payload: `{"data": 1, "priority": 1}`,
			want:    `{"data": 1, "priority": 1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Process(tt.payload)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
