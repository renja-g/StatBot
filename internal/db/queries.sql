-- name: CreatePresenceUpdate :exec
INSERT INTO presence_updates
    (
        timestamp,
        user_id,
        guild_id,
        client_status_desktop,
        client_status_mobile,
        client_status_web,
        activities
    )
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetStatusChangesForDay :many
WITH ranked_updates AS (
    SELECT 
        timestamp,
        user_id,
        COALESCE(client_status_desktop, client_status_mobile, client_status_web) AS status,
        LAG(COALESCE(client_status_desktop, client_status_mobile, client_status_web)) 
            OVER (PARTITION BY user_id ORDER BY timestamp) AS prev_status,
        LEAD(timestamp) 
            OVER (PARTITION BY user_id ORDER BY timestamp) AS next_timestamp
    FROM 
        presence_updates
    WHERE 
        user_id = $1
        AND timestamp >= $2::date
        AND timestamp < ($2::date + interval '1 day')
        AND guild_id = $3
)
SELECT 
    timestamp AS start_time,
    COALESCE(next_timestamp, ($2::date + interval '1 day')) AS end_time,
    status
FROM 
    ranked_updates
WHERE 
    status IS NOT NULL
    AND (prev_status IS NULL OR status != prev_status)
ORDER BY 
    start_time;
