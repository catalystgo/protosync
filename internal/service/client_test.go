package service

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
				w.EXPECT().Write("proto/github.com/catalystgo/protosync/proto/hello.proto", defaulContent, true).Return(nil)
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
				w.EXPECT().Write("proto/github.com/catalystgo/protosync/proto/hello.proto", defaulContent, true).Return(errDummy)
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

			err := s.Download(tc.file, tc.ourDir, "")

			tc.check(t, err)
		})
	}
}

func TestServiceGetConfig(t *testing.T) {
	t.Parallel()

	getFile := func() string {
		return fmt.Sprintf("github.com/catalystgo/proto/experimental/proto/%s.proto", gofakeit.DigitN(20))
	}

	// base64 encoded content of the default config file, used
	// to avoid reading the file from the disk
	base64content := `ZGlyZWN0b3J5OiAidmVuZG9yLnByb3RvIgoKZGVwZW5kZW5jaWVzOgogICMgRXhhbXBsZSBzZXJ2ZXIKICAjIC0gc291cmNlOiBnaXRodWIuY29tL2NhdGFseXN0Z28vcHJvdG9zeW5jL2V4YW1wbGUvcHJvdG8vZWNoby5wcm90b0AxMjAzMzJiCgogICMgR29vZ2xlIHByb3RvYnVmCiAgLSBwYXRoOiBnb29nbGUvcHJvdG9idWYKICAgIHNvdXJjZXM6CiAgICAtIGdpdGh1Yi5jb20vcHJvdG9jb2xidWZmZXJzL3Byb3RvYnVmL3NyYy9nb29nbGUvcHJvdG9idWYvYW55LnByb3RvQGE5NzhiNzUKICAgIC0gZ2l0aHViLmNvbS9wcm90b2NvbGJ1ZmZlcnMvcHJvdG9idWYvc3JjL2dvb2dsZS9wcm90b2J1Zi9lbXB0eS5wcm90b0BhOTc4Yjc1CiAgICAtIGdpdGh1Yi5jb20vcHJvdG9jb2xidWZmZXJzL3Byb3RvYnVmL3NyYy9nb29nbGUvcHJvdG9idWYvc3RydWN0LnByb3RvQGE5NzhiNzUKICAgIC0gZ2l0aHViLmNvbS9wcm90b2NvbGJ1ZmZlcnMvcHJvdG9idWYvc3JjL2dvb2dsZS9wcm90b2J1Zi90aW1lc3RhbXAucHJvdG9AYTk3OGI3NQogICAgLSBnaXRodWIuY29tL3Byb3RvY29sYnVmZmVycy9wcm90b2J1Zi9zcmMvZ29vZ2xlL3Byb3RvYnVmL3dyYXBwZXJzLnByb3RvQGE5NzhiNzUKCiAgIyBHb29nbGUgQVBJCiAgLSBwYXRoOiBnb29nbGUvYXBpCiAgICBzb3VyY2VzOgogICAgLSBnaXRodWIuY29tL2dvb2dsZWFwaXMvZ29vZ2xlYXBpcy9nb29nbGUvYXBpL2Fubm90YXRpb25zLnByb3RvQGY2NWFkNWYKICAgIC0gZ2l0aHViLmNvbS9nb29nbGVhcGlzL2dvb2dsZWFwaXMvZ29vZ2xlL2FwaS9odHRwLnByb3RvQGY2NWFkNWYKCiAgIyBPcGVuQVBJIHYyCiAgLSBwYXRoOiBwcm90b2MtZ2VuLW9wZW5hcGl2Mi9vcHRpb25zCiAgICBzb3VyY2VzOgogICAgICAtIGdpdGh1Yi5jb20vZ3JwYy1lY29zeXN0ZW0vZ3JwYy1nYXRld2F5L3Byb3RvYy1nZW4tb3BlbmFwaXYyL29wdGlvbnMvYW5ub3RhdGlvbnMucHJvdG9ANjcwNzQ5NQogICAgICAtIGdpdGh1Yi5jb20vZ3JwYy1lY29zeXN0ZW0vZ3JwYy1nYXRld2F5L3Byb3RvYy1nZW4tb3BlbmFwaXYyL29wdGlvbnMvb3BlbmFwaXYyLnByb3RvQDY3MDc0OTUK`
	content, err := base64.StdEncoding.DecodeString(base64content)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		file    string
		prepare func(w *mock.MockWriter, file string)
		check   func(t *testing.T, err error)
	}{
		{
			name: "success",
			file: getFile(),
			prepare: func(w *mock.MockWriter, file string) {
				w.EXPECT().Write(file, content, false).Return(nil)
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "write error",
			file: getFile(),
			prepare: func(w *mock.MockWriter, file string) {
				w.EXPECT().Write(file, content, false).Return(fmt.Errorf("some error occurred"))
			},
			check: func(t *testing.T, err error) {
				require.EqualError(t, err, "some error occurred")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			writer := mock.NewMockWriter(ctrl)

			s := New(nil)
			s.writer = writer

			tc.prepare(writer, tc.file)

			err := s.GenConfig(tc.file)

			tc.check(t, err)
		})
	}
}
