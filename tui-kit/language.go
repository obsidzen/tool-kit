package tuikit

import "strings"

const (
	KeyLanguage = "l"
	KeyDetail   = "?"
	LangEnglish = "en"
	LangKorean  = "ko"
)

type Language struct {
	Tag  string
	Name string
}

type MessageCatalog map[string]map[string]string

type Translator struct {
	Language string
	Catalog  MessageCatalog
}

func DefaultLanguages() []Language {
	return []Language{
		{Tag: LangEnglish, Name: "English"},
		{Tag: LangKorean, Name: "한국어"},
	}
}

func NewTranslator(language string, catalog MessageCatalog) Translator {
	language = strings.TrimSpace(language)
	if language == "" {
		language = LangEnglish
	}
	return Translator{Language: language, Catalog: catalog}
}

func (t Translator) T(key, fallback string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return fallback
	}
	if value := t.lookup(t.Language, key); value != "" {
		return value
	}
	if value := t.lookup(LangEnglish, key); value != "" {
		return value
	}
	return fallback
}

func (t Translator) lookup(language, key string) string {
	if t.Catalog == nil {
		return ""
	}
	messages := t.Catalog[language]
	if messages == nil {
		return ""
	}
	return messages[key]
}

func (t Translator) LanguageLabel(languages []Language) string {
	for _, language := range languages {
		if language.Tag == t.Language {
			return language.Name
		}
	}
	return t.Language
}

func LanguageSelectItems(languages []Language) []SelectItem {
	items := make([]SelectItem, len(languages))
	for i, language := range languages {
		items[i] = SelectItem{Key: language.Tag, Label: language.Name}
	}
	return items
}

func SelectLanguageByTag(model SelectModel, tag string) SelectModel {
	for i, item := range model.Items {
		if item.Key == tag {
			model.Index = i
			return model
		}
	}
	return model
}

func CommonMessages() MessageCatalog {
	return MessageCatalog{
		LangKorean: {
			"common.help.select_run_language_quit":        "↑/↓ 선택 · enter 실행 · l 언어 · q 종료",
			"common.help.select_run_detail_language_quit": "↑/↓ 선택 · enter 실행 · ? 설명 · l 언어 · q 종료",
			"common.help.enter_select_back_quit":          "enter 선택 · esc 뒤로 · q 종료",
			"common.help.enter_open_back_quit":            "enter 열기 · esc 뒤로 · q 종료",
			"common.help.any_home_quit":                   "아무 키 홈 · q 종료",
			"common.help.detail":                          "esc 뒤로 · q 종료",
			"common.help.running_quit":                    "실행 중 · q 종료",
			"common.help.language":                        "↑/↓ 언어 선택 · enter 적용 · esc 뒤로 · q 종료",
			"common.language.title":                       "언어 선택",
			"common.detail.title":                         "설명",
			"common.status.completed":                     "완료",
			"common.status.running":                       "실행 중",
			"common.status.language":                      "언어",
			"common.empty.running":                        "(실행 중)",
		},
	}
}

func MergeCatalogs(catalogs ...MessageCatalog) MessageCatalog {
	merged := MessageCatalog{}
	for _, catalog := range catalogs {
		for language, messages := range catalog {
			if merged[language] == nil {
				merged[language] = map[string]string{}
			}
			for key, value := range messages {
				merged[language][key] = value
			}
		}
	}
	return merged
}
