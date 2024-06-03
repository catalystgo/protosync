package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    string
		errMsg  string
		fileObj *File
	}{
		{
			name:   "empty file",
			file:   "",
			errMsg: "invalid source format:  => missing domain",
		},
		{
			name:   "invalid file format (ref missing)",
			file:   "github.com/catalystgo/protosync/server.proto",
			errMsg: "invalid source format: github.com/catalystgo/protosync/server.proto => missing ref",
		},
		{
			name:   "invalid file format (extension not .proto)",
			file:   "github.com/catalystgo/protosync/server.txt@master",
			errMsg: "invalid source format: github.com/catalystgo/protosync/server.txt@master => only .proto extension is allowed",
		},
		{
			name:   "valid file format",
			file:   "github.com/catalystgo/protosync/server.proto@master",
			errMsg: "",
			fileObj: &File{
				Domain: "github.com",
				User:   "catalystgo",
				Repo:   "protosync",
				Path:   "server.proto",
				Ref:    "master",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			f, err := ParseFile(tt.file)
			if tt.errMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(f, tt.fileObj))
		})
	}
}
