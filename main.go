package main

import (
	"fmt"
	"github.com/henrymxu/gomoderator/forum"
	"github.com/henrymxu/gomoderator/moderator"
	"os"
)

func main() {
	//githubBuilder, _ := forum.InitializeGithubBuilderFromConfig("src/github.com/henrymxu/gomoderator/config.toml")
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
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(mod.DoesActionAlreadyExist(0))
	//err = mod.CreateAction(1, "Test Issue for GoSports using GoModerator")
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	mod.StartActionsPollingService()
}

func actionHandler(id int64, resolution string) {
	fmt.Printf("Handling action for %d with resolution %s\n", id, resolution)
}