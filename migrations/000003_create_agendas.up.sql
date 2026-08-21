CREATE TABLE agendas (
    id BIGSERIAL PRIMARY KEY,

    meeting_id BIGINT NOT NULL,

    title VARCHAR(200) NOT NULL,

    purpose TEXT,

    discussion_points TEXT,

    questions TEXT,

    memo TEXT,

    planned_minutes INTEGER NOT NULL,

    actual_start_at TIMESTAMPTZ,
    actual_end_at TIMESTAMPTZ,

    elapsed_seconds INTEGER NOT NULL DEFAULT 0,

    sort_order INTEGER NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'not_started',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_agendas_meeting
        FOREIGN KEY (meeting_id)
        REFERENCES meetings(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_planned_minutes
        CHECK (planned_minutes > 0),

    CONSTRAINT chk_agenda_status
        CHECK (
            status IN (
                'not_started',
                'in_progress',
                'completed'
            )
        ),

    CONSTRAINT uq_agendas_sort_order
        UNIQUE (
            meeting_id,
            sort_order
        )
);

CREATE INDEX idx_agendas_meeting_sort
ON agendas(meeting_id, sort_order);