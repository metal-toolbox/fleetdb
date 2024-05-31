-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.component_firmware_version ADD COLUMN oem BOOL NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.component_firmware_version DROP COLUMN oem;
-- +goose StatementEnd