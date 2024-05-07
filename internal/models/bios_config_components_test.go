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

func testBiosConfigComponentsUpsert(t *testing.T) {
	t.Parallel()

	if len(biosConfigComponentAllColumns) == len(biosConfigComponentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := BiosConfigComponent{}
	if err = randomize.Struct(seed, &o, biosConfigComponentDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BiosConfigComponent: %s", err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, biosConfigComponentDBTypes, false, biosConfigComponentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BiosConfigComponent: %s", err)
	}

	count, err = BiosConfigComponents().Count(ctx, tx)
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

func testBiosConfigComponents(t *testing.T) {
	t.Parallel()

	query := BiosConfigComponents()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testBiosConfigComponentsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
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

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBiosConfigComponentsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := BiosConfigComponents().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBiosConfigComponentsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BiosConfigComponentSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBiosConfigComponentsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := BiosConfigComponentExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if BiosConfigComponent exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BiosConfigComponentExists to return true, but got false.")
	}
}

func testBiosConfigComponentsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	biosConfigComponentFound, err := FindBiosConfigComponent(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if biosConfigComponentFound == nil {
		t.Error("want a record, got nil")
	}
}

func testBiosConfigComponentsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = BiosConfigComponents().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testBiosConfigComponentsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := BiosConfigComponents().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBiosConfigComponentsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	biosConfigComponentOne := &BiosConfigComponent{}
	biosConfigComponentTwo := &BiosConfigComponent{}
	if err = randomize.Struct(seed, biosConfigComponentOne, biosConfigComponentDBTypes, false, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}
	if err = randomize.Struct(seed, biosConfigComponentTwo, biosConfigComponentDBTypes, false, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = biosConfigComponentOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = biosConfigComponentTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BiosConfigComponents().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBiosConfigComponentsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	biosConfigComponentOne := &BiosConfigComponent{}
	biosConfigComponentTwo := &BiosConfigComponent{}
	if err = randomize.Struct(seed, biosConfigComponentOne, biosConfigComponentDBTypes, false, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}
	if err = randomize.Struct(seed, biosConfigComponentTwo, biosConfigComponentDBTypes, false, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = biosConfigComponentOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = biosConfigComponentTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func biosConfigComponentBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func biosConfigComponentAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BiosConfigComponent) error {
	*o = BiosConfigComponent{}
	return nil
}

func testBiosConfigComponentsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &BiosConfigComponent{}
	o := &BiosConfigComponent{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, false); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent object: %s", err)
	}

	AddBiosConfigComponentHook(boil.BeforeInsertHook, biosConfigComponentBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentBeforeInsertHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.AfterInsertHook, biosConfigComponentAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentAfterInsertHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.AfterSelectHook, biosConfigComponentAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentAfterSelectHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.BeforeUpdateHook, biosConfigComponentBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentBeforeUpdateHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.AfterUpdateHook, biosConfigComponentAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentAfterUpdateHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.BeforeDeleteHook, biosConfigComponentBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentBeforeDeleteHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.AfterDeleteHook, biosConfigComponentAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentAfterDeleteHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.BeforeUpsertHook, biosConfigComponentBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentBeforeUpsertHooks = []BiosConfigComponentHook{}

	AddBiosConfigComponentHook(boil.AfterUpsertHook, biosConfigComponentAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	biosConfigComponentAfterUpsertHooks = []BiosConfigComponentHook{}
}

func testBiosConfigComponentsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBiosConfigComponentsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(biosConfigComponentColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBiosConfigComponentToManyFKBiosConfigComponentBiosConfigSettings(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BiosConfigComponent
	var b, c BiosConfigSetting

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, biosConfigSettingDBTypes, false, biosConfigSettingColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, biosConfigSettingDBTypes, false, biosConfigSettingColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.FKBiosConfigComponentID = a.ID
	c.FKBiosConfigComponentID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.FKBiosConfigComponentBiosConfigSettings().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.FKBiosConfigComponentID == b.FKBiosConfigComponentID {
			bFound = true
		}
		if v.FKBiosConfigComponentID == c.FKBiosConfigComponentID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BiosConfigComponentSlice{&a}
	if err = a.L.LoadFKBiosConfigComponentBiosConfigSettings(ctx, tx, false, (*[]*BiosConfigComponent)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FKBiosConfigComponentBiosConfigSettings); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.FKBiosConfigComponentBiosConfigSettings = nil
	if err = a.L.LoadFKBiosConfigComponentBiosConfigSettings(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FKBiosConfigComponentBiosConfigSettings); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBiosConfigComponentToManyAddOpFKBiosConfigComponentBiosConfigSettings(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BiosConfigComponent
	var b, c, d, e BiosConfigSetting

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, biosConfigComponentDBTypes, false, strmangle.SetComplement(biosConfigComponentPrimaryKeyColumns, biosConfigComponentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*BiosConfigSetting{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, biosConfigSettingDBTypes, false, strmangle.SetComplement(biosConfigSettingPrimaryKeyColumns, biosConfigSettingColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*BiosConfigSetting{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFKBiosConfigComponentBiosConfigSettings(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FKBiosConfigComponentID {
			t.Error("foreign key was wrong value", a.ID, first.FKBiosConfigComponentID)
		}
		if a.ID != second.FKBiosConfigComponentID {
			t.Error("foreign key was wrong value", a.ID, second.FKBiosConfigComponentID)
		}

		if first.R.FKBiosConfigComponent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.FKBiosConfigComponent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.FKBiosConfigComponentBiosConfigSettings[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.FKBiosConfigComponentBiosConfigSettings[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.FKBiosConfigComponentBiosConfigSettings().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testBiosConfigComponentToOneBiosConfigSetUsingFKBiosConfigSet(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local BiosConfigComponent
	var foreign BiosConfigSet

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, biosConfigComponentDBTypes, false, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, biosConfigSetDBTypes, false, biosConfigSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigSet struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.FKBiosConfigSetID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.FKBiosConfigSet().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	ranAfterSelectHook := false
	AddBiosConfigSetHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *BiosConfigSet) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := BiosConfigComponentSlice{&local}
	if err = local.L.LoadFKBiosConfigSet(ctx, tx, false, (*[]*BiosConfigComponent)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FKBiosConfigSet == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.FKBiosConfigSet = nil
	if err = local.L.LoadFKBiosConfigSet(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FKBiosConfigSet == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testBiosConfigComponentToOneSetOpBiosConfigSetUsingFKBiosConfigSet(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BiosConfigComponent
	var b, c BiosConfigSet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, biosConfigComponentDBTypes, false, strmangle.SetComplement(biosConfigComponentPrimaryKeyColumns, biosConfigComponentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, biosConfigSetDBTypes, false, strmangle.SetComplement(biosConfigSetPrimaryKeyColumns, biosConfigSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, biosConfigSetDBTypes, false, strmangle.SetComplement(biosConfigSetPrimaryKeyColumns, biosConfigSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*BiosConfigSet{&b, &c} {
		err = a.SetFKBiosConfigSet(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.FKBiosConfigSet != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.FKBiosConfigSetBiosConfigComponents[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.FKBiosConfigSetID != x.ID {
			t.Error("foreign key was wrong value", a.FKBiosConfigSetID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FKBiosConfigSetID))
		reflect.Indirect(reflect.ValueOf(&a.FKBiosConfigSetID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.FKBiosConfigSetID != x.ID {
			t.Error("foreign key was wrong value", a.FKBiosConfigSetID, x.ID)
		}
	}
}

func testBiosConfigComponentsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
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

func testBiosConfigComponentsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BiosConfigComponentSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBiosConfigComponentsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BiosConfigComponents().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	biosConfigComponentDBTypes = map[string]string{`ID`: `uuid`, `FKBiosConfigSetID`: `uuid`, `Name`: `string`, `Vendor`: `string`, `Model`: `string`, `CreatedAt`: `timestamptz`, `UpdatedAt`: `timestamptz`}
	_                          = bytes.MinRead
)

func testBiosConfigComponentsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(biosConfigComponentPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(biosConfigComponentAllColumns) == len(biosConfigComponentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testBiosConfigComponentsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(biosConfigComponentAllColumns) == len(biosConfigComponentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BiosConfigComponent{}
	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BiosConfigComponents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, biosConfigComponentDBTypes, true, biosConfigComponentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BiosConfigComponent struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(biosConfigComponentAllColumns, biosConfigComponentPrimaryKeyColumns) {
		fields = biosConfigComponentAllColumns
	} else {
		fields = strmangle.SetComplement(
			biosConfigComponentAllColumns,
			biosConfigComponentPrimaryKeyColumns,
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

	slice := BiosConfigComponentSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
