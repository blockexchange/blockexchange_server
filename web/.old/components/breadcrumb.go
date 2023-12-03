package components

type BreadcrumbModel struct {
	Entries []BreadcrumbEntry
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
