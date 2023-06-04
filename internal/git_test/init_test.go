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
}

func TestInitWithOptions_DefaultBranch(t *testing.T) {
	dir := utils.CreateEmptyDir(t)

	err := git.InitWithOptions(dir, &git.InitOptions{DefaultBranch: "main"})

	assert.Nil(t, err)
	assert.Equal(t, 6, utils.GetDirEntriesCount(t, dir))
	assert.Equal(t, "ref: refs/heads/main", utils.GetFileContent(t, filepath.Join(dir, "HEAD")))
}
