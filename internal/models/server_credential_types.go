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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ServerCredentialType is an object representing the database table.
type ServerCredentialType struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug      string    `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Builtin   bool      `boil:"builtin" json:"builtin" toml:"builtin" yaml:"builtin"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *serverCredentialTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L serverCredentialTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ServerCredentialTypeColumns = struct {
	ID        string
	Name      string
	Slug      string
	Builtin   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Name:      "name",
	Slug:      "slug",
	Builtin:   "builtin",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ServerCredentialTypeTableColumns = struct {
	ID        string
	Name      string
	Slug      string
	Builtin   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "server_credential_types.id",
	Name:      "server_credential_types.name",
	Slug:      "server_credential_types.slug",
	Builtin:   "server_credential_types.builtin",
	CreatedAt: "server_credential_types.created_at",
	UpdatedAt: "server_credential_types.updated_at",
}

// Generated where

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ServerCredentialTypeWhere = struct {
	ID        whereHelperstring
	Name      whereHelperstring
	Slug      whereHelperstring
	Builtin   whereHelperbool
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperstring{field: "\"server_credential_types\".\"id\""},
	Name:      whereHelperstring{field: "\"server_credential_types\".\"name\""},
	Slug:      whereHelperstring{field: "\"server_credential_types\".\"slug\""},
	Builtin:   whereHelperbool{field: "\"server_credential_types\".\"builtin\""},
	CreatedAt: whereHelpertime_Time{field: "\"server_credential_types\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"server_credential_types\".\"updated_at\""},
}

// ServerCredentialTypeRels is where relationship names are stored.
var ServerCredentialTypeRels = struct {
	ServerCredentials string
}{
	ServerCredentials: "ServerCredentials",
}

// serverCredentialTypeR is where relationships are stored.
type serverCredentialTypeR struct {
	ServerCredentials ServerCredentialSlice `boil:"ServerCredentials" json:"ServerCredentials" toml:"ServerCredentials" yaml:"ServerCredentials"`
}

// NewStruct creates a new relationship struct
func (*serverCredentialTypeR) NewStruct() *serverCredentialTypeR {
	return &serverCredentialTypeR{}
}

func (r *serverCredentialTypeR) GetServerCredentials() ServerCredentialSlice {
	if r == nil {
		return nil
	}
	return r.ServerCredentials
}

// serverCredentialTypeL is where Load methods for each relationship are stored.
type serverCredentialTypeL struct{}

var (
	serverCredentialTypeAllColumns            = []string{"id", "name", "slug", "builtin", "created_at", "updated_at"}
	serverCredentialTypeColumnsWithoutDefault = []string{"name", "slug", "created_at", "updated_at"}
	serverCredentialTypeColumnsWithDefault    = []string{"id", "builtin"}
	serverCredentialTypePrimaryKeyColumns     = []string{"id"}
	serverCredentialTypeGeneratedColumns      = []string{}
)

type (
	// ServerCredentialTypeSlice is an alias for a slice of pointers to ServerCredentialType.
	// This should almost always be used instead of []ServerCredentialType.
	ServerCredentialTypeSlice []*ServerCredentialType
	// ServerCredentialTypeHook is the signature for custom ServerCredentialType hook methods
	ServerCredentialTypeHook func(context.Context, boil.ContextExecutor, *ServerCredentialType) error

	serverCredentialTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	serverCredentialTypeType                 = reflect.TypeOf(&ServerCredentialType{})
	serverCredentialTypeMapping              = queries.MakeStructMapping(serverCredentialTypeType)
	serverCredentialTypePrimaryKeyMapping, _ = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, serverCredentialTypePrimaryKeyColumns)
	serverCredentialTypeInsertCacheMut       sync.RWMutex
	serverCredentialTypeInsertCache          = make(map[string]insertCache)
	serverCredentialTypeUpdateCacheMut       sync.RWMutex
	serverCredentialTypeUpdateCache          = make(map[string]updateCache)
	serverCredentialTypeUpsertCacheMut       sync.RWMutex
	serverCredentialTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var serverCredentialTypeAfterSelectHooks []ServerCredentialTypeHook

var serverCredentialTypeBeforeInsertHooks []ServerCredentialTypeHook
var serverCredentialTypeAfterInsertHooks []ServerCredentialTypeHook

var serverCredentialTypeBeforeUpdateHooks []ServerCredentialTypeHook
var serverCredentialTypeAfterUpdateHooks []ServerCredentialTypeHook

var serverCredentialTypeBeforeDeleteHooks []ServerCredentialTypeHook
var serverCredentialTypeAfterDeleteHooks []ServerCredentialTypeHook

var serverCredentialTypeBeforeUpsertHooks []ServerCredentialTypeHook
var serverCredentialTypeAfterUpsertHooks []ServerCredentialTypeHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ServerCredentialType) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ServerCredentialType) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ServerCredentialType) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ServerCredentialType) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ServerCredentialType) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ServerCredentialType) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ServerCredentialType) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ServerCredentialType) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ServerCredentialType) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverCredentialTypeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddServerCredentialTypeHook registers your hook function for all future operations.
func AddServerCredentialTypeHook(hookPoint boil.HookPoint, serverCredentialTypeHook ServerCredentialTypeHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		serverCredentialTypeAfterSelectHooks = append(serverCredentialTypeAfterSelectHooks, serverCredentialTypeHook)
	case boil.BeforeInsertHook:
		serverCredentialTypeBeforeInsertHooks = append(serverCredentialTypeBeforeInsertHooks, serverCredentialTypeHook)
	case boil.AfterInsertHook:
		serverCredentialTypeAfterInsertHooks = append(serverCredentialTypeAfterInsertHooks, serverCredentialTypeHook)
	case boil.BeforeUpdateHook:
		serverCredentialTypeBeforeUpdateHooks = append(serverCredentialTypeBeforeUpdateHooks, serverCredentialTypeHook)
	case boil.AfterUpdateHook:
		serverCredentialTypeAfterUpdateHooks = append(serverCredentialTypeAfterUpdateHooks, serverCredentialTypeHook)
	case boil.BeforeDeleteHook:
		serverCredentialTypeBeforeDeleteHooks = append(serverCredentialTypeBeforeDeleteHooks, serverCredentialTypeHook)
	case boil.AfterDeleteHook:
		serverCredentialTypeAfterDeleteHooks = append(serverCredentialTypeAfterDeleteHooks, serverCredentialTypeHook)
	case boil.BeforeUpsertHook:
		serverCredentialTypeBeforeUpsertHooks = append(serverCredentialTypeBeforeUpsertHooks, serverCredentialTypeHook)
	case boil.AfterUpsertHook:
		serverCredentialTypeAfterUpsertHooks = append(serverCredentialTypeAfterUpsertHooks, serverCredentialTypeHook)
	}
}

