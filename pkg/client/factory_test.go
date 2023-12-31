package client_test

import (
	"context"
	"testing"

	"github.com/garethjevans/scm/pkg/client"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/assert"
)

func TestFromRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://fake.com/myorg/myrepo.git", "https://user:token@fake.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	assert.Equal(t, "//fake.com/", c.BaseURL.String())

	orgs, _, err := c.Organizations.List(context.Background(), &scm.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, orgs, 5)
}

func TestFromGithubRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://github.com/garethjevans/scm.git", "https://user:token@github.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	assert.Equal(t, "https://api.github.com/", c.BaseURL.String())
}

func TestFromGithubEnterpriseRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://my.ghe.com/garethjevans/scm.git", "https://user:token@my.ghe.com", "github")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	assert.Equal(t, "https://my.ghe.com/api/v3/", c.BaseURL.String())
}

func TestFromGitlabRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://gitlab.com/garethjevans/scm.git", "https://user:token@gitlab.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	assert.Equal(t, "https://gitlab.com/", c.BaseURL.String())
}

func TestFromGitlabInternalRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://gl.eng.xxx/garethjevans/scm.git", "https://user:token@gl.eng.xxx", "gitlab")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	t.Log(c)

	assert.Equal(t, "https://gl.eng.xxx/", c.BaseURL.String())
}

func TestFromAzureRepoURL(t *testing.T) {
	c, username, token, err := client.FromRepoURL("https://dev.azure.com/garethjevans/_/scm.git", "https://user:token@dev.azure.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "token", token)
	assert.Equal(t, "user", username)

	assert.Equal(t, "https://dev.azure.com/", c.BaseURL.String())
}
