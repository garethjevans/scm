package cmd

import (
	"os"

	"github.com/garethjevans/scm/pkg/client"
	"github.com/jenkins-x/go-scm/scm"
)

var (
	Path  string
	Host  string
	Owner string
	Repo  string
	Kind  string
)

func GetScmClient(repositoryURL string) (*scm.Client, string, string, error) {
	b, err := os.ReadFile(Path)
	if err != nil {
		return nil, "", "", err
	}

	scmClient, username, token, err := client.FromRepoURL(repositoryURL, string(b), Kind)
	return scmClient, username, token, err
}
