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
	ctx context.Context
	configuration
}

type configuration struct {
	accountName string
	owner       string
	repo        string
	labels      *[]string
}

func (g *Github) GetResolution(comment *Comment) (string, bool) {

	return "", false
}

func (g *Github) CloseAction(id int64) error {
	closed := "closed"
	issueRequest := &github.IssueRequest{
		State: &closed,
	}
	_, _, err := g.client.Issues.Edit(g.ctx, g.owner, g.repo, int(id), issueRequest)
	return err
}

func (g *Github) CreateAction(action *Action) error {
	issueRequest := &github.IssueRequest{
		Title: &action.Title,
		Labels: g.labels,
		Body: &action.Body,
	}
	_, err := g.postAction(issueRequest)
	return err
}

func (g *Github) GetActions(state string, timestamp time.Time) ([]*Action, error) {
	actions := make([]*Action, 0)
	pageNumber := 0
	for len(actions) % 100 == 0 && len(actions) > 0 || pageNumber == 0 {
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
	options := &github.IssueListCommentsOptions {
		Since: &timestamp,
	}
	comments, _, err := g.client.Issues.ListComments(g.ctx, g.owner, g.repo, int(action.ID), options)
	return convertIssueCommentsToComment(comments), err
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

func (g *Github) getActions(filterState string, pageNumber int, timestamp time.Time) ([]*github.Issue, error){
	options := &github.IssueListByRepoOptions{
		Creator: g.accountName,
		Sort:    IssueSort,
		State:   filterState,
		ListOptions: github.ListOptions{
			PerPage: MaxPageSize,
			Page: pageNumber,
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
		ID:issue.GetID(),
		Title:issue.GetTitle(),
		Body:issue.GetBody(),
		User:issue.GetUser().GetLogin(),
		Labels:labels,
		Comments:issue.GetComments(),
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
			Body:issueComment.GetBody(),
			User:issueComment.GetUser().GetLogin(),
		}
	}
	return comments
}

