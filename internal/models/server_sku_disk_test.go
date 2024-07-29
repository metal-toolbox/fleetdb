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

func testServerSkuDisksUpsert(t *testing.T) {
	t.Parallel()

	if len(serverSkuDiskAllColumns) == len(serverSkuDiskPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ServerSkuDisk{}
	if err = randomize.Struct(seed, &o, serverSkuDiskDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSkuDisk: %s", err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, serverSkuDiskDBTypes, false, serverSkuDiskPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSkuDisk: %s", err)
	}

	count, err = ServerSkuDisks().Count(ctx, tx)
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

func testServerSkuDisks(t *testing.T) {
	t.Parallel()

	query := ServerSkuDisks()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testServerSkuDisksDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
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

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuDisksQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ServerSkuDisks().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuDisksSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSkuDiskSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSkuDisksExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ServerSkuDiskExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ServerSkuDisk exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ServerSkuDiskExists to return true, but got false.")
	}
}

func testServerSkuDisksFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	serverSkuDiskFound, err := FindServerSkuDisk(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if serverSkuDiskFound == nil {
		t.Error("want a record, got nil")
	}
}

func testServerSkuDisksBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ServerSkuDisks().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testServerSkuDisksOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ServerSkuDisks().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testServerSkuDisksAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	serverSkuDiskOne := &ServerSkuDisk{}
	serverSkuDiskTwo := &ServerSkuDisk{}
	if err = randomize.Struct(seed, serverSkuDiskOne, serverSkuDiskDBTypes, false, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSkuDiskTwo, serverSkuDiskDBTypes, false, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSkuDiskOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSkuDiskTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSkuDisks().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testServerSkuDisksCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	serverSkuDiskOne := &ServerSkuDisk{}
	serverSkuDiskTwo := &ServerSkuDisk{}
	if err = randomize.Struct(seed, serverSkuDiskOne, serverSkuDiskDBTypes, false, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSkuDiskTwo, serverSkuDiskDBTypes, false, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSkuDiskOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSkuDiskTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func serverSkuDiskBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func serverSkuDiskAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSkuDisk) error {
	*o = ServerSkuDisk{}
	return nil
}

func testServerSkuDisksHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ServerSkuDisk{}
	o := &ServerSkuDisk{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk object: %s", err)
	}

	AddServerSkuDiskHook(boil.BeforeInsertHook, serverSkuDiskBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskBeforeInsertHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.AfterInsertHook, serverSkuDiskAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskAfterInsertHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.AfterSelectHook, serverSkuDiskAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskAfterSelectHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.BeforeUpdateHook, serverSkuDiskBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskBeforeUpdateHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.AfterUpdateHook, serverSkuDiskAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskAfterUpdateHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.BeforeDeleteHook, serverSkuDiskBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskBeforeDeleteHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.AfterDeleteHook, serverSkuDiskAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskAfterDeleteHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.BeforeUpsertHook, serverSkuDiskBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskBeforeUpsertHooks = []ServerSkuDiskHook{}

	AddServerSkuDiskHook(boil.AfterUpsertHook, serverSkuDiskAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	serverSkuDiskAfterUpsertHooks = []ServerSkuDiskHook{}
}

func testServerSkuDisksInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSkuDisksInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(serverSkuDiskColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSkuDiskToOneServerSkuUsingSku(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ServerSkuDisk
	var foreign ServerSku

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, serverSkuDiskDBTypes, false, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
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

	slice := ServerSkuDiskSlice{&local}
	if err = local.L.LoadSku(ctx, tx, false, (*[]*ServerSkuDisk)(&slice), nil); err != nil {
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

func testServerSkuDiskToOneSetOpServerSkuUsingSku(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ServerSkuDisk
	var b, c ServerSku

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverSkuDiskDBTypes, false, strmangle.SetComplement(serverSkuDiskPrimaryKeyColumns, serverSkuDiskColumnsWithoutDefault)...); err != nil {
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

		if x.R.SkuServerSkuDisks[0] != &a {
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

func testServerSkuDisksReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
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

func testServerSkuDisksReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSkuDiskSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testServerSkuDisksSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSkuDisks().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	serverSkuDiskDBTypes = map[string]string{`ID`: `uuid`, `SkuID`: `uuid`, `Vendor`: `string`, `Model`: `string`, `Bytes`: `int8`, `Protocol`: `string`, `Count`: `int8`, `CreatedAt`: `timestamptz`, `UpdatedAt`: `timestamptz`}
	_                    = bytes.MinRead
)

func testServerSkuDisksUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(serverSkuDiskPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(serverSkuDiskAllColumns) == len(serverSkuDiskPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testServerSkuDisksSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(serverSkuDiskAllColumns) == len(serverSkuDiskPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSkuDisk{}
	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSkuDisks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSkuDiskDBTypes, true, serverSkuDiskPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSkuDisk struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(serverSkuDiskAllColumns, serverSkuDiskPrimaryKeyColumns) {
		fields = serverSkuDiskAllColumns
	} else {
		fields = strmangle.SetComplement(
			serverSkuDiskAllColumns,
			serverSkuDiskPrimaryKeyColumns,
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

	slice := ServerSkuDiskSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
