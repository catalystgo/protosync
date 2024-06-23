package service

import (
	"encoding/base64"
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

			err := s.Download(tc.file, tc.ourDir, "")

			tc.check(t, err)
		})
	}
}

func TestServiceGetConfig(t *testing.T) {
	t.Parallel()

	file := "github.com/catalystgo/proto/experimental/proto/echo.proto"

	// base64 encoded content of the default config file, used
	// to avoid reading the file from the disk
	base64content := `ZGlyZWN0b3J5OiAidmVuZG9yLnByb3RvIgoKZGVwZW5kZW5jaWVzOgogICMgRXhhbXBsZSBzZXJ2ZXIKICAtIHNvdXJjZTogZ2l0aHViLmNvbS9jYXRhbHlzdGdvL3Byb3Rvc3luYy9leGFtcGxlL3Byb3RvL2VjaG8ucHJvdG9AMTIwMzMyYgoKICAjIEdvb2dsZSBwcm90b2J1ZgogIC0gcGF0aDogZ29vZ2xlL3Byb3RvYnVmCiAgICBzb3VyY2VzOgogICAgLSBnaXRodWIuY29tL3Byb3RvY29sYnVmZmVycy9wcm90b2J1Zi9zcmMvZ29vZ2xlL3Byb3RvYnVmL2FueS5wcm90b0BhOTc4Yjc1CiAgICAtIGdpdGh1Yi5jb20vcHJvdG9jb2xidWZmZXJzL3Byb3RvYnVmL3NyYy9nb29nbGUvcHJvdG9idWYvZW1wdHkucHJvdG9AYTk3OGI3NQogICAgLSBnaXRodWIuY29tL3Byb3RvY29sYnVmZmVycy9wcm90b2J1Zi9zcmMvZ29vZ2xlL3Byb3RvYnVmL3N0cnVjdC5wcm90b0BhOTc4Yjc1CiAgICAtIGdpdGh1Yi5jb20vcHJvdG9jb2xidWZmZXJzL3Byb3RvYnVmL3NyYy9nb29nbGUvcHJvdG9idWYvdGltZXN0YW1wLnByb3RvQGE5NzhiNzUKICAgIC0gZ2l0aHViLmNvbS9wcm90b2NvbGJ1ZmZlcnMvcHJvdG9idWYvc3JjL2dvb2dsZS9wcm90b2J1Zi93cmFwcGVycy5wcm90b0BhOTc4Yjc1CgogICMgR29vZ2xlIEFQSQogIC0gcGF0aDogZ29vZ2xlL2FwaQogICAgc291cmNlczoKICAgIC0gZ2l0aHViLmNvbS9nb29nbGVhcGlzL2dvb2dsZWFwaXMvZ29vZ2xlL2FwaS9hbm5vdGF0aW9ucy5wcm90b0BmNjVhZDVmCiAgICAtIGdpdGh1Yi5jb20vZ29vZ2xlYXBpcy9nb29nbGVhcGlzL2dvb2dsZS9hcGkvaHR0cC5wcm90b0BmNjVhZDVmCgogICMgT3BlbkFQSSB2MgogIC0gcGF0aDogcHJvdG9jLWdlbi1vcGVuYXBpdjIvb3B0aW9ucwogICAgc291cmNlczoKICAgICAgLSBnaXRodWIuY29tL2dycGMtZWNvc3lzdGVtL2dycGMtZ2F0ZXdheS9wcm90b2MtZ2VuLW9wZW5hcGl2Mi9vcHRpb25zL2Fubm90YXRpb25zLnByb3RvQDY3MDc0OTUKICAgICAgLSBnaXRodWIuY29tL2dycGMtZWNvc3lzdGVtL2dycGMtZ2F0ZXdheS9wcm90b2MtZ2VuLW9wZW5hcGl2Mi9vcHRpb25zL29wZW5hcGl2Mi5wcm90b0A2NzA3NDk1Cg==`
	content, err := base64.StdEncoding.DecodeString(base64content)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		file    string
		prepare func(w *mock.MockWriter)
		check   func(t *testing.T, err error)
	}{
		{
			name: "success",
			file: file,
			prepare: func(w *mock.MockWriter) {
				w.EXPECT().Write(file, content).Return(nil)
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "write error",
			file: file,
			prepare: func(w *mock.MockWriter) {
				w.EXPECT().Write(file, content).Return(fmt.Errorf("some error occurred"))
			},
			check: func(t *testing.T, err error) {
				require.EqualError(t, err, "some error occurred")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			writer := mock.NewMockWriter(ctrl)

			s := New(nil)
			s.writer = writer

			tc.prepare(writer)

			err := s.GenConfig(tc.file)

			tc.check(t, err)
		})
	}
}
