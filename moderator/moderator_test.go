package moderator

import (
	"fmt"
	"testing"
)

func TestModerator_ParseItemIdFromAction(t *testing.T) {
	format := "Here is an item ID %d, try to parse it"
	id := int64(27)
	title := fmt.Sprintf(format, id)
	parsedId, err := parseItemIdFromAction(title, format)
	if err != nil || parsedId != id {
		t.Error(err)
	}
}

func TestModerator_ParseItemIdFromAction2(t *testing.T) {
	format := "%d"
	id := int64(27)
	title := fmt.Sprintf(format, id)
	parsedId, err := parseItemIdFromAction(title, format)
	if err != nil || parsedId != id {
		t.Error(err)
	}
}




