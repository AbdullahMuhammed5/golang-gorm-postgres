DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type where pg_type.typname = 'ticket_status'
    ) THEN
        CREATE TYPE ticket_status AS ENUM (
            'NEW',
            'IN_PROGRESS',
            'RESOLVED'
        );
    END IF;
END
$$