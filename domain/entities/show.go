package entities

type Show struct {
	ID       int
	Name     string
	Overview string
	Seasons  []Season
}
