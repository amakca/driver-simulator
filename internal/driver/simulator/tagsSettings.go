package driver

import (
	"encoding/json"
	"fmt"
	m "practice/internal/models"
	"strings"
	"time"
)

type TagSettings struct {
	PollTime  time.Duration `json:"poll-time"`
	GenConfig string        `json:"generator-config"`
}

func (t *TagSettings) String() string {
	return fmt.Sprint(t.PollTime, m.DELIMITER, t.GenConfig)
}

func (t *TagSettings) BytesJSON() ([]byte, error) {
	return json.Marshal(t)
}

func parseTags(v any) (TagSettings, error) {
	switch x := v.(type) {
	case string:
		return parseTagsString(x)
	case []byte:
		return parseTagsJSON(x)
	case json.RawMessage:
		return parseTagsJSON(x)
	default:
		return parseTagsStruct(x)
	}
}

func parseTagsJSON(input []byte) (TagSettings, error) {
	tagSet := &TagSettings{}
	if err := json.Unmarshal(input, tagSet); err != nil {
		return TagSettings{}, err
	}
	if tagSet.PollTime < m.MIN_POLL_TIME {
		return TagSettings{}, m.ErrPollTimeSmall
	}
	return *tagSet, nil
}

func parseTagsString(input string) (TagSettings, error) {
	var err error

	idx := strings.Index(input, m.DELIMITER)
	if idx < 0 {
		return TagSettings{}, m.ErrInvalidSettings
	}

	tagSet := TagSettings{}
	tagSet.PollTime, err = time.ParseDuration(input[:idx])
	if err != nil {
		return TagSettings{}, err
	}
	if tagSet.PollTime < m.MIN_POLL_TIME {
		return TagSettings{}, m.ErrPollTimeSmall
	}

	tagSet.GenConfig = input[idx+1:]
	return tagSet, nil
}

func parseTagsStruct(v any) (TagSettings, error) {
	tagSet, ok := v.(*TagSettings)
	if !ok || tagSet == nil {
		return TagSettings{}, m.ErrInvalidSettings
	}
	if tagSet.PollTime < m.MIN_POLL_TIME {
		return TagSettings{}, m.ErrPollTimeSmall
	}
	return *tagSet, nil
}
