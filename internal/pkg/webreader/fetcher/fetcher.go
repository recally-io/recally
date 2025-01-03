package fetcher

type FecherType string

const (
	TypeHttp       FecherType = "http"
	TypeJinaReader FecherType = "jinaReader"
	TypeBrowser    FecherType = "browser"
	TypeNil        FecherType = ""
)
