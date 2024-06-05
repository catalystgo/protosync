package downloader

import (
	"fmt"
	"testing"

	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/protosync/internal/downloader/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGitlab(t *testing.T) {
	t.Parallel()

	errDummy := fmt.Errorf("dummy error")

	testCases := []struct {
		name    string
		file    *domain.File
		prepare func(t *testing.T, c *mock.MockhttpClient)
		check   func(t *testing.T, content []byte, err error)
	}{
		{
			name: "success",
			file: &domain.File{
				Domain: "gitlab.com",
				User:   "user",
				Repo:   "repo",
				Path:   "path/to/file/hello.proto",
				Ref:    "08c4336",
			},
			prepare: func(t *testing.T, c *mock.MockhttpClient) {
				c.EXPECT().Get("https://gitlab.com/user/repo/-/raw/08c4336/path/to/file/hello.proto").Return([]byte("hello"), nil)
			},
			check: func(t *testing.T, content []byte, err error) {
				require.NoError(t, err)
				require.Equal(t, []byte("hello"), content)
			},
		},
		{
			name: "error",
			file: &domain.File{
				Domain: "gitlab.com",
				User:   "user",
				Repo:   "repo",
				Path:   "path/to/file/hello.proto",
				Ref:    "master",
			},
			prepare: func(t *testing.T, c *mock.MockhttpClient) {
				c.EXPECT().Get("https://gitlab.com/user/repo/-/raw/master/path/to/file/hello.proto").Return(nil, errDummy)
			},
			check: func(t *testing.T, content []byte, err error) {
				require.ErrorIs(t, err, errDummy)
				require.Nil(t, content)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
		
			ctrl := gomock.NewController(t)

			httpClient := mock.NewMockhttpClient(ctrl)
			tc.prepare(t, httpClient)

			d := NewGitlab(httpClient)
			content, err := d.GetFile(tc.file)

			tc.check(t, content, err)
		})
	}
}
