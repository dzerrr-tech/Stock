CREATE TABLE stock_movement (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id BIGINT NOT NULL,
    type    TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    note TEXT NOT NULL DEFAULT'',
    create_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT stock_movements_product_fk FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE RESTRICT,
    CONSTRAINT stock_movements_type_valid CHECK (type IN('IN', 'OUT')),
    CONSTRAINT stock_movements_quantity_positive CHECK (quantity > 0),
    CONSTRAINT stock_movements_note_length CHECK (char_length(note) <= 500)
);

CREATE INDEX stock_movements_product_created_idx
    ON stock_movements (product_id, create_at DESC);
