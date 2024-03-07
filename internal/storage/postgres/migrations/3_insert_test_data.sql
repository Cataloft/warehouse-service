-- +goose Up
-- +goose StatementBegin
INSERT INTO warehouses(id, name, availability)
VALUES
    (1, 'Склад 1', true),
    (2, 'Склад 2', false),
    (3, 'Склад 3', true),
    (4, 'Склад 4', false),
    (5, 'Склад 5', true);

INSERT INTO goods(id, warehouse_id, name, unique_code, amount)
VALUES
    (1, 1, 'Adidas', '12345678', 115),
    (2, 1, 'Nike', '22345678', 105),
    (3, 1, 'Lewis', '32345678', 40),
    (4, 2, 'Reebok', '42345678', 95),
    (5, 3, 'Asics', '52345678', 35),
    (6, 3, 'NB', '62345678', 100),
    (7, 4, 'Puma', '72345678', 70),
    (8, 5, 'Balenciaga', '82345678', 25);
-- +goose StatementEnd
