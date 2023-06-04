package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateEmptyDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "")
	require.Nil(t, err)

	return dir
}

func CreateEmptyFile(t *testing.T, dir string) string {
	file := filepath.Join(dir, uuid.NewString())
	err := os.WriteFile(file, []byte(""), os.ModePerm)
	require.Nil(t, err)

	return file
}

func GetDirEntriesCount(t *testing.T, dir string) int {
	entries, err := os.ReadDir(dir)
	require.Nil(t, err)

	return len(entries)
}

func GetFileContent(t *testing.T, file string) string {
	b, err := os.ReadFile(file)
	require.Nil(t, err)

	return string(b)
}

func GetFileSize(t *testing.T, file string) int64 {
	f, err := os.Stat(file)
	require.Nil(t, err)

	return f.Size()
}
