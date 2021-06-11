DROP TRIGGER IF EXISTS todo_before_update ON todo;
DROP TABLE IF EXISTS todo;
DROP FUNCTION IF EXISTS todo_touch_updated_at;

CREATE FUNCTION todo_touch_updated_at() RETURNS TRIGGER AS $func$
    BEGIN
        IF NEW IS DISTINCT FROM OLD THEN
            NEW.updated_at = now();
        END IF;
        RETURN NEW;
    END;
$func$ LANGUAGE plpgsql;

CREATE TABLE todo(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    text text NOT NULL,
    done boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER todo_before_update BEFORE UPDATE
ON todo
FOR EACH ROW
EXECUTE PROCEDURE todo_touch_updated_at();
