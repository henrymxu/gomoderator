package forum

import (
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"time"
)

const MaxPageSize = 100
const IssueSort = "updated"

type Github struct {
	client *github.Client
	ctx    context.Context
	configuration
}

type configuration struct {
	accountName string
	owner       string
	repo        string
	labels      *[]string
}

func (g *Github) GetVotingResolution(comments []*Comment) (string, bool) {
	votes := make(map[string]int)
	for _, comment := range comments {
		if comment.User == g.accountName {
			reactionsTotal := 0
			for _, val := range comment.Reactions {
				reactionsTotal += val
			}
			votes[comment.Body] = reactionsTotal
		}
	}
	mostVotes := 0
	optionsWithMostVotes := make([]string, 0)
	for option, value := range votes {
		if value > mostVotes {
			optionsWithMostVotes = []string{option}
			mostVotes = value
		} else if value == mostVotes {
			optionsWithMostVotes = append(optionsWithMostVotes, option)
		}
	}
	if optionsWithMostVotes == nil || len(optionsWithMostVotes) != 1 {
		return "", false
	}
	return optionsWithMostVotes[0], true
}

func (g *Github) GetCommentatingResolution(comments []*Comment, resolutions map[string]struct{}, moderators map[string]struct{}) (string, bool) {
	for _, comment := range comments {
		if _, ok := moderators[comment.User]; ok {
			if _, ok := resolutions[comment.Body]; ok {
				return comment.Body, true
			}
		}
	}
	return "", false
}

func (g *Github) CloseAction(number int) error {
	closed := "closed"
	issueRequest := &github.IssueRequest{
		State: &closed,
	}
	_, _, err := g.client.Issues.Edit(g.ctx, g.owner, g.repo, number, issueRequest)
	return err
}

func (g *Github) CreateAction(action *Action) (int, error) {
	issueRequest := &github.IssueRequest{
		Title:  &action.Title,
		Labels: g.labels,
		Body:   &action.Body,
	}
	issue, err := g.postAction(issueRequest)
	if err != nil {
		return 0, nil
	}
	return issue.GetNumber(), nil
}

func (g *Github) GetActions(state string, timestamp time.Time) ([]*Action, error) {
	actions := make([]*Action, 0)
	pageNumber := 0
	for len(actions)%100 == 0 && len(actions) > 0 || pageNumber == 0 {
		newIssues, err := g.getActions(state, pageNumber, timestamp)
		if err != nil {
			return nil, err
		}
		actions = append(actions, convertIssuesToActions(newIssues)...)
		pageNumber++
	}
	return actions, nil
}

func (g *Github) GetNewComments(action *Action, timestamp time.Time) ([]*Comment, error) {
	options := &github.IssueListCommentsOptions{
		Since: &timestamp,
	}
	comments, _, err := g.client.Issues.ListComments(g.ctx, g.owner, g.repo, action.ID, options)
	return convertIssueCommentsToComment(comments), err
}

func (g *Github) PostComment(body string, id int) error {
	_, err := g.postComment(body, id)
	return err
}

func (g *Github) isAction(issue *github.Issue) bool {
	if issue.IsPullRequest() {
		return false
	}
	userName := issue.GetUser().GetLogin()
	return userName == g.accountName
}

func (g *Github) postAction(request *github.IssueRequest) (*github.Issue, error) {
	issue, _, err := g.client.Issues.Create(g.ctx, g.owner, g.repo, request)
	return issue, err
}

func (g *Github) postComment(contents string, issueNumber int) (*github.IssueComment, error) {
	issueComment := github.IssueComment{
		Body: &contents,
	}
	res, _, err := g.client.Issues.CreateComment(g.ctx, g.owner, g.repo, issueNumber, &issueComment)
	return res, err
}

func (g *Github) getActions(filterState string, pageNumber int, timestamp time.Time) ([]*github.Issue, error) {
	options := &github.IssueListByRepoOptions{
		Creator: g.accountName,
		Sort:    IssueSort,
		State:   filterState,
		ListOptions: github.ListOptions{
			PerPage: MaxPageSize,
			Page:    pageNumber,
		},
		Since: timestamp,
	}
	issues, _, err := g.client.Issues.ListByRepo(g.ctx, g.owner, g.repo, options)
	return issues, err
}

func convertIssueToAction(issue *github.Issue) *Action {
	labels := make([]string, len(issue.Labels))
	for i, label := range issue.Labels {
		labels[i] = *label.Name
	}
	action := Action{
		ID:       issue.GetNumber(),
		Title:    issue.GetTitle(),
		Body:     issue.GetBody(),
		User:     issue.GetUser().GetLogin(),
		Labels:   labels,
		Comments: issue.GetComments(),
	}
	return &action
}

func convertIssuesToActions(issues []*github.Issue) []*Action {
	actions := make([]*Action, len(issues))
	for i, issue := range issues {
		actions[i] = convertIssueToAction(issue)
	}
	return actions
}

func convertIssueCommentsToComment(issueComments []*github.IssueComment) []*Comment {
	comments := make([]*Comment, len(issueComments))
	for i, issueComment := range issueComments {
		comments[i] = &Comment{
			Body:      issueComment.GetBody(),
			User:      issueComment.GetUser().GetLogin(),
			Reactions: convertReactionsToMap(issueComment.Reactions),
		}
	}
	return comments
}

func convertReactionsToMap(reactions *github.Reactions) map[string]int {
	result := make(map[string]int)
	result["Confused"] = reactions.GetConfused()
	result["Heart"] = reactions.GetHeart()
	result["Hooray"] = reactions.GetHooray()
	result["Laugh"] = reactions.GetLaugh()
	result["MinusOne"] = reactions.GetMinusOne()
	result["PlusOne"] = reactions.GetPlusOne()
	return result
}
