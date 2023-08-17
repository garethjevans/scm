package client

import (
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"net/url"
	"strings"
)

// FromRepoURL parses a URL of the form https://:authtoken@host/ and attempts to
// determine the driver and creates a client to authenticate to the endpoint.
func FromRepoURL(repoURL string, credentials string) (*scm.Client, string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, "", err
	}

	auth := ""
	if password, ok := u.User.Password(); ok {
		auth = password
	} else {
		fmt.Println("[DEBUG] Token is not available from the url, falling back to .git-credentials")
		token, err := DetermineToken(credentials, repoURL)
		if err != nil {
			return nil, "", err
		}
		auth = token
	}

	// these seem to be the most sensible mappings
	identifier := factory.NewDriverIdentifier(
		factory.Mapping("dev.azure.com", "azure"),
		factory.Mapping("bitbucket.org", "bitbucketcloud"),
		factory.Mapping("fake.com", "fake"),
	)

	driver, err := identifier.Identify(u.Host)
	if err != nil {
		return nil, "", err
	}

	u.Path = "/"
	u.User = nil

	client, err := factory.NewClient(driver, u.String(), auth)
	return client, auth, err
}

func DetermineToken(credentials string, repositoryUrl string) (string, error) {
	lines := strings.Split(credentials, "\n")
	for _, line := range lines {
		u, err := url.Parse(strings.TrimSpace(line))
		if err != nil {
			return "", err
		}

		h, err := url.Parse(repositoryUrl)
		if err != nil {
			return "", err
		}

		if h.Host == u.Host {
			// we have found a host that matches
			password, ok := u.User.Password()

			if ok {
				return password, nil
			}
		}
	}

	return "", fmt.Errorf("unable to locate a token for %s", repositoryUrl)
}
