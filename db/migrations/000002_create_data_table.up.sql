CREATE TABLE IF NOT EXISTS data_items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type text NOT NULL CHECK (type IN ('login_password', 'text', 'binary', 'bank_card')),
    data text NOT NULL,
    metadata text DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS data_items_user_id_idx ON data_items(user_id);
CREATE INDEX IF NOT EXISTS data_items_type_idx ON data_items(type);
CREATE INDEX IF NOT EXISTS data_items_created_at_idx ON data_items(created_at);

COMMENT ON TABLE data_items IS 'Таблица для хранения пользовательских данных';
COMMENT ON COLUMN data_items.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN data_items.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN data_items.type IS 'Тип данных: login_password, text, binary, bank_card';
COMMENT ON COLUMN data_items.data IS 'Данные в JSON формате';
COMMENT ON COLUMN data_items.metadata IS 'Метаинформация о данных';
COMMENT ON COLUMN data_items.created_at IS 'Дата создания записи';
COMMENT ON COLUMN data_items.updated_at IS 'Дата последнего обновления записи';
