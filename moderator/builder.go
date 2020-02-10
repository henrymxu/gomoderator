package moderator

import (
	"errors"
	"github.com/henrymxu/gomoderator/forum"
	"golang.org/x/net/context"
	"strings"
)

type Builder struct {
	forumBuilder *forum.Builder
	titleFormat  string
	moderators   []string
}

func NewModeratorBuilder() *Builder {
	return &Builder{}
}

func (m *Builder) SetForumBuilder(forumBuilder forum.Builder) {
	m.forumBuilder = &forumBuilder
}

func (m *Builder) SetTitleFormat(format string) error {
	if !strings.Contains(format, "%d") {
		return errors.New("invalid title format, must contain %d for item ID placeholder")
	}
	m.titleFormat = format
	return nil
}

func (m *Builder) SetModerators(mods ...string) {
	m.moderators = mods
}

func (m *Builder) BuildModerator() (*Moderator, error) {
	if m.titleFormat == "" {
		return nil, errors.New("moderator must contain a title format")
	}
	if m.moderators == nil {
		return nil, errors.New("moderator must authorize at least one user account to be a moderator")
	}
	if m.forumBuilder == nil {
		return nil, errors.New("moderator must contain a forum builder")
	}
	return buildModerator(m)
}

func buildModerator(builder *Builder) (*Moderator, error) {
	builtForum, err := (*builder.forumBuilder).Build()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	moderatorsMap := make(map[string]struct{})
	for _, moderator := range builder.moderators {
		moderatorsMap[moderator] = struct{}{}
	}
	moderator := Moderator{
		ctx:         ctx,
		forum:       builtForum,
		moderators:  moderatorsMap,
		titleFormat: builder.titleFormat,
	}
	return &moderator, nil
}
