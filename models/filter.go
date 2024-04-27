package models

import "strings"

type (
	KeySort  string
	SortType string

	Filter struct {
		Keyword     string
		WarehouseID int
		KeySort     KeySort
		SortType    SortType
	}
)

const (
	SortTypeDesc SortType = "desc"
	SortTypeAsc  SortType = "asc"

	KeySortStatus KeySort = "status"
	KeySortName   KeySort = "name"
)

func (m *Filter) GetKeyword() string {
	m.Keyword = strings.ReplaceAll(m.Keyword, `"`, "")
	m.Keyword = strings.ReplaceAll(m.Keyword, `;`, "")
	m.Keyword = strings.ReplaceAll(m.Keyword, `'`, "")
	m.Keyword = strings.ReplaceAll(m.Keyword, `1=1`, "")
	return m.Keyword
}

func (t SortType) Value() bool {
	switch t {
	case SortTypeDesc:
		return true
	case SortTypeAsc:
		return false
	}
	return true
}
