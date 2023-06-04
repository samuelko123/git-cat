package git_test

import (
	"path/filepath"
	"testing"

	"github.com/samuelko123/git-cat/git"
	"github.com/samuelko123/git-cat/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestInit_EmptyDir(t *testing.T) {
	dir := utils.CreateEmptyDir(t)

	err := git.Init(dir)

	assert.Nil(t, err)
	assert.Equal(t, 6, utils.GetDirEntriesCount(t, dir))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "hooks")))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "info")))
	assert.Equal(t, 2, utils.GetDirEntriesCount(t, filepath.Join(dir, "objects")))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "objects", "info")))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "objects", "pack")))
	assert.Equal(t, 2, utils.GetDirEntriesCount(t, filepath.Join(dir, "refs")))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "refs", "head")))
	assert.Equal(t, 0, utils.GetDirEntriesCount(t, filepath.Join(dir, "refs", "tags")))
	assert.Equal(t, int64(105), utils.GetFileSize(t, filepath.Join(dir, "config")))
	assert.Equal(t, "ref: refs/heads/master", utils.GetFileContent(t, filepath.Join(dir, "HEAD")))
}

func TestInit_NonEmptyDir(t *testing.T) {
	dir := utils.CreateEmptyDir(t)
	utils.CreateEmptyFile(t, dir)

	err := git.Init(dir)

	assert.Nil(t, err)
	assert.Equal(t, 7, utils.GetDirEntriesCount(t, dir))
}

func TestInit_InvalidDir(t *testing.T) {
	dir := utils.CreateEmptyDir(t)
	file := utils.CreateEmptyFile(t, dir)

	err := git.Init(file)

	assert.NotNil(t, err)
	assert.Equal(t, "mkdir "+file+": The system cannot find the path specified.", err.Error())
}

func TestInitWithOptions_DefaultBranch(t *testing.T) {
	dir := utils.CreateEmptyDir(t)

	err := git.InitWithOptions(dir, &git.InitOptions{DefaultBranch: "main"})

	assert.Nil(t, err)
	assert.Equal(t, "ref: refs/heads/main", utils.GetFileContent(t, filepath.Join(dir, "HEAD")))
}
