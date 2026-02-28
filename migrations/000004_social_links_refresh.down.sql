BEGIN;

DELETE FROM social_links WHERE platform IN (
  'Telegram Channel',
  'Upwork',
  'Twitter',
  'Reddit',
  'Dev.to',
  'YouTube'
);

COMMIT;

