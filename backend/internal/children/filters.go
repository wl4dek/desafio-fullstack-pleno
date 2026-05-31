package children

type Filters struct {
	Name         string
	Neighborhood string
	Alert        string
	Reviewed     *bool
	HasAlert     *bool
	Page         int
	PerPage      int
}

func (f *Filters) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PerPage < 10 {
		f.PerPage = 10
	}
	if f.PerPage > 50 {
		f.PerPage = 50
	}
}

func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PerPage
}
