INSERT INTO tasks (
    id,
    title,
    status,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);
