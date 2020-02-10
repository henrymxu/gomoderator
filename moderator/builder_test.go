package moderator

import (
	"github.com/henrymxu/gomoderator/forum"
	"testing"
)

func createGithubBuilder() *forum.GithubBuilder {
	githubBuilder:= forum.NewGithubBuilder()
	githubBuilder.AccessToken = "test"
	githubBuilder.RepositoryName = "test"
	githubBuilder.RepositoryOwner = "test"
	githubBuilder.AccountName = "test"
	return githubBuilder
}

func TestBuilder_SetTitleFormat(t *testing.T) {
	builder := NewModeratorBuilder()
	err := builder.SetTitleFormat("Moderator action for %d")
	if err != nil {
		t.Error(err)
	}
}

func TestBuilder_SetTitleFormatFail(t *testing.T) {
	builder := NewModeratorBuilder()
	err := builder.SetTitleFormat("Moderator action for %s")
	if err == nil {
		t.Error(err)
	}
}

func TestBuilder_BuildModerator(t *testing.T) {
	builder := NewModeratorBuilder()
	builder.SetForumBuilder(createGithubBuilder())
	_ = builder.SetTitleFormat("Valid Title Format %d")
	builder.SetModerators("testModerator")
	_, err := builder.BuildModerator()
	if err != nil {
		t.Error(err)
	}
}

func TestBuilder_BuildModeratorFail(t *testing.T) {
	builder := NewModeratorBuilder()
	_, err := builder.BuildModerator()
	if err == nil {
		t.Error(err)
	}
}
