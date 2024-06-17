package fleetdbapi

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	// MaxPaginationSize represents the maximum number of records that can be returned per page
	MaxPaginationSize = 1000
	// DefaultPaginationSize represents the default number of records that are returned per page
	DefaultPaginationSize = 100
)

// PaginationParams allow you to paginate the results.
// Some tables can have multiple preloadable child tables. Set Preload to true to load them.
// TODO; Preload should probably be moved over to the params of each individual endpoint. Example: ServerListParams.
type PaginationParams struct {
	Limit   int    `json:"limit,omitempty"`
	Page    int    `json:"page,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	Preload bool   `json:"preload,omitempty"`
	OrderBy string `json:"orderby,omitempty"`
}

type paginationData struct {
	pageCount  int
	totalCount int64
	pager      PaginationParams
}

// TODO; Replace with a query parser like done here: bio_config_set_params.go:parseBiosConfigSetListParams()
func parsePagination(c *gin.Context) (PaginationParams, error) {
	// Initializing default
	var err error
	limit := DefaultPaginationSize
	page := 1
	query := c.Request.URL.Query()
	preload := false
	orderby := ""

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, err = strconv.Atoi(queryValue)
		case "page":
			page, err = strconv.Atoi(queryValue)
		case "preload":
			preload, err = strconv.ParseBool(queryValue)
		case "orderby":
			orderby = queryValue
		}

		if err != nil {
			return PaginationParams{}, err
		}
	}

	return PaginationParams{
		Limit:   limit,
		Page:    page,
		Preload: preload,
		OrderBy: orderby,
	}, nil
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *PaginationParams) queryMods() []qm.QueryMod {
	if p == nil {
		p = &PaginationParams{}
	}

	mods := []qm.QueryMod{}

	mods = append(mods, qm.Limit(p.limitUsed()))

	if p.Page != 0 {
		mods = append(mods, qm.Offset(p.offset()))
	}

	// TODO; match the old functionality for now, will handle order and load as params later
	if p.OrderBy != "" {
		mods = append(mods, qm.OrderBy(p.OrderBy))
	}

	return mods
}

// serverQueryMods queryMods converts the list params into sql conditions that can be added to sql queries
func (p *PaginationParams) serverQueryMods() []qm.QueryMod {
	mods := p.queryMods()

	if p.Preload {
		preload := []qm.QueryMod{
			qm.Load("Attributes"),
			qm.Load("VersionedAttributes", qm.Where("(server_id, namespace, created_at) IN (select server_id, namespace, max(created_at) from versioned_attributes group by server_id, namespace)")),
			qm.Load("ServerComponents.Attributes"),
			qm.Load("ServerComponents.ServerComponentType"),
		}
		mods = append(mods, preload...)
	}

	return mods
}

// serverComponentQueryMods converts the server component list params into sql conditions that can be added to sql queries
func (p *PaginationParams) serverComponentsQueryMods() []qm.QueryMod {
	mods := p.queryMods()

	preload := []qm.QueryMod{
		qm.Load("Attributes"),
		qm.Load("VersionedAttributes", qm.Where("(server_component_id, namespace, created_at) IN (select server_component_id, namespace, max(created_at) from versioned_attributes group by server_component_id, namespace)")),
		qm.Load("ServerComponentType"),
	}
	mods = append(mods, preload...)

	return mods
}

func (p *PaginationParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}

	if p.Page != 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}

	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}

	if p.Preload { // No need to send preload=false. It will default to that on the other side if its empty
		q.Set("preload", "1") // strconv.ParseBool() will convert 1 to TRUE. No reason to send extra bytes when you dont have to.
	}

	if p.OrderBy != "" {
		q.Set("orderby", p.OrderBy)
	}
}

func (p *PaginationParams) limitUsed() int {
	limit := p.Limit

	switch {
	case limit > MaxPaginationSize:
		limit = MaxPaginationSize
	case limit <= 0:
		limit = DefaultPaginationSize
	}

	return limit
}

func (p *PaginationParams) offset() int {
	page := p.Page
	if page == 0 {
		page = 1
	}

	return (page - 1) * p.limitUsed()
}
