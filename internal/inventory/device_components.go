//nolint:all  // XXX remove this!
package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (dv *DeviceView) ComposeComponents(ctx context.Context, exec boil.ContextExecutor,
	srv uuid.UUID, method string) error {
	if err := dv.writeDimms(ctx, exec, srv, method); err != nil {
		return err
	}
	return nil
}

func (dv *DeviceView) writeDimms(ctx context.Context, exec boil.ContextExecutor,
	srv uuid.UUID, method string) error {
	return errors.New("unimplemented")
}
