CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

CREATE TABLE IF NOT EXISTS metrics (
    ID SERIAL NOT NULL,
    metric_type metric_type NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    delta INTEGER NULL,
    value DOUBLE PRECISION NULL,
    PRIMARY KEY (ID),
    CONSTRAINT chk_one_value CHECK (
        (delta IS NOT NULL AND value IS NULL) OR
        (delta IS NULL AND value IS NOT NULL)
    )
);

COMMENT ON TABLE metrics is 'Таблица метрик';
COMMENT ON COLUMN metrics.ID is 'Индекнтификатор метрики';
COMMENT ON COLUMN metrics.metric_type is 'Тип метрики: gauge или counter';
COMMENT ON COLUMN metrics.metric_name is 'Название метрики';
COMMENT ON COLUMN metrics.delta is 'Целочисленное значение метрики';
COMMENT ON COLUMN metrics.value is 'Вещественное значение метрики';