package domain

import (
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
		{
			name:   "invalid file format (domain missing)",
			file:   "",
			errMsg: "invalid source format:  => missing domain",
		},
		{
			name:   "invalid file format (domain empty)",
			file:   "/catalystgo/protosync/server.proto@master",
			errMsg: "invalid source format: /catalystgo/protosync/server.proto@master => domain is empty",
		},
		{
			name:   "invalid file format (ref missing)",
			file:   "github.com/catalystgo/protosync/server.proto",
			errMsg: "invalid source format: github.com/catalystgo/protosync/server.proto => missing ref",
		},
		{
			name:   "invalid file format (ref empty)",
			file:   "github.com/catalystgo/protosync/server.proto@",
			errMsg: "invalid source format: github.com/catalystgo/protosync/server.proto@ => ref is empty",
		},
		{
			name:   "invalid file format (path missing)",
			file:   "github.com/@master",
			errMsg: "invalid source format: github.com/@master => path is empty",
		},
		{
			name:   "invalid file format (path invalid)", // repo is missing
			file:   "github.com/catalystgo/server.proto@master",
			errMsg: "invalid source format: github.com/catalystgo/server.proto@master => path is invalid",
		},
		{
			name:   "invalid file format (user is empty)",
			file:   "github.com//protosync/server.proto@master",
			errMsg: "invalid source format: github.com//protosync/server.proto@master => user is empty",
		},
		{
			name:   "invalid file format (repo is empty)",
			file:   "github.com/catalystgo//server.proto@master",
			errMsg: "invalid source format: github.com/catalystgo//server.proto@master => repo is empty",
		},
		{
			name:   "invalid file format (extension not .proto)",
			file:   "github.com/catalystgo/protosync/server.txt@master",
			errMsg: "invalid source format: github.com/catalystgo/protosync/server.txt@master => only .proto extension is allowed",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f, err := ParseFile(tt.file)
			if tt.errMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, *f, *tt.fileObj)
		})
	}
}
