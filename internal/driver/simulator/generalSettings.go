package driver

import (
	"encoding/json"
	"fmt"
	m "practice/internal/models"
	"strconv"
	"strings"
	"time"
)

type GeneralSettings struct {
	ProgramLiveTime time.Duration `json:"program-live-time"`
	GenOptimization bool          `json:"use-generator-optimization"`
}

func (g *GeneralSettings) String() string {
	return fmt.Sprint(g.ProgramLiveTime, m.DELIMITER, g.GenOptimization)
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
	if genSet.ProgramLiveTime > m.MAX_LIVE_TIME {
		return GeneralSettings{}, m.ErrLiveTimeLong
	}
	return *genSet, nil
}

func parseGeneralString(input string) (GeneralSettings, error) {
	parts := strings.Split(input, m.DELIMITER)
	if len(parts) != 2 {
		return GeneralSettings{}, m.ErrInvalidSettings
	}
	genSet := GeneralSettings{}
	genSet.ProgramLiveTime, _ = time.ParseDuration(parts[0])
	if genSet.ProgramLiveTime > m.MAX_LIVE_TIME {
		return GeneralSettings{}, m.ErrLiveTimeLong
	}
	genSet.GenOptimization, _ = strconv.ParseBool(parts[1])
	return genSet, nil
}

func parseGeneralStruct(v any) (GeneralSettings, error) {
	generalSettings, ok := v.(*GeneralSettings)
	if !ok || generalSettings == nil {
		return GeneralSettings{}, m.ErrInvalidSettings
	}
	if generalSettings.ProgramLiveTime > m.MAX_LIVE_TIME {
		return GeneralSettings{}, m.ErrLiveTimeLong
	}
	return *generalSettings, nil
}
