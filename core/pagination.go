package core

type PaginationParams struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}

func (p *PaginationParams) Default() {
	var (
		defaultLimit  = 100
		defaultOffset = 0
	)

	if p.Limit == nil {
		p.Limit = &defaultLimit
	}

	if p.Offset == nil {
		p.Offset = &defaultOffset
	}

	if *p.Limit < 0 {
		p.Limit = &defaultLimit
	}

	if *p.Offset < 0 {
		p.Offset = &defaultOffset
	}
}

type Metadata struct {
	Self     PaginationParams  `json:"self"`
	Next     *PaginationParams `json:"next,omitempty"`
	Previous *PaginationParams `json:"previous,omitempty"`
}

type PagedResponse[T any] struct {
	Metadata Metadata `json:"metadata"`
	Total    int      `json:"total"`
	Items    []T      `json:"items"`
}

func Page[T any](items []T, total, limit, offset int) PagedResponse[T] {
	p := PagedResponse[T]{}

	p.Total = total
	p.Items = items

	p.Metadata.Self.Limit = &limit
	p.Metadata.Self.Offset = &offset

	nextPage := (offset + limit)
	if nextPage < total {
		p.Metadata.Next = &PaginationParams{
			Limit:  &limit,
			Offset: &nextPage,
		}
	}

	if offset > 0 {
		previousPage := (offset - limit)
		p.Metadata.Previous = &PaginationParams{
			Limit:  &limit,
			Offset: &previousPage,
		}
	}

	return p
}
