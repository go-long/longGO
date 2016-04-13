package translation

import (
	"jex/cn/longGo/fb/i18n/language"
)

type singleTranslation struct {
	id       string
	template *template
}

func (st *singleTranslation) MarshalInterface() interface{} {
	return map[string]interface{}{
		"id":          st.id,
		"translation": st.template,
	}
}

func (st *singleTranslation) ID() string {
	return st.id
}

func (st *singleTranslation) Template(pc language.Plural) *template {
	return st.template
}

func (st *singleTranslation) UntranslatedCopy() Translation {
	return &singleTranslation{st.id, mustNewTemplate("")}
}

func (st *singleTranslation) Normalize(language *language.Language) Translation {
	return st
}

func (st *singleTranslation) Backfill(src Translation) Translation {
	if st.template == nil || st.template.src == "" {
		st.template = src.Template(language.Other)
	}
	return st
}

func (st *singleTranslation) Merge(t Translation) Translation {
	other, ok := t.(*singleTranslation)
	if !ok || st.ID() != t.ID() {
		return t
	}
	if other.template != nil && other.template.src != "" {
		st.template = other.template
	}
	return st
}

func (st *singleTranslation) Incomplete(l *language.Language) bool {
	return st.template == nil || st.template.src == ""
}

func (st *singleTranslation)Templates()map[string]string{
	ts :=make(map[string]string)
	 ts[""]=  st.template.src

	return ts
}

var _ = Translation(&singleTranslation{})
