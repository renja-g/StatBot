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

-- name: GetUserStatusTimeline :many
WITH params AS (
    -- Parameters: $1=guild_id, $2=user_id, $3=target_date
    SELECT
        $1::BIGINT AS p_guild_id,
        $2::BIGINT AS p_user_id,
        $3::DATE AS p_target_date
),
-- Calculate the precise time boundaries for the target date
time_boundaries AS (
    SELECT
        p_target_date::TIMESTAMP AS start_of_day,
        (p_target_date + INTERVAL '1 day')::TIMESTAMP AS end_of_day,
        p_target_date = CURRENT_DATE AS is_today,
        -- Determine the final end point for the query range
        CASE
            WHEN p_target_date = CURRENT_DATE THEN NOW() -- If today, range extends to the current time
            ELSE (p_target_date + INTERVAL '1 day')::TIMESTAMP -- If past day, range extends to end of that day
        END AS end_boundary
    FROM params
),
-- Find the last known status strictly BEFORE the target day began.
-- This defines the status at the very beginning of the target day.
last_status_before AS (
    SELECT
        -- This status is effective from the start of the target day
        (SELECT start_of_day FROM time_boundaries) AS timestamp,
        -- Prioritize status: desktop > mobile > web. If all are NULL, consider 'offline'.
        COALESCE(t.client_status_desktop, t.client_status_mobile, t.client_status_web, 'offline'::discord_status) AS effective_status
    FROM presence_updates t, params p, time_boundaries tb
    WHERE
        t.user_id = p.p_user_id
        AND t.guild_id = p.p_guild_id
        AND t.timestamp < tb.start_of_day
    ORDER BY t.timestamp DESC -- Get the most recent one before the day started
    LIMIT 1
),
-- Find all status updates that occurred ON the target date.
status_updates_today AS (
    SELECT
        t.timestamp,
        -- Prioritize status: desktop > mobile > web. If all are NULL, consider 'offline'.
        COALESCE(t.client_status_desktop, t.client_status_mobile, t.client_status_web, 'offline'::discord_status) AS effective_status
    FROM presence_updates t, params p, time_boundaries tb
    WHERE
        t.user_id = p.p_user_id
        AND t.guild_id = p.p_guild_id
        -- Timestamp must be within the target day: [start_of_day, end_of_day)
        AND t.timestamp >= tb.start_of_day
        AND t.timestamp < tb.end_of_day
),
-- Combine the initial status (from before the day) and the updates during the day.
-- Also, remove consecutive duplicate statuses, as they don't represent a change.
combined_updates AS (
    SELECT
        timestamp,
        effective_status
    FROM (
        SELECT
            timestamp,
            effective_status,
            -- Look at the previous status in the time-ordered sequence
            LAG(effective_status) OVER (ORDER BY timestamp) as prev_status
        FROM (
            -- Union the status from before the day with statuses from the day
            SELECT timestamp, effective_status FROM last_status_before
            UNION ALL
            SELECT timestamp, effective_status FROM status_updates_today
        ) AS u
    ) AS with_prev_status
    -- Only keep rows where the status is different from the previous one
    WHERE effective_status IS DISTINCT FROM prev_status
),
-- Calculate the end time for each status interval using LEAD.
-- The end time of one interval is the start time of the next.
status_intervals AS (
    SELECT
        timestamp AS effective_start_time, -- When this status began
        effective_status AS status,
        -- Find the timestamp of the next status change
        LEAD(timestamp) OVER (ORDER BY timestamp) AS next_change_time
    FROM combined_updates
)
-- Final selection: Filter intervals to the target day range and adjust start/end times.
SELECT
    -- The interval's start time within the day is the later of its effective start or the day's start.
    GREATEST(si.effective_start_time, tb.start_of_day) AS start_time,
    -- The interval's end time is the earlier of the next change time (or the ultimate boundary)
    -- and the ultimate boundary (end_of_day or now()).
    LEAST(
        -- If no next change time, the status lasts until the end_boundary.
        COALESCE(si.next_change_time, tb.end_boundary),
        -- Ensure the calculated end time doesn't go past the end_boundary.
        tb.end_boundary
     ) AS end_time,
    si.status
FROM status_intervals si
CROSS JOIN time_boundaries tb -- Makes boundary values available in each row
WHERE
    -- Keep intervals that overlap with the target period [start_of_day, end_boundary).
    -- Condition 1: Interval must start before the target period ends.
    si.effective_start_time < tb.end_boundary
    -- Condition 2: Interval must end after the target period starts.
    AND COALESCE(si.next_change_time, tb.end_boundary) > tb.start_of_day
ORDER BY start_time; -- Ensure chronological order
