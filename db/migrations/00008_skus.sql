-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.server_sku (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    version STRING NOT NULL,
    vendor STRING NOT NULL,
    chassis STRING NOT NULL,
    bmc_model STRING NOT NULL,
    motherboard_model STRING NOT NULL,
    cpu_vendor STRING NOT NULL,
    cpu_model STRING NOT NULL,
    cpu_cores INTEGER NOT NULL,
    cpu_hertz BIGINT NOT NULL,
    cpu_count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    UNIQUE (name, version)
);

CREATE TABLE public.server_sku_disk (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES public.server_sku(id) ON DELETE CASCADE,
    bytes BIGINT NOT NULL,
    protocol STRING NOT NULL, -- SATA vs NVMe vs PCIE
    count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE public.server_sku_memory (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES public.server_sku(id) ON DELETE CASCADE,
    bytes BIGINT NOT NULL,
    count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE public.server_sku_nic (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES public.server_sku(id) ON DELETE CASCADE,
    port_bandwidth BIGINT NOT NULL,
    port_count INTEGER NOT NULL,
    count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE public.server_sku_aux_device (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES public.server_sku(id) ON DELETE CASCADE,
    vendor STRING NOT NULL,
    model STRING NOT NULL,
    device_type STRING NOT NULL, -- GPU vs. other?
    details JSON NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.server_sku;
DROP TABLE public.server_sku_disk;
DROP TABLE public.server_sku_memory;
DROP TABLE public.server_sku_nic;
DROP TABLE public.server_sku_aux_device;
-- +goose StatementEnd