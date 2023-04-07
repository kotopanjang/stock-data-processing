package utils

type Pagination struct {
	Data       any
	Count      int
	TotalCount int
	Page       int
	TotalPage  int
	Error      error
}
