package driver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GeneralSettings struct {
	MaxLiveTime   time.Duration `json:"max-live-time"`
	UseGenManager bool          `json:"flag-generator-manager"`
}

func (g *GeneralSettings) String() string {
	return fmt.Sprint(g.MaxLiveTime, g.UseGenManager)
}

func (g *GeneralSettings) BytesJSON() ([]byte, error) {
	return json.Marshal(g)
}

func parseGeneral(v any) (GeneralSettings, error) {
	switch x := v.(type) {
	case string:
		return parseGeneralString(x)
	case []byte:
		return parseGeneralJSON(x)
	case json.RawMessage:
		return parseGeneralJSON(x)
	default:
		return parseGeneralStruct(x)
	}
}

func parseGeneralJSON(input []byte) (GeneralSettings, error) {
	genSet := &GeneralSettings{}
	if err := json.Unmarshal(input, genSet); err != nil {
		return GeneralSettings{}, err
	}
	return *genSet, nil
}

func parseGeneralString(input string) (GeneralSettings, error) {
	parts := strings.Split(input, delimiter)
	if len(parts) != 2 {
		return GeneralSettings{}, errInvalidSettings
	}
	genSet := GeneralSettings{}
	genSet.MaxLiveTime, _ = time.ParseDuration(parts[0])
	if genSet.MaxLiveTime > MaxLiveTime {
		return GeneralSettings{}, errLiveTimeLong
	}
	genSet.UseGenManager, _ = strconv.ParseBool(parts[1])
	return genSet, nil
}

func parseGeneralStruct(v any) (GeneralSettings, error) {
	generalSettings, ok := v.(*GeneralSettings)
	if !ok || generalSettings == nil {
		return GeneralSettings{}, errInvalidSettings
	}
	if generalSettings.MaxLiveTime > MaxLiveTime {
		return GeneralSettings{}, errLiveTimeLong
	}
	return *generalSettings, nil
}
