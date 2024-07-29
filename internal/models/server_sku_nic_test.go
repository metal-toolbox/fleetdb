// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

func testServerSkuNicsUpsert(t *testing.T) {
	t.Parallel()

	if len(serverSkuNicAllColumns) == len(serverSkuNicPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ServerSkuNic{}
	if err = randomize.Struct(seed, &o, serverSkuNicDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSkuNic: %s", err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, serverSkuNicDBTypes, false, serverSkuNicPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSkuNic: %s", err)
	}

	count, err = ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testServerSkuNics(t *testing.T) {
	t.Parallel()

	query := ServerSkuNics()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testServerSkuNicsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuNicsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ServerSkuNics().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuNicsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSkuNicSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuNicsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ServerSkuNicExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ServerSkuNic exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ServerSkuNicExists to return true, but got false.")
	}
}

func testServerSkuNicsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	serverSkuNicFound, err := FindServerSkuNic(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if serverSkuNicFound == nil {
		t.Error("want a record, got nil")
	}
}

func testServerSkuNicsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ServerSkuNics().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testServerSkuNicsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ServerSkuNics().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testServerSkuNicsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	serverSkuNicOne := &ServerSkuNic{}
	serverSkuNicTwo := &ServerSkuNic{}
	if err = randomize.Struct(seed, serverSkuNicOne, serverSkuNicDBTypes, false, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSkuNicTwo, serverSkuNicDBTypes, false, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSkuNicOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSkuNicTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSkuNics().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testServerSkuNicsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	serverSkuNicOne := &ServerSkuNic{}
	serverSkuNicTwo := &ServerSkuNic{}
	if err = randomize.Struct(seed, serverSkuNicOne, serverSkuNicDBTypes, false, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSkuNicTwo, serverSkuNicDBTypes, false, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSkuNicOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSkuNicTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func serverSkuNicBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func serverSkuNicAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuNic) error {
	*o = ServerSkuNic{}
	return nil
}

func testServerSkuNicsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ServerSkuNic{}
	o := &ServerSkuNic{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic object: %s", err)
	}

	AddServerSkuNicHook(boil.BeforeInsertHook, serverSkuNicBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	serverSkuNicBeforeInsertHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.AfterInsertHook, serverSkuNicAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	serverSkuNicAfterInsertHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.AfterSelectHook, serverSkuNicAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	serverSkuNicAfterSelectHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.BeforeUpdateHook, serverSkuNicBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	serverSkuNicBeforeUpdateHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.AfterUpdateHook, serverSkuNicAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	serverSkuNicAfterUpdateHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.BeforeDeleteHook, serverSkuNicBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	serverSkuNicBeforeDeleteHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.AfterDeleteHook, serverSkuNicAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	serverSkuNicAfterDeleteHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.BeforeUpsertHook, serverSkuNicBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	serverSkuNicBeforeUpsertHooks = []ServerSkuNicHook{}

	AddServerSkuNicHook(boil.AfterUpsertHook, serverSkuNicAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	serverSkuNicAfterUpsertHooks = []ServerSkuNicHook{}
}

func testServerSkuNicsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSkuNicsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(serverSkuNicColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSkuNicToOneServerSkuUsingSku(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ServerSkuNic
	var foreign ServerSku

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, serverSkuNicDBTypes, false, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, serverSkuDBTypes, false, serverSkuColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSku struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SkuID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Sku().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	ranAfterSelectHook := false
	AddServerSkuHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *ServerSku) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := ServerSkuNicSlice{&local}
	if err = local.L.LoadSku(ctx, tx, false, (*[]*ServerSkuNic)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Sku == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Sku = nil
	if err = local.L.LoadSku(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Sku == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testServerSkuNicToOneSetOpServerSkuUsingSku(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ServerSkuNic
	var b, c ServerSku

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverSkuNicDBTypes, false, strmangle.SetComplement(serverSkuNicPrimaryKeyColumns, serverSkuNicColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, serverSkuDBTypes, false, strmangle.SetComplement(serverSkuPrimaryKeyColumns, serverSkuColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, serverSkuDBTypes, false, strmangle.SetComplement(serverSkuPrimaryKeyColumns, serverSkuColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ServerSku{&b, &c} {
		err = a.SetSku(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Sku != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.SkuServerSkuNics[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SkuID != x.ID {
			t.Error("foreign key was wrong value", a.SkuID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SkuID))
		reflect.Indirect(reflect.ValueOf(&a.SkuID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SkuID != x.ID {
			t.Error("foreign key was wrong value", a.SkuID, x.ID)
		}
	}
}

func testServerSkuNicsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testServerSkuNicsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSkuNicSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testServerSkuNicsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSkuNics().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	serverSkuNicDBTypes = map[string]string{`ID`: `uuid`, `SkuID`: `uuid`, `Vendor`: `string`, `Model`: `string`, `PortBandwidth`: `int8`, `PortCount`: `int8`, `Count`: `int8`, `CreatedAt`: `timestamptz`, `UpdatedAt`: `timestamptz`}
	_                   = bytes.MinRead
)

func testServerSkuNicsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(serverSkuNicPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(serverSkuNicAllColumns) == len(serverSkuNicPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testServerSkuNicsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(serverSkuNicAllColumns) == len(serverSkuNicPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuNic{}
	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuNics().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSkuNicDBTypes, true, serverSkuNicPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuNic struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(serverSkuNicAllColumns, serverSkuNicPrimaryKeyColumns) {
		fields = serverSkuNicAllColumns
	} else {
		fields = strmangle.SetComplement(
			serverSkuNicAllColumns,
			serverSkuNicPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := ServerSkuNicSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
