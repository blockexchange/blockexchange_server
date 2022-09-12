package components

type BreadcrumbModel struct {
	Entries []BreadcrumbEntry
	BaseURL string
}

type BreadcrumbEntry struct {
	Link   string
	Name   string
	Active bool
}

func Breadcrumb(entries ...BreadcrumbEntry) *BreadcrumbModel {
	m := &BreadcrumbModel{
		Entries: entries,
	}

	return m
}
