package git

import (
	"os"
	"path/filepath"

	"github.com/samuelko123/git-cat/internal/utils"
)

type InitOptions struct {
	DefaultBranch string
}

func Init(dir string) (err error) {
	return InitWithOptions(dir, &InitOptions{})
}

func InitWithOptions(dir string, opts *InitOptions) (err error) {
	defer utils.ReturnError(&err)

	defaultBranch := opts.DefaultBranch
	if defaultBranch == "" {
		defaultBranch = "master"
	}

	config := "" +
		"[core]\n" +
		"\trepositoryformatversion = 0\n" +
		"\tfilemode = false\n" +
		"\tbare = false\n" +
		"\tsymlinks = false\n" +
		"\tignorecase = true\n"

	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "hooks"), os.ModePerm))
	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "info"), os.ModePerm))
	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "objects", "info"), os.ModePerm))
	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "objects", "pack"), os.ModePerm))
	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "refs", "head"), os.ModePerm))
	utils.PanicIfError(os.MkdirAll(filepath.Join(dir, "refs", "tags"), os.ModePerm))

	utils.PanicIfError(os.WriteFile(filepath.Join(dir, "config"), []byte(config), os.ModePerm))
	utils.PanicIfError(os.WriteFile(filepath.Join(dir, "HEAD"), []byte("ref: refs/heads/"+defaultBranch), os.ModePerm))

	return nil
}
