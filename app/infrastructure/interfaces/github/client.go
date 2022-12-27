package github

import (
	"fmt"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/go-git/go-git/v5"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/rs/xid"
)

type GithubClient interface {
	Clone() error
	ClearnDirectory()
	GetOutputPath() *string
	Scan()
	GetRules()
}

type githubClient struct {
	RootDir    string
	URL        string
	OutputPath *string
}

func (c *githubClient) Scan() {

}
func (c *githubClient) GetRules() {

}
func (c *githubClient) GetOutputPath() *string {
	return c.OutputPath
}

func (c *githubClient) ClearnDirectory() {
	if c.OutputPath != nil {
		os.RemoveAll(*c.OutputPath)
	}
}

func (c *githubClient) Clone() error {
	repoID := xid.New().String()
	fullPath := fmt.Sprintf("%v/%v", c.RootDir, repoID)

	_, err := git.PlainClone(fullPath, false, &git.CloneOptions{
		URL:      c.URL,
		Progress: os.Stdout,
	})
	if err != nil {
		return gerrors.NewInternalError(err)
	}

	c.OutputPath = pointer.ToString(fullPath)
	return nil

}

func New(rootDir string, url string) GithubClient {
	return &githubClient{
		RootDir: rootDir,
		URL:     url,
	}
}
