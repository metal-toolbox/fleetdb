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

func testBMCMacAddressesUpsert(t *testing.T) {
	t.Parallel()

	if len(bmcMacAddressAllColumns) == len(bmcMacAddressPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := BMCMacAddress{}
	if err = randomize.Struct(seed, &o, bmcMacAddressDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BMCMacAddress: %s", err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, bmcMacAddressDBTypes, false, bmcMacAddressPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BMCMacAddress: %s", err)
	}

	count, err = BMCMacAddresses().Count(ctx, tx)
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

func testBMCMacAddresses(t *testing.T) {
	t.Parallel()

	query := BMCMacAddresses()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testBMCMacAddressesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
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

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBMCMacAddressesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := BMCMacAddresses().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBMCMacAddressesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BMCMacAddressSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBMCMacAddressesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := BMCMacAddressExists(ctx, tx, o.BMCMacAddress)
	if err != nil {
		t.Errorf("Unable to check if BMCMacAddress exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BMCMacAddressExists to return true, but got false.")
	}
}

func testBMCMacAddressesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	bmcMacAddressFound, err := FindBMCMacAddress(ctx, tx, o.BMCMacAddress)
	if err != nil {
		t.Error(err)
	}

	if bmcMacAddressFound == nil {
		t.Error("want a record, got nil")
	}
}

func testBMCMacAddressesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = BMCMacAddresses().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testBMCMacAddressesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := BMCMacAddresses().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBMCMacAddressesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	bmcMacAddressOne := &BMCMacAddress{}
	bmcMacAddressTwo := &BMCMacAddress{}
	if err = randomize.Struct(seed, bmcMacAddressOne, bmcMacAddressDBTypes, false, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}
	if err = randomize.Struct(seed, bmcMacAddressTwo, bmcMacAddressDBTypes, false, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = bmcMacAddressOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = bmcMacAddressTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BMCMacAddresses().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBMCMacAddressesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	bmcMacAddressOne := &BMCMacAddress{}
	bmcMacAddressTwo := &BMCMacAddress{}
	if err = randomize.Struct(seed, bmcMacAddressOne, bmcMacAddressDBTypes, false, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}
	if err = randomize.Struct(seed, bmcMacAddressTwo, bmcMacAddressDBTypes, false, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = bmcMacAddressOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = bmcMacAddressTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func bmcMacAddressBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func bmcMacAddressAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BMCMacAddress) error {
	*o = BMCMacAddress{}
	return nil
}

func testBMCMacAddressesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &BMCMacAddress{}
	o := &BMCMacAddress{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, false); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress object: %s", err)
	}

	AddBMCMacAddressHook(boil.BeforeInsertHook, bmcMacAddressBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressBeforeInsertHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.AfterInsertHook, bmcMacAddressAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressAfterInsertHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.AfterSelectHook, bmcMacAddressAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressAfterSelectHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.BeforeUpdateHook, bmcMacAddressBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressBeforeUpdateHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.AfterUpdateHook, bmcMacAddressAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressAfterUpdateHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.BeforeDeleteHook, bmcMacAddressBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressBeforeDeleteHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.AfterDeleteHook, bmcMacAddressAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressAfterDeleteHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.BeforeUpsertHook, bmcMacAddressBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressBeforeUpsertHooks = []BMCMacAddressHook{}

	AddBMCMacAddressHook(boil.AfterUpsertHook, bmcMacAddressAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	bmcMacAddressAfterUpsertHooks = []BMCMacAddressHook{}
}

func testBMCMacAddressesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBMCMacAddressesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(bmcMacAddressColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBMCMacAddressToOneBomInfoUsingSerialNumBomInfo(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local BMCMacAddress
	var foreign BomInfo

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, bmcMacAddressDBTypes, false, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, bomInfoDBTypes, false, bomInfoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BomInfo struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SerialNum = foreign.SerialNum
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.SerialNumBomInfo().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.SerialNum != foreign.SerialNum {
		t.Errorf("want: %v, got %v", foreign.SerialNum, check.SerialNum)
	}

	ranAfterSelectHook := false
	AddBomInfoHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *BomInfo) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := BMCMacAddressSlice{&local}
	if err = local.L.LoadSerialNumBomInfo(ctx, tx, false, (*[]*BMCMacAddress)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.SerialNumBomInfo == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.SerialNumBomInfo = nil
	if err = local.L.LoadSerialNumBomInfo(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.SerialNumBomInfo == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testBMCMacAddressToOneSetOpBomInfoUsingSerialNumBomInfo(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BMCMacAddress
	var b, c BomInfo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bmcMacAddressDBTypes, false, strmangle.SetComplement(bmcMacAddressPrimaryKeyColumns, bmcMacAddressColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, bomInfoDBTypes, false, strmangle.SetComplement(bomInfoPrimaryKeyColumns, bomInfoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, bomInfoDBTypes, false, strmangle.SetComplement(bomInfoPrimaryKeyColumns, bomInfoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*BomInfo{&b, &c} {
		err = a.SetSerialNumBomInfo(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.SerialNumBomInfo != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.SerialNumBMCMacAddresses[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SerialNum != x.SerialNum {
			t.Error("foreign key was wrong value", a.SerialNum)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SerialNum))
		reflect.Indirect(reflect.ValueOf(&a.SerialNum)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SerialNum != x.SerialNum {
			t.Error("foreign key was wrong value", a.SerialNum, x.SerialNum)
		}
	}
}

func testBMCMacAddressesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
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

func testBMCMacAddressesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BMCMacAddressSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBMCMacAddressesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BMCMacAddresses().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	bmcMacAddressDBTypes = map[string]string{`BMCMacAddress`: `string`, `SerialNum`: `string`}
	_                    = bytes.MinRead
)

func testBMCMacAddressesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(bmcMacAddressPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(bmcMacAddressAllColumns) == len(bmcMacAddressPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testBMCMacAddressesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(bmcMacAddressAllColumns) == len(bmcMacAddressPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BMCMacAddress{}
	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BMCMacAddresses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, bmcMacAddressDBTypes, true, bmcMacAddressPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BMCMacAddress struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(bmcMacAddressAllColumns, bmcMacAddressPrimaryKeyColumns) {
		fields = bmcMacAddressAllColumns
	} else {
		fields = strmangle.SetComplement(
			bmcMacAddressAllColumns,
			bmcMacAddressPrimaryKeyColumns,
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

	slice := BMCMacAddressSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
