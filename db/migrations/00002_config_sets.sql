-- +goose Up
-- +goose StatementBegin

CREATE TABLE public.config_sets (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING  UNIQUE NOT NULL,
    version STRING NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE public.config_components (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    fk_config_set_id UUID NOT NULL REFERENCES public.config_sets(id) ON DELETE CASCADE,
    name STRING NOT NULL,
    vendor STRING NULL,
    model STRING NULL,
    serial STRING NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    UNIQUE (fk_config_set_id, name)
);

CREATE TABLE public.config_component_settings (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    fk_component_id UUID NOT NULL REFERENCES public.config_components(id) ON DELETE CASCADE,
    settings_key STRING NOT NULL,
    settings_value STRING NOT NULL,
    custom JSONB NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    UNIQUE (fk_component_id, settings_key)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE public.config_component_settings;
DROP TABLE public.config_components;
DROP TABLE public.config_sets;

-- +goose StatementEnd