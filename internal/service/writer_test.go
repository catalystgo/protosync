package service

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderWrite(t *testing.T) {
	t.Parallel()

	// Create a temporary directory
	dir := t.TempDir()

	// Create a temporary file
	file := path.Join(dir, "test.txt")

	// Create a temporary file content
	content := []byte("test")

	// Create a new writer
	p := NewWriteProvider()

	// Write the content to the file
	err := p.Write(file, content)
	require.NoError(t, err)

	// Read the file content
	got, err := os.ReadFile(file)
	require.NoError(t, err)

	// Compare the content
	require.Equal(t, content, got)
}
