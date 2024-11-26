CREATE OR REPLACE FUNCTION delete_empty_group()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем, остались ли песни в группе
    IF NOT EXISTS (
        SELECT 1 FROM "Song" WHERE group_id = OLD.group_id
    ) THEN
        -- Удаляем группу, если песен больше нет
        DELETE FROM "Group" WHERE group_id = OLD.group_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_empty_group
AFTER DELETE ON "Song"
FOR EACH ROW
EXECUTE FUNCTION delete_empty_group();