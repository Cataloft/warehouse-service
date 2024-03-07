-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods
(
    id           bigserial
        constraint goods_pk
            primary key,
    warehouse_id bigint      not null
        references warehouses,
    name         varchar(50) not null,
    unique_code  varchar(8)  not null,
    amount       bigint      not null
);

CREATE INDEX IF NOT EXISTS idx_goods_unique_code
    ON goods (unique_code);
-- +goose StatementEnd