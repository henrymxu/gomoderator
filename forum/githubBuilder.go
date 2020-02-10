package forum

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

type GithubBuilder struct {
	AccessToken string
	AccountName string
	RepositoryOwner string
	RepositoryName string
	IssueLabels *[]string
}

func NewGithubBuilder() *GithubBuilder {
	return &GithubBuilder{}
}

func InitializeGithubBuilderFromConfig(configPath string) (*GithubBuilder, error) {
	var githubBuilder GithubBuilder
	_, err := toml.DecodeFile(configPath, &githubBuilder)
	if err != nil {
		return nil, err
	}
	return &githubBuilder, nil
}

func (g *GithubBuilder) Build() (*Forum, error) {
	if g.AccessToken == "" {
		return nil, errors.New("github forum must have access token")
	}
	if g.AccountName == "" {
		return nil, errors.New("github forum must contain account name")
	}
	if g.RepositoryName == "" || g.RepositoryOwner == "" {
		return nil, errors.New("github forum must contain repository owner and repository name")
	}
	ctx := context.Background()
	client := github.NewClient(createAuthClient(ctx, g.AccessToken))
	var githubForum Forum
	githubForum = &Github{
		ctx: ctx,
		client: client,
		configuration: configuration{
			accountName:g.AccountName,
			owner: g.RepositoryOwner,
			repo: g.RepositoryName,
			labels: g.IssueLabels,
		},
	}
	return &githubForum, nil
}

