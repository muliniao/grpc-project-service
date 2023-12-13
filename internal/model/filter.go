package model

type Filter struct {
	Offset int64
	Limit  int64
}

func DefaultFilter() Filter {
	return Filter{
		Offset: 0,
		Limit:  20,
	}
}
