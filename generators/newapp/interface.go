package newapp

type Gen interface {
	RootPath(curpath string) string
	ModelPath(curpath string) string
}
