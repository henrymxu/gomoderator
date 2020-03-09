package moderator

import (
	"errors"
	"fmt"
	"github.com/henrymxu/gomoderator/forum"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

const (
	votingMode     = "voting"
	commentingMode = "commenting"
)

type ActionHandlerFunc func(id int64, resolution string)

type Moderator struct {
	ctx           context.Context
	forum         *forum.Forum
	moderators    map[string]struct{}
	resolutions   map[string]struct{}
	actionCache   map[int64]struct{}
	actionHandler *ActionHandlerFunc
	titleFormat   string
	mode          string
	pollingFrequency time.Duration
}

// CreateAction attempts to create a moderator action based on a unique ID.
// It will check if the item is a duplicate issue, and if so return an error.
// If the network request fails, it will return an error as well.
func (m *Moderator) CreateAction(id int64, body string) error {
	if m.DoesActionAlreadyExist(id) {
		return errors.New(fmt.Sprintf("Issue %d already exists", id))
	}
	title := fmt.Sprintf(m.titleFormat, id)
	action := &forum.Action{
		Title: title,
		Body:  body,
	}
	actionId, err := (*m.forum).CreateAction(action)
	if err != nil {
		return err
	}
	if m.mode == votingMode {
		for resolution := range m.resolutions {
			_ = (*m.forum).PostComment(resolution, actionId)
		}
	}
	m.actionCache[id] = struct{}{}
	return nil
}

// DoesActionAlreadyExist checks if the action has been created before.
// It will reference a cache of IDs that will be created upon startup.
func (m *Moderator) DoesActionAlreadyExist(id int64) bool {
	if m.actionCache != nil {
		_, ok := m.actionCache[id]
		return ok
	}
	actions, err := (*m.forum).GetActions("all", time.Unix(0, 0))
	if err != nil {
		return false
	}
	result := false
	cache := make(map[int64]struct{})
	for _, action := range actions {
		if itemId, err := parseItemIdFromAction(action.Title, m.titleFormat); err == nil {
			cache[itemId] = struct{}{}
			if id == itemId {
				result = true
			}
		}
	}
	m.actionCache = cache
	return result
}

func (m *Moderator) StartActionsPollingService() {
	go m.actionHandlerService()
}

func (m *Moderator) RegisterActionHandler(handlerFunc ActionHandlerFunc) {
	m.actionHandler = &handlerFunc
}

// actionHandlerService runs the findAndHandleNewlyResolvedActions function every N duration
// It caches the timestamp of the last runtime in order to improve speed.
func (m *Moderator) actionHandlerService() {
	ticker := time.NewTicker(m.pollingFrequency)
	timestamp := time.Unix(0, 0)
	for ; true; <-ticker.C {
		m.findAndHandleNewlyResolvedActions(timestamp)
		timestamp = time.Now()
	}
}

// findAndHandleNewlyResolvedActions synchronously finds all newly resolved actions and
// calls the associated handler function.
func (m *Moderator) findAndHandleNewlyResolvedActions(timestamp time.Time) {
	actions, err := (*m.forum).GetActions("open", timestamp)
	if err != nil {
		return
	}
	for _, action := range actions {
		if action.Comments == 0 {
			continue
		}
		comments, err := (*m.forum).GetNewComments(action, timestamp)
		if err != nil {
			continue
		}
		if resolution, ok := m.getResolution(comments); ok {
			if m.actionHandler != nil {
				if err := (*m.forum).CloseAction(action.ID); err == nil {
					if id, err := parseItemIdFromAction(action.Title, m.titleFormat); err == nil {
						(*m.actionHandler)(id, resolution)
					}
				}
			}
		}
	}
}

func (m *Moderator) getResolution(comments []*forum.Comment) (string, bool) {
	if m.mode == votingMode {
		return (*m.forum).GetVotingResolution(comments)
	} else if m.mode == commentingMode {
		return (*m.forum).GetCommentatingResolution(comments, m.resolutions, m.moderators)
	}
	return "", false
}

// parseItemIdFromAction retrieves the item ID from the title of the post.
func parseItemIdFromAction(title string, format string) (int64, error) {
	index := strings.Index(format, "%d")
	prefix := format[0:index]
	suffix := format[index+2:]
	idString := strings.TrimSuffix(strings.TrimPrefix(title, prefix), suffix)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
