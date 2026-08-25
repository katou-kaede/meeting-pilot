CREATE TABLE meeting_members (
    meeting_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    role VARCHAR(20) NOT NULL DEFAULT 'viewer',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (meeting_id, user_id),

    CONSTRAINT fk_meeting_members_meeting
        FOREIGN KEY (meeting_id)
        REFERENCES meetings(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_meeting_members_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT chk_meeting_member_role
        CHECK (role IN ('editor', 'viewer'))
);

CREATE INDEX idx_meeting_members_user_id
ON meeting_members(user_id);

CREATE UNIQUE INDEX uq_meeting_members_single_editor
ON meeting_members (meeting_id)
WHERE role = 'editor';