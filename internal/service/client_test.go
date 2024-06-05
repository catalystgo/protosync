package service

import (
	"fmt"
	"testing"

	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/protosync/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	t.Parallel()

	var (
		defaulContent  = []byte(`hello world`)
		defaultOutDir  = "proto"
		defaultFileURL = "github.com/catalystgo/protosync/proto/hello.proto@559bac64"
		defaultFile    = &domain.File{
			Domain: "github.com",
			User:   "catalystgo",
			Repo:   "protosync",
			Path:   "proto/hello.proto",
			Ref:    "559bac64",
		}

		errDummy = fmt.Errorf("dummy error")
	)

	testCases := []struct {
		name    string
		file    string
		ourDir  string
		content []byte

		prepare func(w *mock.MockWriter, d *mock.MockDownloader)
		check   func(t *testing.T, err error)
	}{
		{
			name:    "success",
			file:    defaultFileURL,
			ourDir:  defaultOutDir,
			content: defaulContent,
			prepare: func(w *mock.MockWriter, d *mock.MockDownloader) {

				w.EXPECT().Write("proto/github.com/catalystgo/protosync/proto/hello.proto", defaulContent).Return(nil)
				d.EXPECT().GetFile(defaultFile).Return(defaulContent, nil)
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:    "download error",
			file:    defaultFileURL,
			ourDir:  defaultOutDir,
			content: defaulContent,
			prepare: func(_ *mock.MockWriter, d *mock.MockDownloader) {
				d.EXPECT().GetFile(defaultFile).Return(nil, errDummy)
			},
			check: func(t *testing.T, err error) {
				require.ErrorIs(t, err, errDummy)
			},
		},
		{
			name:    "write error",
			file:    defaultFileURL,
			ourDir:  defaultOutDir,
			content: defaulContent,
			prepare: func(w *mock.MockWriter, d *mock.MockDownloader) {
				w.EXPECT().Write("proto/github.com/catalystgo/protosync/proto/hello.proto", defaulContent).Return(errDummy)
				d.EXPECT().GetFile(defaultFile).Return(defaulContent, nil)
			},
			check: func(t *testing.T, err error) {
				require.ErrorIs(t, err, errDummy)
			},
		},
		{
			name:    "invalid file",
			file:    "gitlab.com/catalystgo/protosync/proto/hello@559bac64",
			ourDir:  defaultOutDir,
			content: defaulContent,
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name:    "unknown domain",
			file:    "unknown.com/catalystgo/protosync/proto/hello.proto@559bac64",
			ourDir:  defaultOutDir,
			content: defaulContent,
			check: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			downloader := mock.NewMockDownloader(ctrl)
			writer := mock.NewMockWriter(ctrl)
			s := New(writer)

			// Use the same downloader for all domains for simplicity

			s.Register("github.com", downloader)
			s.Register("gitlab.com", downloader)
			s.Register("bitbucket.org", downloader)

			if tc.prepare != nil {
				tc.prepare(writer, downloader)
			}

			err := s.Download(tc.file, tc.ourDir)

			tc.check(t, err)
		})
	}
}
