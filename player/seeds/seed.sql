INSERT INTO players (id, created_at, updated_at, age)
VALUES (3, '2022-01-01', '2022-01-01', 1)
ON CONFLICT DO NOTHING;

INSERT INTO players (age)
SELECT ROUND(RANDOM() * 100)
FROM GENERATE_SERIES(1, 200)
ON CONFLICT DO NOTHING;

INSERT INTO food_changes (id, change, player_id)
VALUES (1, 200, 3),
       (2, -100, 3),
       (3, 50, 3)
ON CONFLICT DO NOTHING;

INSERT INTO wood_changes (id, change, player_id)
VALUES (1, 200, 3),
       (2, 100, 3),
       (3, -50, 3)
ON CONFLICT DO NOTHING;

INSERT INTO gold_changes (id, change, player_id)
VALUES (1, 200, 3),
       (2, 300, 3),
       (3, 50, 3)
ON CONFLICT DO NOTHING;

INSERT INTO stone_changes (id, change, player_id)
VALUES (1, 200, 3),
       (2, -10, 3),
       (3, 50, 3)
ON CONFLICT DO NOTHING;
