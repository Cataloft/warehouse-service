-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS warehouses
(
    id           bigserial
        constraint warehouses_pk
            primary key,
    name         varchar(50) not null,
    availability boolean     not null
);

CREATE INDEX IF NOT EXISTS idx_warehouses_availability
    ON warehouses(availability);
-- +goose StatementEnd