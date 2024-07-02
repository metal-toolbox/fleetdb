-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.event_history (
    event_id UUID NOT NULL PRIMARY KEY,
    event_type STRING NOT NULL,
    event_start TIMESTAMPTZ NOT NULL,
    event_end TIMESTAMPTZ NOT NULL,
    target_server UUID NOT NULL REFERENCES public.servers(id) ON DELETE CASCADE,
    parameters JSON,
    final_state STRING NOT NULL,
    final_status JSON
);

CREATE INDEX evt_history_target ON public.event_history (target_server ASC) INCLUDE (event_type, event_start, event_end);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.event_history;
DROP INDEX evt_history_target;
-- +goose StatementEnd
