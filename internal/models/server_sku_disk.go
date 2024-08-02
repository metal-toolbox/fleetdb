// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ServerSkuDisk is an object representing the database table.
type ServerSkuDisk struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	SkuID     string    `boil:"sku_id" json:"sku_id" toml:"sku_id" yaml:"sku_id"`
	Bytes     int64     `boil:"bytes" json:"bytes" toml:"bytes" yaml:"bytes"`
	Protocol  string    `boil:"protocol" json:"protocol" toml:"protocol" yaml:"protocol"`
	Count     int64     `boil:"count" json:"count" toml:"count" yaml:"count"`
	CreatedAt null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *serverSkuDiskR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L serverSkuDiskL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ServerSkuDiskColumns = struct {
	ID        string
	SkuID     string
	Bytes     string
	Protocol  string
	Count     string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	SkuID:     "sku_id",
	Bytes:     "bytes",
	Protocol:  "protocol",
	Count:     "count",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ServerSkuDiskTableColumns = struct {
	ID        string
	SkuID     string
	Bytes     string
	Protocol  string
	Count     string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "server_sku_disk.id",
	SkuID:     "server_sku_disk.sku_id",
	Bytes:     "server_sku_disk.bytes",
	Protocol:  "server_sku_disk.protocol",
	Count:     "server_sku_disk.count",
	CreatedAt: "server_sku_disk.created_at",
	UpdatedAt: "server_sku_disk.updated_at",
}

// Generated where

var ServerSkuDiskWhere = struct {
	ID        whereHelperstring
	SkuID     whereHelperstring
	Bytes     whereHelperint64
	Protocol  whereHelperstring
	Count     whereHelperint64
	CreatedAt whereHelpernull_Time
	UpdatedAt whereHelpernull_Time
}{
	ID:        whereHelperstring{field: "\"server_sku_disk\".\"id\""},
	SkuID:     whereHelperstring{field: "\"server_sku_disk\".\"sku_id\""},
	Bytes:     whereHelperint64{field: "\"server_sku_disk\".\"bytes\""},
	Protocol:  whereHelperstring{field: "\"server_sku_disk\".\"protocol\""},
	Count:     whereHelperint64{field: "\"server_sku_disk\".\"count\""},
	CreatedAt: whereHelpernull_Time{field: "\"server_sku_disk\".\"created_at\""},
	UpdatedAt: whereHelpernull_Time{field: "\"server_sku_disk\".\"updated_at\""},
}

// ServerSkuDiskRels is where relationship names are stored.
var ServerSkuDiskRels = struct {
	Sku string
}{
	Sku: "Sku",
}

// serverSkuDiskR is where relationships are stored.
type serverSkuDiskR struct {
	Sku *ServerSku `boil:"Sku" json:"Sku" toml:"Sku" yaml:"Sku"`
}

// NewStruct creates a new relationship struct
func (*serverSkuDiskR) NewStruct() *serverSkuDiskR {
	return &serverSkuDiskR{}
}

func (r *serverSkuDiskR) GetSku() *ServerSku {
	if r == nil {
		return nil
	}
	return r.Sku
}

// serverSkuDiskL is where Load methods for each relationship are stored.
type serverSkuDiskL struct{}

var (
	serverSkuDiskAllColumns            = []string{"id", "sku_id", "bytes", "protocol", "count", "created_at", "updated_at"}
	serverSkuDiskColumnsWithoutDefault = []string{"sku_id", "bytes", "protocol", "count"}
	serverSkuDiskColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	serverSkuDiskPrimaryKeyColumns     = []string{"id"}
	serverSkuDiskGeneratedColumns      = []string{}
)

type (
	// ServerSkuDiskSlice is an alias for a slice of pointers to ServerSkuDisk.
	// This should almost always be used instead of []ServerSkuDisk.
	ServerSkuDiskSlice []*ServerSkuDisk
	// ServerSkuDiskHook is the signature for custom ServerSkuDisk hook methods
	ServerSkuDiskHook func(context.Context, boil.ContextExecutor, *ServerSkuDisk) error

	serverSkuDiskQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	serverSkuDiskType                 = reflect.TypeOf(&ServerSkuDisk{})
	serverSkuDiskMapping              = queries.MakeStructMapping(serverSkuDiskType)
	serverSkuDiskPrimaryKeyMapping, _ = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, serverSkuDiskPrimaryKeyColumns)
	serverSkuDiskInsertCacheMut       sync.RWMutex
	serverSkuDiskInsertCache          = make(map[string]insertCache)
	serverSkuDiskUpdateCacheMut       sync.RWMutex
	serverSkuDiskUpdateCache          = make(map[string]updateCache)
	serverSkuDiskUpsertCacheMut       sync.RWMutex
	serverSkuDiskUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var serverSkuDiskAfterSelectHooks []ServerSkuDiskHook

var serverSkuDiskBeforeInsertHooks []ServerSkuDiskHook
var serverSkuDiskAfterInsertHooks []ServerSkuDiskHook

var serverSkuDiskBeforeUpdateHooks []ServerSkuDiskHook
var serverSkuDiskAfterUpdateHooks []ServerSkuDiskHook

var serverSkuDiskBeforeDeleteHooks []ServerSkuDiskHook
var serverSkuDiskAfterDeleteHooks []ServerSkuDiskHook

var serverSkuDiskBeforeUpsertHooks []ServerSkuDiskHook
var serverSkuDiskAfterUpsertHooks []ServerSkuDiskHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ServerSkuDisk) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ServerSkuDisk) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ServerSkuDisk) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ServerSkuDisk) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ServerSkuDisk) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ServerSkuDisk) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ServerSkuDisk) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ServerSkuDisk) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ServerSkuDisk) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSkuDiskAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddServerSkuDiskHook registers your hook function for all future operations.
func AddServerSkuDiskHook(hookPoint boil.HookPoint, serverSkuDiskHook ServerSkuDiskHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		serverSkuDiskAfterSelectHooks = append(serverSkuDiskAfterSelectHooks, serverSkuDiskHook)
	case boil.BeforeInsertHook:
		serverSkuDiskBeforeInsertHooks = append(serverSkuDiskBeforeInsertHooks, serverSkuDiskHook)
	case boil.AfterInsertHook:
		serverSkuDiskAfterInsertHooks = append(serverSkuDiskAfterInsertHooks, serverSkuDiskHook)
	case boil.BeforeUpdateHook:
		serverSkuDiskBeforeUpdateHooks = append(serverSkuDiskBeforeUpdateHooks, serverSkuDiskHook)
	case boil.AfterUpdateHook:
		serverSkuDiskAfterUpdateHooks = append(serverSkuDiskAfterUpdateHooks, serverSkuDiskHook)
	case boil.BeforeDeleteHook:
		serverSkuDiskBeforeDeleteHooks = append(serverSkuDiskBeforeDeleteHooks, serverSkuDiskHook)
	case boil.AfterDeleteHook:
		serverSkuDiskAfterDeleteHooks = append(serverSkuDiskAfterDeleteHooks, serverSkuDiskHook)
	case boil.BeforeUpsertHook:
		serverSkuDiskBeforeUpsertHooks = append(serverSkuDiskBeforeUpsertHooks, serverSkuDiskHook)
	case boil.AfterUpsertHook:
		serverSkuDiskAfterUpsertHooks = append(serverSkuDiskAfterUpsertHooks, serverSkuDiskHook)
	}
}