// One returns a single serverCredentialType record from the query.
func (q serverCredentialTypeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ServerCredentialType, error) {
	o := &ServerCredentialType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for server_credential_types")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ServerCredentialType records from the query.
func (q serverCredentialTypeQuery) All(ctx context.Context, exec boil.ContextExecutor) (ServerCredentialTypeSlice, error) {
	var o []*ServerCredentialType

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ServerCredentialType slice")
	}

	if len(serverCredentialTypeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ServerCredentialType records in the query.
func (q serverCredentialTypeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count server_credential_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q serverCredentialTypeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if server_credential_types exists")
	}

	return count > 0, nil
}

// ServerCredentials retrieves all the server_credential's ServerCredentials with an executor.
func (o *ServerCredentialType) ServerCredentials(mods ...qm.QueryMod) serverCredentialQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"server_credentials\".\"server_credential_type_id\"=?", o.ID),
	)

	return ServerCredentials(queryMods...)
}

// LoadServerCredentials allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (serverCredentialTypeL) LoadServerCredentials(ctx context.Context, e boil.ContextExecutor, singular bool, maybeServerCredentialType interface{}, mods queries.Applicator) error {
	var slice []*ServerCredentialType
	var object *ServerCredentialType

	if singular {
		var ok bool
		object, ok = maybeServerCredentialType.(*ServerCredentialType)
		if !ok {
			object = new(ServerCredentialType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeServerCredentialType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeServerCredentialType))
			}
		}
	} else {
		s, ok := maybeServerCredentialType.(*[]*ServerCredentialType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeServerCredentialType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeServerCredentialType))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &serverCredentialTypeR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &serverCredentialTypeR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`server_credentials`),
		qm.WhereIn(`server_credentials.server_credential_type_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load server_credentials")
	}

	var resultSlice []*ServerCredential
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice server_credentials")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on server_credentials")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for server_credentials")
	}

	if len(serverCredentialAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.ServerCredentials = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &serverCredentialR{}
			}
			foreign.R.ServerCredentialType = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ServerCredentialTypeID {
				local.R.ServerCredentials = append(local.R.ServerCredentials, foreign)
				if foreign.R == nil {
					foreign.R = &serverCredentialR{}
				}
				foreign.R.ServerCredentialType = local
				break
			}
		}
	}

	return nil
}

// AddServerCredentials adds the given related objects to the existing relationships
// of the server_credential_type, optionally inserting them as new records.
// Appends related to o.R.ServerCredentials.
// Sets related.R.ServerCredentialType appropriately.
func (o *ServerCredentialType) AddServerCredentials(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*ServerCredential) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ServerCredentialTypeID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"server_credentials\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"server_credential_type_id"}),
				strmangle.WhereClause("\"", "\"", 2, serverCredentialPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ServerCredentialTypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &serverCredentialTypeR{
			ServerCredentials: related,
		}
	} else {
		o.R.ServerCredentials = append(o.R.ServerCredentials, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &serverCredentialR{
				ServerCredentialType: o,
			}
		} else {
			rel.R.ServerCredentialType = o
		}
	}
	return nil
}

// ServerCredentialTypes retrieves all the records using an executor.
func ServerCredentialTypes(mods ...qm.QueryMod) serverCredentialTypeQuery {
	mods = append(mods, qm.From("\"server_credential_types\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"server_credential_types\".*"})
	}

	return serverCredentialTypeQuery{q}
}

// FindServerCredentialType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindServerCredentialType(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ServerCredentialType, error) {
	serverCredentialTypeObj := &ServerCredentialType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"server_credential_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, serverCredentialTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from server_credential_types")
	}

	if err = serverCredentialTypeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return serverCredentialTypeObj, err
	}

	return serverCredentialTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ServerCredentialType) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_credential_types provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverCredentialTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	serverCredentialTypeInsertCacheMut.RLock()
	cache, cached := serverCredentialTypeInsertCache[key]
	serverCredentialTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			serverCredentialTypeAllColumns,
			serverCredentialTypeColumnsWithDefault,
			serverCredentialTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"server_credential_types\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"server_credential_types\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into server_credential_types")
	}

	if !cached {
		serverCredentialTypeInsertCacheMut.Lock()
		serverCredentialTypeInsertCache[key] = cache
		serverCredentialTypeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ServerCredentialType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ServerCredentialType) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	serverCredentialTypeUpdateCacheMut.RLock()
	cache, cached := serverCredentialTypeUpdateCache[key]
	serverCredentialTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			serverCredentialTypeAllColumns,
			serverCredentialTypePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update server_credential_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"server_credential_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, serverCredentialTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, append(wl, serverCredentialTypePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update server_credential_types row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for server_credential_types")
	}

	if !cached {
		serverCredentialTypeUpdateCacheMut.Lock()
		serverCredentialTypeUpdateCache[key] = cache
		serverCredentialTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q serverCredentialTypeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for server_credential_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for server_credential_types")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ServerCredentialTypeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverCredentialTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"server_credential_types\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, serverCredentialTypePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in serverCredentialType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all serverCredentialType")
	}
	return rowsAff, nil
}

// Delete deletes a single ServerCredentialType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ServerCredentialType) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ServerCredentialType provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), serverCredentialTypePrimaryKeyMapping)
	sql := "DELETE FROM \"server_credential_types\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from server_credential_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for server_credential_types")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q serverCredentialTypeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no serverCredentialTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from server_credential_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_credential_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ServerCredentialTypeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(serverCredentialTypeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverCredentialTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"server_credential_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverCredentialTypePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from serverCredentialType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_credential_types")
	}

	if len(serverCredentialTypeAfterDeleteHooks) != 0 {
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
func (o *ServerCredentialType) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindServerCredentialType(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ServerCredentialTypeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ServerCredentialTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverCredentialTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"server_credential_types\".* FROM \"server_credential_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverCredentialTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ServerCredentialTypeSlice")
	}

	*o = slice

	return nil
}

// ServerCredentialTypeExists checks if the ServerCredentialType row exists.
func ServerCredentialTypeExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"server_credential_types\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if server_credential_types exists")
	}

	return exists, nil
}

// Exists checks if the ServerCredentialType row exists.
func (o *ServerCredentialType) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ServerCredentialTypeExists(ctx, exec, o.ID)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ServerCredentialType) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_credential_types provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverCredentialTypeColumnsWithDefault, o)

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

	serverCredentialTypeUpsertCacheMut.RLock()
	cache, cached := serverCredentialTypeUpsertCache[key]
	serverCredentialTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			serverCredentialTypeAllColumns,
			serverCredentialTypeColumnsWithDefault,
			serverCredentialTypeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			serverCredentialTypeAllColumns,
			serverCredentialTypePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert server_credential_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(serverCredentialTypePrimaryKeyColumns))
			copy(conflict, serverCredentialTypePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryCockroachDB(dialect, "\"server_credential_types\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(serverCredentialTypeType, serverCredentialTypeMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert server_credential_types")
	}

	if !cached {
		serverCredentialTypeUpsertCacheMut.Lock()
		serverCredentialTypeUpsertCache[key] = cache
		serverCredentialTypeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
