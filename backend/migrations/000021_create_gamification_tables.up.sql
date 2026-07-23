CREATE TABLE user_stats (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    xp INTEGER NOT NULL DEFAULT 0,
    total_seconds_active BIGINT NOT NULL DEFAULT 0,
    last_heartbeat_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE achievements (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    tier SMALLINT NOT NULL,
    metric VARCHAR(50) NOT NULL,
    threshold INTEGER NOT NULL,
    xp_reward INTEGER NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_achievements (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id BIGINT NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,

    unlocked_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, achievement_id)
);
