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

func TestBuilder_SetResolutions(t *testing.T) {
	builder := NewModeratorBuilder()
	err := builder.SetResolutions("pass", "fail", "nothing")
	if err != nil {
		t.Error(err)
	}
}

func TestBuilder_SetResolutionsFail(t *testing.T) {
	builder := NewModeratorBuilder()
	err := builder.SetResolutions( "fail")
	if err == nil {
		t.Error(err)
	}
}

func TestBuilder_BuildModerator(t *testing.T) {
	builder := NewModeratorBuilder()
	builder.SetForumBuilder(createGithubBuilder())
	_ = builder.SetTitleFormat("Valid Title Format %d")
	builder.SetModerators("testModerator")
	_ = builder.SetResolutions("pass", "fail")
	builder.SetModeToCommenting()
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
	_ = builder.SetTitleFormat("Test Title Format %d")
	_, err = builder.BuildModerator()
	if err == nil {
		t.Error(err)
	}
	builder.SetModerators("testModerator")
	_, err = builder.BuildModerator()
	if err == nil {
		t.Error(err)
	}
	builder.SetForumBuilder(createGithubBuilder())
	_, err = builder.BuildModerator()
	if err == nil {
		t.Error(err)
	}
	_ = builder.SetResolutions("pass", "fail")
	_, err = builder.BuildModerator()
	if err == nil {
		t.Error(err)
	}
	builder.SetModeToCommenting()
	_, err = builder.BuildModerator()
	if err != nil {
		t.Error(err)
	}
}
