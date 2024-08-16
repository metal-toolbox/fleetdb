-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.firmware_set_validation_facts(
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    firmware_set_id UUID NOT NULL REFERENCES public.component_firmware_set(id) ON DELETE CASCADE,
    target_server_id UUID NOT NULL,
    performed_on TIMESTAMPTZ NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.firmware_set_validation_facts
-- +goose StatementEnd
