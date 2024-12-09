CREATE TYPE metric_type_enum AS ENUM ('gauge', 'counter');

COMMENT ON TYPE metric_type_enum is 'Тип метрики: gauge или counter';

CREATE TABLE IF NOT EXISTS metric_type (
    ID VARCHAR(36) PRIMARY KEY,
    metric_type metric_type_enum NOT NULL UNIQUE
);

COMMENT ON TABLE metric_type is 'Таблица типов метрик';
COMMENT ON COLUMN metric_type.ID is 'Индекнтификатор метрики';
COMMENT ON COLUMN metric_type.metric_type is 'Тип метрики: gauge или counter';

CREATE TABLE IF NOT EXISTS metric (
    ID VARCHAR(36) PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    delta INTEGER NULL,
    value DOUBLE PRECISION NULL,
    CONSTRAINT chk_one_value CHECK (
        (delta IS NOT NULL AND value IS NULL) OR
        (delta IS NULL AND value IS NOT NULL)
    ),
    CONSTRAINT fk_metric_type FOREIGN KEY (ID) REFERENCES metric_type(ID) ON DELETE RESTRICT ON UPDATE CASCADE
);

COMMENT ON TABLE metric is 'Таблица метрик';
COMMENT ON COLUMN metric.ID is 'Индекнтификатор метрики';
COMMENT ON COLUMN metric.metric_name is 'Название метрики';
COMMENT ON COLUMN metric.delta is 'Целочисленное значение метрики';
COMMENT ON COLUMN metric.value is 'Вещественное значение метрики';

CREATE TABLE IF NOT EXISTS uuid (
    UUID VARCHAR(36) UNIQUE,
    is_free BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now()
);

COMMENT ON TABLE uuid is 'Таблица для генерации уникальных идентификаторов метрик';
COMMENT ON COLUMN uuid.UUID is 'Уникальный индекнтификатор метрики';
COMMENT ON COLUMN uuid.is_free is 'Флаг: свободен ли UUID';
COMMENT ON COLUMN uuid.created_at is 'Timestamp создания UUID';

CREATE OR REPLACE FUNCTION get_first_free_generated_metric_UUID()
RETURNS VARCHAR(36)
LANGUAGE plpgsql
AS
$$
DECLARE
    searched_UUID VARCHAR(36);
BEGIN
    SELECT uuid
    INTO searched_UUID
    FROM uuid
    WHERE is_free = true
    ORDER BY created_at ASC
    LIMIT 1;
    return searched_UUID;
END;
$$;

COMMENT ON FUNCTION get_first_free_generated_metric_UUID is 'Функция, возвращающая первый свободный сгенерированный UUID для метрики';

CREATE OR REPLACE FUNCTION get_metric_uuid()
RETURNS VARCHAR(36)
LANGUAGE plpgsql
AS
$$
DECLARE
    searched_UUID VARCHAR(36);
	UUID_from_table VARCHAR(36);
    is_set BOOLEAN;
BEGIN
    SELECT get_first_free_generated_metric_UUID() INTO searched_UUID;
    IF searched_UUID IS NOT NULL THEN
        return searched_UUID;
    END IF;

    FOR i IN 1..10 LOOP
        SELECT FALSE INTO is_set;
        WHILE NOT is_set LOOP
            SELECT gen_random_uuid() INTO searched_UUID;
            BEGIN
                INSERT INTO uuid("UUID") VALUES (searched_UUID);
                EXCEPTION WHEN OTHERS THEN
                    CONTINUE;
            END;
            SELECT TRUE INTO is_set;
        END LOOP;
    END LOOP;

    return searched_UUID;
END;
$$;

COMMENT ON FUNCTION get_metric_uuid is 'Функция, возвращающая первый свободный сгенерированный UUID для метрики. При его отсутствии добавлет в таблицу 10 новых свободных UUID';
