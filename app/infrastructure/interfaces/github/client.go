package github

import (
	"fmt"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/go-git/go-git/v5"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
	"github.com/rs/xid"
)

// directory in which the github repository will be cloned
// (this should be listed in gitignore)
const rootDir = "scanning"

type GithubClient interface {
	Clone() error
	CleanDirectory()
	GetOutputPath() *string
	Scan()
	GetRules()
}

type githubClient struct {
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

func (c *githubClient) CleanDirectory() {
	if c.OutputPath != nil {
		os.RemoveAll(*c.OutputPath)
	}
}

func (c *githubClient) Clone() error {
	repoID := xid.New().String()
	fullPath := fmt.Sprintf("%v/%v", rootDir, repoID)

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

func New(url string) GithubClient {
	return &githubClient{
		URL: url,
	}
}
