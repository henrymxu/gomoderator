package main

import (
	"fmt"
	"github.com/henrymxu/gomoderator/forum"
	"github.com/henrymxu/gomoderator/moderator"
	"os"
)

func ExampleModerator_DoesActionAlreadyExist() {
	githubBuilder := forum.NewGithubBuilder()
	githubBuilder.AccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	githubBuilder.AccountName = "henrymxu"
	githubBuilder.RepositoryOwner = "henrymxu"
	githubBuilder.RepositoryName = "gosports"

	builder := moderator.NewModeratorBuilder()
	builder.SetForumBuilder(githubBuilder)
	builder.SetModerators("henrymxu")
	_ = builder.SetResolutions("pass", "fail")
	builder.RegisterActionHandler(actionHandler)
	builder.SetModeToCommenting()
	_ = builder.SetTitleFormat("Action required for %d")
	mod, err := builder.BuildModerator()
	if err != nil {
		panic(err)
	}
	fmt.Println(mod.DoesActionAlreadyExist(0))
	// Output: true
}

func actionHandler(id int64, resolution string) {
	fmt.Printf("Handling action for %d with resolution %s\n", id, resolution)
}