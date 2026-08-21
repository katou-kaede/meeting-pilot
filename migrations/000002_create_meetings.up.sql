CREATE TABLE meetings (
    id BIGSERIAL PRIMARY KEY,

    title VARCHAR(200) NOT NULL,
    target_name VARCHAR(200),

    description TEXT,

    scheduled_start_at TIMESTAMPTZ,
    planned_minutes INTEGER NOT NULL,

    decisions TEXT,

    todo TEXT,

    actual_start_at TIMESTAMPTZ,
    actual_end_at TIMESTAMPTZ,

    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',

    current_agenda_id BIGINT,

    created_by BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_meetings_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id),

    CONSTRAINT chk_meeting_status
        CHECK (
            status IN (
                'scheduled',
                'in_progress',
                'completed'
            )
        ),

    CONSTRAINT chk_meeting_planned_minutes
        CHECK (planned_minutes > 0)
);