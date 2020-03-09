package forum

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type Forum interface {
	GetVotingResolution(comments []*Comment) (string, bool)
	GetCommentatingResolution(comments []*Comment, resolutions map[string]struct{}, moderators map[string]struct{}) (string, bool)
	CloseAction(number int) error
	CreateAction(action *Action) (int, error)
	GetActions(state string, prevTime time.Time) ([]*Action, error)
	GetNewComments(action *Action, prevTime time.Time) ([]*Comment, error)
	PostComment(body string, id int) error
}

type Builder interface {
	Build() (*Forum, error)
}

type Action struct {
	ID        int
	Title     string
	Body      string
	User      string
	Labels    []string
	Comments  int
	Reactions map[string]int
}

type Comment struct {
	Body      string
	User      string
	Reactions map[string]int
}

func createAuthClient(ctx context.Context, accessToken string) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}
