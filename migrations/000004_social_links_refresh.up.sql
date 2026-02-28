BEGIN;

UPDATE social_links
SET url = 'https://t.me/Naod2i', sort_order = 2, visible = TRUE
WHERE platform = 'Telegram';

UPDATE social_links
SET url = 'https://www.linkedin.com/in/fkremariam-fentahun-b9a902390', sort_order = 4, visible = TRUE
WHERE platform = 'LinkedIn';

UPDATE social_links
SET url = 'https://github.com/naodEthiop', sort_order = 5, visible = TRUE
WHERE platform = 'GitHub';

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'Telegram Channel', 'https://t.me/naodbuilds', 3, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'Telegram Channel');

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'Upwork', 'https://www.upwork.com/freelancers/~01531376e3abd2dc50?mp_source=share', 6, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'Upwork');

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'Twitter', 'https://twitter.com/naodEthiop', 7, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'Twitter');

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'Reddit', 'https://www.reddit.com/user/naodEthiop', 8, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'Reddit');

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'Dev.to', 'https://dev.to/naodEthiop', 9, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'Dev.to');

INSERT INTO social_links (platform, url, sort_order, visible)
SELECT 'YouTube', 'https://www.youtube.com/@NaodBuilds', 10, TRUE
WHERE NOT EXISTS (SELECT 1 FROM social_links WHERE platform = 'YouTube');

COMMIT;

