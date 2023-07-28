package driver

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TagSettings struct {
	PollTime time.Duration `json:"poll-time"`
	Settings string        `json:"settings"`
}

func (t *TagSettings) String() string {
	return fmt.Sprint(t.PollTime, t.Settings)
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
	return *tagSet, nil
}

func parseTagsString(input string) (TagSettings, error) {
	idx := strings.Index(input, delimiter)
	if idx < 0 {
		return TagSettings{}, errInvalidSettings
	}

	tagSet := TagSettings{}
	tagSet.PollTime, _ = time.ParseDuration(input[:idx])
	if tagSet.PollTime < maxPrescaler {
		return TagSettings{}, errPrescallerSmall
	}
	tagSet.Settings = input[idx+1:]
	return tagSet, nil
}

func parseTagsStruct(v any) (TagSettings, error) {
	tagSet, ok := v.(*TagSettings)
	if !ok || tagSet == nil {
		return TagSettings{}, errInvalidSettings
	}
	return *tagSet, nil
}
