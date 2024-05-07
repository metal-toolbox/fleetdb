-- +goose Up
-- +goose StatementBegin

CREATE TABLE public.bios_config_sets (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING UNIQUE NOT NULL,
    version STRING NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE public.bios_config_components (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    fk_bios_config_set_id UUID NOT NULL REFERENCES public.bios_config_sets(id) ON DELETE CASCADE,
    name STRING NOT NULL,
    vendor STRING NOT NULL,
    model STRING NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    UNIQUE (fk_bios_config_set_id, name)
);

CREATE TABLE public.bios_config_settings (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    fk_bios_config_component_id UUID NOT NULL REFERENCES public.bios_config_components(id) ON DELETE CASCADE,
    settings_key STRING NOT NULL,
    settings_value STRING NOT NULL,
    raw JSONB NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    UNIQUE (fk_bios_config_component_id, settings_key)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE public.bios_config_settings;
DROP TABLE public.bios_config_components;
DROP TABLE public.bios_config_sets;

-- +goose StatementEnd