package gql

import (
	"github.com/gelleson/changescout/changescout/internal/api/gql/model"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/gelleson/changescout/changescout/pkg/diff"
	"net/http"
)

func buildSetting(input *model.SettingInput) domain.Setting {
	var setting model.SettingInput
	if input != nil {
		setting = *input
	}

	return domain.Setting{
		Headers:       http.Header{},
		UserAgent:     transform.ToValueOrDefault(setting.UserAgent, ""),
		Referer:       transform.ToValueOrDefault(setting.Referer, ""),
		Template:      diff.GetUpdatedValueWithPointer(input.Template, input.Template),
		Method:        setting.Method.String(),
		Selectors:     setting.Selectors,
		Deduplication: transform.ToValueOrDefault(setting.Deduplication, false),
		Trim:          transform.ToValueOrDefault(setting.Trim, false),
		Sort:          transform.ToValueOrDefault(setting.Sort, false),
		JSONPath:      setting.JSONPath,
	}
}
