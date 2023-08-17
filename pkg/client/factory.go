package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

// FromRepoURL parses a URL of the form https://:authtoken@host/ and attempts to
// determine the driver and creates a client to authenticate to the endpoint.
func FromRepoURL(repoURL string, credentials string, kind string) (*scm.Client, string, string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, "", "", err
	}

	var username string
	var auth string

	if password, ok := u.User.Password(); ok {
		auth = password
		username = u.User.Username()
	} else {
		fmt.Println("[DEBUG] Token is not available from the url, falling back to .git-credentials")
		user, token, err := DetermineAuth(credentials, repoURL)
		if err != nil {
			return nil, "", "", err
		}
		auth = token
		username = user
	}

	// these seem to be the most sensible mappings
	identifier := factory.NewDriverIdentifier(
		factory.Mapping("dev.azure.com", "azure"),
		factory.Mapping("bitbucket.org", "bitbucketcloud"),
		factory.Mapping("fake.com", "fake"),
	)

	var driver string

	if kind == "" {
		driver, err = identifier.Identify(u.Host)
		if err != nil {
			return nil, "", "", err
		}
	} else {
		driver = kind
	}

	u.Path = "/"
	u.User = nil

	client, err := factory.NewClient(driver, u.String(), auth)
	client.Username = username

	return client, username, auth, err
}

func DetermineAuth(credentials string, repositoryURL string) (string, string, error) {
	lines := strings.Split(credentials, "\n")
	for _, line := range lines {
		u, err := url.Parse(strings.TrimSpace(line))
		if err != nil {
			return "", "", err
		}

		h, err := url.Parse(repositoryURL)
		if err != nil {
			return "", "", err
		}

		if h.Host == u.Host {
			// we have found a host that matches
			password, ok := u.User.Password()

			if ok {
				return u.User.Username(), password, nil
			}
		}
	}

	return "", "", fmt.Errorf("unable to locate a token for %s", repositoryURL)
}
