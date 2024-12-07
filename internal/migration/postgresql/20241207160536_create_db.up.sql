CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

CREATE TABLE IF NOT EXISTS metric_types (
    ID SERIAL PRIMARY KEY,
    metric_type metric_type NOT NULL UNIQUE
);

COMMENT ON TABLE metric_types is 'Таблица типов метрик';
COMMENT ON COLUMN metric_types.ID is 'Индекнтификатор метрики';
COMMENT ON COLUMN metric_types.metric_type is 'Тип метрики: gauge или counter';

CREATE TABLE IF NOT EXISTS metric (
    ID SERIAL PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    delta INTEGER NULL,
    value DOUBLE PRECISION NULL,
    CONSTRAINT chk_one_value CHECK (
        (delta IS NOT NULL AND value IS NULL) OR
        (delta IS NULL AND value IS NOT NULL)
    ),
    CONSTRAINT fk_metric_type FOREIGN KEY (ID) REFERENCES metric_types(ID) ON DELETE RESTRICT ON UPDATE CASCADE
);

COMMENT ON TABLE metric is 'Таблица метрик';
COMMENT ON COLUMN metric.ID is 'Индекнтификатор метрики';
COMMENT ON COLUMN metric.metric_name is 'Название метрики';
COMMENT ON COLUMN metric.delta is 'Целочисленное значение метрики';
COMMENT ON COLUMN metric.value is 'Вещественное значение метрики';