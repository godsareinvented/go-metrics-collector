DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS metric_type;
DROP TYPE IF EXISTS metric_type_enum;
DROP FUNCTION IF EXISTS check_metric_value;
DO
$$
    BEGIN
        IF EXISTS (SELECT 1
                   FROM information_schema.tables
                   WHERE table_name = 'metrics') THEN
            IF EXISTS (SELECT 1
                       FROM pg_trigger
                       WHERE tgname = 'trg_check_metric_value') THEN
                EXECUTE 'DROP TRIGGER IF EXISTS trg_check_metric_value ON metrics';
            END IF;
        ELSE
            RAISE NOTICE 'Table metrics does not exist';
        END IF;
    END
$$;