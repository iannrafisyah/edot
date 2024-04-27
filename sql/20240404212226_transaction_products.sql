-- +goose Up
create table transaction_products (
    id              bigserial primary key,
    name            varchar(50),
    price           float,
    qty             int,
    product_id      int,
    transaction_id  int,
    warehouse_id    int,
    to_warehouse_id int default null,
    snapshot        jsonb,
    updated_at      timestamptz default now(),
    created_at      timestamptz default now(),
    foreign key (product_id) references products (id),
    foreign key (transaction_id) references transactions (id),
    foreign key (warehouse_id) references warehouses (id),
    foreign key (to_warehouse_id) references warehouses (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
drop table transaction_products;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
