-- +goose Up

ALTER TABLE seed ADD COLUMN img_plant VARCHAR(255);
ALTER TABLE seed ALTER COLUMN icon TYPE VARCHAR(255);

INSERT INTO seed (id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at) VALUES
(1, 'Пшеница', '/assets/seeds/wheat.png', '/assets/plants/wheat', 1, 4, 'common', 0.25, 5, 40, NOW()),
(2, 'Баклажан', '/assets/seeds/aubergine.png', '/assets/plants/aubergine', 1, 8, 'uncommon', 0.5, 12, 150, NOW())
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    img_plant = EXCLUDED.img_plant,
    level_required = EXCLUDED.level_required,
    target_growth = EXCLUDED.target_growth,
    rarity = EXCLUDED.rarity,
    modification = EXCLUDED.modification,
    gold_reward = EXCLUDED.gold_reward,
    xp_reward = EXCLUDED.xp_reward,
    created_at = EXCLUDED.created_at;

SELECT setval('seed_id_seq', (SELECT MAX(id) FROM seed));

-- +goose Down
DELETE FROM seed WHERE id BETWEEN 1 AND 2;

ALTER TABLE seed DROP COLUMN img_plant;
ALTER TABLE seed ALTER COLUMN icon TYPE VARCHAR(100);