# users
- email
    - UNIQUE
    - ログインIDとして利用
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,

    password_hash VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_users_email UNIQUE (email)
);
```


# meetings
```sql
CREATE TABLE meetings (
    id BIGSERIAL PRIMARY KEY,

    title VARCHAR(200) NOT NULL,
    description TEXT,
    target_name VARCHAR(200),

    scheduled_start_at TIMESTAMPTZ,
    planned_minutes INTEGER NOT NULL,

    actual_start_at TIMESTAMPTZ,
    actual_end_at TIMESTAMPTZ,

    status VARCHAR(20) NOT NULL DEFAULT 'draft',

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
        )

    CONSTRAINT chk_meeting_planned_minutes
        CHECK (planned_minutes > 0)
);
```

# agendas
```sql
CREATE TABLE agendas (
    id BIGSERIAL PRIMARY KEY,

    meeting_id BIGINT NOT NULL,

    title VARCHAR(200) NOT NULL,

    purpose TEXT,

    discussion_points TEXT,

    key_messages TEXT,

    questions TEXT,

    decisions TEXT,

    memo TEXT,

    actual_start_at TIMESTAMPTZ,
    planned_minutes INTEGER NOT NULL,

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
```

# インデックス
```sql
CREATE INDEX idx_meetings_created_by
    ON meetings(created_by);

CREATE INDEX idx_agendas_meeting_id
    ON agendas(meeting_id);

CREATE INDEX idx_agendas_meeting_sort
    ON agendas(meeting_id, sort_order);
```

# ER図
```
users
  1
  │
  │
  N
meetings
  1
  │
  │
  N
agendas
```