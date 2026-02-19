SELECT
    id,
    title,
    status,
    created_at,
    updated_at
FROM
    tasks
WHERE
    id = $1
;