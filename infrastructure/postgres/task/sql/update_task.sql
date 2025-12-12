UPDATE
    tasks 
SET (
    title,
    status,
    created_at,
    updated_at
) = (
    $2,
    $3,
    $4,
    $5
)
WHERE
    id = $1
;