// One returns a single serverSkuDisk record from the query.
func (q serverSkuDiskQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ServerSkuDisk, error) {
	o := &ServerSkuDisk{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for server_sku_disk")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ServerSkuDisk records from the query.
func (q serverSkuDiskQuery) All(ctx context.Context, exec boil.ContextExecutor) (ServerSkuDiskSlice, error) {
	var o []*ServerSkuDisk

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ServerSkuDisk slice")
	}

	if len(serverSkuDiskAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ServerSkuDisk records in the query.
func (q serverSkuDiskQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count server_sku_disk rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q serverSkuDiskQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if server_sku_disk exists")
	}

	return count > 0, nil
}

// Sku pointed to by the foreign key.
func (o *ServerSkuDisk) Sku(mods ...qm.QueryMod) serverSkuQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SkuID),
	}

	queryMods = append(queryMods, mods...)

	return ServerSkus(queryMods...)
}

// LoadSku allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (serverSkuDiskL) LoadSku(ctx context.Context, e boil.ContextExecutor, singular bool, maybeServerSkuDisk interface{}, mods queries.Applicator) error {
	var slice []*ServerSkuDisk
	var object *ServerSkuDisk

	if singular {
		var ok bool
		object, ok = maybeServerSkuDisk.(*ServerSkuDisk)
		if !ok {
			object = new(ServerSkuDisk)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeServerSkuDisk)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeServerSkuDisk))
			}
		}
	} else {
		s, ok := maybeServerSkuDisk.(*[]*ServerSkuDisk)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeServerSkuDisk)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeServerSkuDisk))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &serverSkuDiskR{}
		}
		args = append(args, object.SkuID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &serverSkuDiskR{}
			}

			for _, a := range args {
				if a == obj.SkuID {
					continue Outer
				}
			}

			args = append(args, obj.SkuID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`server_sku`),
		qm.WhereIn(`server_sku.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ServerSku")
	}

	var resultSlice []*ServerSku
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ServerSku")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for server_sku")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for server_sku")
	}

	if len(serverSkuAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Sku = foreign
		if foreign.R == nil {
			foreign.R = &serverSkuR{}
		}
		foreign.R.SkuServerSkuDisks = append(foreign.R.SkuServerSkuDisks, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SkuID == foreign.ID {
				local.R.Sku = foreign
				if foreign.R == nil {
					foreign.R = &serverSkuR{}
				}
				foreign.R.SkuServerSkuDisks = append(foreign.R.SkuServerSkuDisks, local)
				break
			}
		}
	}

	return nil
}

// SetSku of the serverSkuDisk to the related item.
// Sets o.R.Sku to related.
// Adds o to related.R.SkuServerSkuDisks.
func (o *ServerSkuDisk) SetSku(ctx context.Context, exec boil.ContextExecutor, insert bool, related *ServerSku) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"server_sku_disk\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sku_id"}),
		strmangle.WhereClause("\"", "\"", 2, serverSkuDiskPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SkuID = related.ID
	if o.R == nil {
		o.R = &serverSkuDiskR{
			Sku: related,
		}
	} else {
		o.R.Sku = related
	}

	if related.R == nil {
		related.R = &serverSkuR{
			SkuServerSkuDisks: ServerSkuDiskSlice{o},
		}
	} else {
		related.R.SkuServerSkuDisks = append(related.R.SkuServerSkuDisks, o)
	}

	return nil
}

// ServerSkuDisks retrieves all the records using an executor.
func ServerSkuDisks(mods ...qm.QueryMod) serverSkuDiskQuery {
	mods = append(mods, qm.From("\"server_sku_disk\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"server_sku_disk\".*"})
	}

	return serverSkuDiskQuery{q}
}

// FindServerSkuDisk retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindServerSkuDisk(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ServerSkuDisk, error) {
	serverSkuDiskObj := &ServerSkuDisk{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"server_sku_disk\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, serverSkuDiskObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from server_sku_disk")
	}

	if err = serverSkuDiskObj.doAfterSelectHooks(ctx, exec); err != nil {
		return serverSkuDiskObj, err
	}

	return serverSkuDiskObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ServerSkuDisk) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_sku_disk provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverSkuDiskColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	serverSkuDiskInsertCacheMut.RLock()
	cache, cached := serverSkuDiskInsertCache[key]
	serverSkuDiskInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			serverSkuDiskAllColumns,
			serverSkuDiskColumnsWithDefault,
			serverSkuDiskColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"server_sku_disk\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"server_sku_disk\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into server_sku_disk")
	}

	if !cached {
		serverSkuDiskInsertCacheMut.Lock()
		serverSkuDiskInsertCache[key] = cache
		serverSkuDiskInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ServerSkuDisk.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ServerSkuDisk) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	serverSkuDiskUpdateCacheMut.RLock()
	cache, cached := serverSkuDiskUpdateCache[key]
	serverSkuDiskUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			serverSkuDiskAllColumns,
			serverSkuDiskPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update server_sku_disk, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"server_sku_disk\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, serverSkuDiskPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, append(wl, serverSkuDiskPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update server_sku_disk row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for server_sku_disk")
	}

	if !cached {
		serverSkuDiskUpdateCacheMut.Lock()
		serverSkuDiskUpdateCache[key] = cache
		serverSkuDiskUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q serverSkuDiskQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for server_sku_disk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for server_sku_disk")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ServerSkuDiskSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSkuDiskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"server_sku_disk\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, serverSkuDiskPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in serverSkuDisk slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all serverSkuDisk")
	}
	return rowsAff, nil
}

// Delete deletes a single ServerSkuDisk record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ServerSkuDisk) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ServerSkuDisk provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), serverSkuDiskPrimaryKeyMapping)
	sql := "DELETE FROM \"server_sku_disk\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from server_sku_disk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for server_sku_disk")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q serverSkuDiskQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no serverSkuDiskQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from server_sku_disk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_sku_disk")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ServerSkuDiskSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(serverSkuDiskBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSkuDiskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"server_sku_disk\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverSkuDiskPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from serverSkuDisk slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_sku_disk")
	}

	if len(serverSkuDiskAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ServerSkuDisk) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindServerSkuDisk(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ServerSkuDiskSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ServerSkuDiskSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSkuDiskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"server_sku_disk\".* FROM \"server_sku_disk\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverSkuDiskPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ServerSkuDiskSlice")
	}

	*o = slice

	return nil
}

// ServerSkuDiskExists checks if the ServerSkuDisk row exists.
func ServerSkuDiskExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"server_sku_disk\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if server_sku_disk exists")
	}

	return exists, nil
}

// Exists checks if the ServerSkuDisk row exists.
func (o *ServerSkuDisk) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ServerSkuDiskExists(ctx, exec, o.ID)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ServerSkuDisk) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_sku_disk provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverSkuDiskColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	serverSkuDiskUpsertCacheMut.RLock()
	cache, cached := serverSkuDiskUpsertCache[key]
	serverSkuDiskUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			serverSkuDiskAllColumns,
			serverSkuDiskColumnsWithDefault,
			serverSkuDiskColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			serverSkuDiskAllColumns,
			serverSkuDiskPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert server_sku_disk, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(serverSkuDiskPrimaryKeyColumns))
			copy(conflict, serverSkuDiskPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryCockroachDB(dialect, "\"server_sku_disk\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(serverSkuDiskType, serverSkuDiskMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		_, _ = fmt.Fprintln(boil.DebugWriter, cache.query)
		_, _ = fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // CockcorachDB doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert server_sku_disk")
	}

	if !cached {
		serverSkuDiskUpsertCacheMut.Lock()
		serverSkuDiskUpsertCache[key] = cache
		serverSkuDiskUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
