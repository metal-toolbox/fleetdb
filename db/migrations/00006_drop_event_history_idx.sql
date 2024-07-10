-- +goose Up
-- +goose StatementBegin
DROP INDEX evt_history_target;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE INDEX evt_history_target ON public.event_history (target_server ASC) INCLUDE (event_type, event_start, event_end);
-- +goose StatementEnd
