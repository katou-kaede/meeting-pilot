INSERT INTO users (
    name,
    email,
    password_hash
)
VALUES (
    'Employee1',
    'employee1@example.com',
    -- password123
    '$2a$10$3/oiYZqc56MwmtqY1WtkCeXgx6NOsF8lXSfD80dpQEXxWkIQg/y.y'
);

INSERT INTO meetings (
    title,
    description,
    target_name,
    planned_minutes,
    decisions,
    todo,
    status,
    created_by
)
VALUES (
    '週次開発定例',
    '開発チームの進捗確認',
    '開発チーム',
    60,
    '次回リリースを8月末に決定',
    '・一覧画面のUI改善
    ・詳細画面の実装
    ・テストケース作成',
    'scheduled',
    1
);

INSERT INTO agendas (
    meeting_id,
    title,
    purpose,
    discussion_points,
    questions,
    memo,
    planned_minutes,
    sort_order
)
VALUES
(
    1,
    '進捗確認',
    '各担当の進捗共有',
    '課題を早期に洗い出す',
    'リリースに間に合うか？',
    '認証機能が少し遅れ気味',
    20,
    1
),
(
    1,
    '課題共有',
    'ボトルネック解消',
    '仕様を統一する',
    '画面遷移はこれで良いか？',
    '',
    20,
    2
),
(
    1,
    '今後の予定',
    '次回までの作業整理',
    '期限を明確化する',
    '',
    '',
    20,
    3
);