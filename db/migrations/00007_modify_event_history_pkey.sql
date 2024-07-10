-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.event_history DROP CONSTRAINT event_history_pkey;

ALTER TABLE public.event_history ADD PRIMARY KEY (event_id, event_type, target_server);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.event_history DROP CONSTRAINT event_history_pkey;

ALTER TABLE public.event_history ADD PRIMARY KEY (event_id);
-- +goose StatementEnd
