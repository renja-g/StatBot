CREATE TYPE discord_status AS ENUM ('online', 'idle', 'dnd', 'offline');

CREATE TABLE presence_updates (
    timestamp TIMESTAMPTZ NOT NULL,

    user_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,


    -- Status per client type (can be NULL if the user isn't logged in on that client)
    client_status_desktop discord_status,
    client_status_mobile discord_status,
    client_status_web discord_status,

    -- Store the full activities array as JSONB
    -- See: https://discord.com/developers/docs/topics/gateway-events#activity-object
    activities JSONB
);

-- 2. Turn the regular table into a TimescaleDB hypertable
-- This partitions the data based on the 'timestamp' column for performance
SELECT create_hypertable('presence_updates', 'timestamp');

-- Index for looking up specific user/guild presences over time
CREATE INDEX idx_presence_user_guild_time ON presence_updates (user_id, guild_id, timestamp DESC);

-- Index for querying data within the activities JSONB column (if needed)
-- CREATE INDEX idx_presence_activities ON presence_updates USING GIN (activities);