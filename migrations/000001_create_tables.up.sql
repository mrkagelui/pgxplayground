CREATE TABLE IF NOT EXISTS players
(
    id         BIGSERIAL
        CONSTRAINT players_pk
            PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    age        INT                       NOT NULL
);

CREATE TABLE IF NOT EXISTS food_changes
(
    id         BIGSERIAL
        CONSTRAINT food_changes_pk
            PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    change     INT                       NOT NULL,
    player_id  INT                       NOT NULL
);

CREATE TABLE IF NOT EXISTS wood_changes
(
    id         BIGSERIAL
        CONSTRAINT wood_changes_pk
            PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    change     INT                       NOT NULL,
    player_id  INT                       NOT NULL
);

CREATE TABLE IF NOT EXISTS gold_changes
(
    id         BIGSERIAL
        CONSTRAINT gold_changes_pk
            PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    change     INT                       NOT NULL,
    player_id  INT                       NOT NULL
);

CREATE TABLE IF NOT EXISTS stone_changes
(
    id         BIGSERIAL
        CONSTRAINT stone_changes_pk
            PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    change     INT                       NOT NULL,
    player_id  INT                       NOT NULL
);
