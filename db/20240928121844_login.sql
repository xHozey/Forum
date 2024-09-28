-- +goose Up
-- +goose StatementBegin
CREATE TABLE login (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mail TEXT NOT NULL,
    pasword TEXT NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE login;
-- +goose StatementEnd
