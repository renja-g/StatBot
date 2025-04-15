// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPresenceUpdate = `-- name: CreatePresenceUpdate :exec
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
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreatePresenceUpdateParams struct {
	Timestamp           pgtype.Timestamptz
	UserID              int64
	GuildID             int64
	ClientStatusDesktop NullDiscordStatus
	ClientStatusMobile  NullDiscordStatus
	ClientStatusWeb     NullDiscordStatus
	Activities          []byte
}

func (q *Queries) CreatePresenceUpdate(ctx context.Context, arg CreatePresenceUpdateParams) error {
	_, err := q.db.Exec(ctx, createPresenceUpdate,
		arg.Timestamp,
		arg.UserID,
		arg.GuildID,
		arg.ClientStatusDesktop,
		arg.ClientStatusMobile,
		arg.ClientStatusWeb,
		arg.Activities,
	)
	return err
}

const getStatusChangesForDay = `-- name: GetStatusChangesForDay :many
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
    start_time
`

type GetStatusChangesForDayParams struct {
	UserID  int64
	Column2 pgtype.Date
	GuildID int64
}

type GetStatusChangesForDayRow struct {
	StartTime pgtype.Timestamptz
	EndTime   interface{}
	Status    NullDiscordStatus
}

func (q *Queries) GetStatusChangesForDay(ctx context.Context, arg GetStatusChangesForDayParams) ([]GetStatusChangesForDayRow, error) {
	rows, err := q.db.Query(ctx, getStatusChangesForDay, arg.UserID, arg.Column2, arg.GuildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStatusChangesForDayRow
	for rows.Next() {
		var i GetStatusChangesForDayRow
		if err := rows.Scan(&i.StartTime, &i.EndTime, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
