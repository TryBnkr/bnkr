package paginator

import "math"

type Paginator struct {
	CurrentPage int
	PerPage     int
	TotalCount  int
}

func (p *Paginator) Offset() int {
	// Assuming 20 items per page:
	// page 1 has an offset of 0    (1-1) * 20
	// page 2 has an offset of 20   (2-1) * 20
	//   in other words, page 2 starts with item 21
	return (p.CurrentPage - 1) * p.PerPage
}

func (p *Paginator) TotalPages() int {
	return int(math.Ceil(float64(p.TotalCount) / float64(p.PerPage)))
}

func (p *Paginator) PreviousPage() int {
	return p.CurrentPage - 1
}

func (p *Paginator) NextPage() int {
	return p.CurrentPage + 1
}

func (p *Paginator) HasPreviousPage() bool {
	return p.PreviousPage() >= 1
}

func (p *Paginator) HasNextPage() bool {
	return p.NextPage() <= p.TotalPages()
}
