package forum

import "testing"

func TestInitializeGithubBuilderFromConfig(t *testing.T) {
	githubBuilder, _ := InitializeGithubBuilderFromConfig("src/github.com/henrymxu/gomoderator/config.toml")
	_, err := githubBuilder.Build()
	if err != nil {
		t.Error()
	}
}

func TestGithubBuilder_Build(t *testing.T) {
	githubBuilder := NewGithubBuilder()
	githubBuilder.AccessToken = "access_token"
	githubBuilder.AccountName = "account_name"
	githubBuilder.RepositoryName = "repo_name"
	githubBuilder.RepositoryOwner = "repo_owner"
	_, err := githubBuilder.Build()
	if err != nil {
		t.Fail()
	}
}

func TestGithubBuilder_BuildFail(t *testing.T) {
	githubBuilder := NewGithubBuilder()
	_, err := githubBuilder.Build()
	if err == nil {
		t.Error(err)
	}
	githubBuilder.AccessToken = "access_token"
	_, err = githubBuilder.Build()
	if err == nil {
		t.Error(err)
	}
	githubBuilder.AccountName = "account_name"
	_, err = githubBuilder.Build()
	if err == nil {
		t.Error(err)
	}
	githubBuilder.RepositoryName = "repo_name"
	githubBuilder.RepositoryOwner = "repo_owner"
	_, err = githubBuilder.Build()
	if err != nil {
		t.Error(err)
	}
}
