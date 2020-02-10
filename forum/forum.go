package forum

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type Forum interface {
	GetResolution(comment *Comment) (string, bool)
	CloseAction(id int64) error
	CreateAction(action *Action) error
	GetActions(state string, timestamp time.Time) ([]*Action, error)
	GetNewComments(action *Action, timestamp time.Time) ([]*Comment, error)
}

type Builder interface {
	Build() (*Forum, error)
}

type Action struct {
	ID        int64
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
