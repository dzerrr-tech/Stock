CREATE TABLE PRODUCTS (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sku         TEXT NOT NULL,
    name        TEXT NOT NULL,
    stock       INTEGER NOT NULL DEFAULT 0,
    min_stock   INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT products_sku_unique UNIQUE (sku),
    CONSTRAINT products_sku_not_blank CHECK (btrim(sku) <> ''),
    CONSTRAINT products_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT products_stock_non_negative CHECK (stock >= 0),
    CONSTRAINT products_min_stock_non_negative CHECK (min_stock >= 0)
);