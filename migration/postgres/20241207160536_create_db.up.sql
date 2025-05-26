DO
$$
    BEGIN
        CREATE TYPE metric_type_enum AS ENUM ('gauge', 'counter');
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN
    END
$$;

COMMENT ON TYPE metric_type_enum is 'Тип метрики: gauge или counter';


CREATE TABLE IF NOT EXISTS metric_type
(
    ID          SMALLSERIAL,
    metric_type metric_type_enum NOT NULL UNIQUE,
    PRIMARY KEY (ID)
);

COMMENT ON TABLE metric_type is 'Таблица типов метрик';
COMMENT ON COLUMN metric_type.ID is 'Автоматически увеличивающийся идентификатор-счётчик';
COMMENT ON COLUMN metric_type.metric_type is 'Тип метрики: gauge или counter';

INSERT INTO metric_type (metric_type)
VALUES ('gauge');
INSERT INTO metric_type (metric_type)
VALUES ('counter');


CREATE TABLE IF NOT EXISTS metrics
(
    ID             VARCHAR(100) UNIQUE,
    delta BIGINT NULL,
    value          DOUBLE PRECISION NULL,
    metric_type_id SMALLINT,
    PRIMARY KEY (ID),
    CONSTRAINT fk_metric_type FOREIGN KEY (metric_type_id) REFERENCES metric_type (ID) ON DELETE RESTRICT ON UPDATE CASCADE
);

COMMENT ON TABLE metrics is 'Таблица метрик';
COMMENT ON COLUMN metrics.ID is 'Название метрики';
COMMENT ON COLUMN metrics.delta is 'Целочисленное значение метрики типа counter';
COMMENT ON COLUMN metrics.value is 'Вещественное значение метрики типа gauge';


CREATE OR REPLACE FUNCTION check_metric_value()
    RETURNS TRIGGER AS
$$
DECLARE
    mtype metric_type_enum;
BEGIN
    SELECT metric_type
    INTO mtype
    FROM metric_type
    WHERE ID = NEW.metric_type_id;

    IF mtype = 'counter' AND (NEW.delta IS NULL OR NEW.value IS NOT NULL) THEN
        RAISE EXCEPTION 'For counter type, delta must not be NULL and value must be NULL';
    ELSIF mtype = 'gauge' AND (NEW.value IS NULL OR NEW.delta IS NOT NULL) THEN
        RAISE EXCEPTION 'For gauge type, value must not be NULL and delta must be NULL';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1
                       FROM pg_trigger
                       WHERE tgname = 'trg_check_metric_value') THEN
            CREATE TRIGGER trg_check_metric_value
                BEFORE INSERT OR UPDATE
                ON metrics
                FOR EACH ROW
            EXECUTE FUNCTION check_metric_value();

            RAISE NOTICE 'Trigger trg_check_metric_value created';
        ELSE
            RAISE NOTICE 'Trigger trg_check_metric_value already exists';
        END IF;
    END
$$;