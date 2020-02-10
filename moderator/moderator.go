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

type ActionHandlerFunc func(id int64, resolution string)

type Moderator struct {
	ctx            context.Context
	forum          *forum.Forum
	actionCache    map[int64]struct{}
	actionHandlers map[string]*ActionHandlerFunc
	moderators     map[string]struct{}
	titleFormat    string
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
	err := (*m.forum).CreateAction(action)
	if err != nil {
		return err
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

func (m *Moderator) PollActions() {
	go m.actionHandlerService()
}

// actionHandlerService runs the handleNewlyResolvedActions function every N duration
// It caches the timestamp of the last runtime in order to improve speed.
func (m *Moderator) actionHandlerService() {
	ticker := time.NewTicker(1 * time.Hour)
	timestamp := time.Unix(0, 0)
	for ; true; <-ticker.C {
		m.handleNewlyResolvedActions(timestamp)
		timestamp = time.Now()
	}
}

// handleNewlyResolvedActions synchronously finds all newly resolved actions and
// calls the associated handler function.
func (m *Moderator) handleNewlyResolvedActions(timestamp time.Time) {
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
		for _, comment := range comments {
			if _, ok := m.moderators[comment.User]; ok {
				if resolution, ok := (*m.forum).GetResolution(comment); ok {
					actionHandler := *m.actionHandlers[resolution]
					if actionHandler != nil {
						if err := (*m.forum).CloseAction(action.ID); err == nil {
							if id, err := parseItemIdFromAction(action.Title, m.titleFormat); err == nil {
								actionHandler(id, resolution)
							}
						}
					}
				}
			}
		}
	}
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
