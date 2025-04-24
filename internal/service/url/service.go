package url

type Service interface {
	Shorten(original string) (code string, err error)
	Resolve(code string) (original string, err error)
}
