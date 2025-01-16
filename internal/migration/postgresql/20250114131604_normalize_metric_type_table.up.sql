
-- По-правильному, стоило бы мигрировть через temp-таблицу, т.к. данные уже могут быть в таблице на момент второй миграции.
-- (Потому что старый строковый uuid не приести к новому типу поля-идентификатора: INT)
-- Здесь такой подход игнорируется, потому что данные можно спокойно удалить.

DELETE FROM metric_type;
DELETE FROM metric;

ALTER TABLE IF EXISTS metric_type
    ALTER COLUMN id SET DATA TYPE INT USING (id::INT),
    ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY,
    ADD CONSTRAINT unique_metric_type UNIQUE (metric_type),
    DROP CONSTRAINT IF EXISTS fk_metric_type;

ALTER TABLE IF EXISTS metric
    ADD COLUMN metric_type_id INT NOT NULL REFERENCES metric_type(id),
    ADD CONSTRAINT fk_metric FOREIGN KEY (metric_type_id) REFERENCES metric_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;

COMMENT ON COLUMN metric.metric_type_id is 'Индекнтификатор типа метрики из таблицы metric_type';
