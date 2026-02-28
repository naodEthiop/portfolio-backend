BEGIN;

-- Ensure Bingo repo is linkable
UPDATE projects
SET
  repo_url = 'https://github.com/naodEthiop/Multiplayer-Bingo-Game',
  demo_url = 'https://github.com/naodEthiop/Multiplayer-Bingo-Game'
WHERE slug = 'bingo-game';

-- Add Lalibela CLI as a manual project to guarantee it shows (GitHub data will merge when available)
INSERT INTO projects (
  title, slug, short_description, description, status,
  tech_stack, achievements, demo_url, repo_url, featured
)
SELECT
  'Lalibela CLI',
  'lalibela-cli',
  'Go scaffolder engine for production-grade services.',
  'A Go-based scaffolding CLI that generates clean, production-oriented project structures and API foundations.',
  'complete',
  ARRAY['Go (Golang)','CLI','Scaffolding','Templates'],
  ARRAY['Production-grade scaffolding','Fast project bootstrap'],
  'https://github.com/naodEthiop/lalibela-cli',
  'https://github.com/naodEthiop/lalibela-cli',
  TRUE
WHERE NOT EXISTS (SELECT 1 FROM projects WHERE slug = 'lalibela-cli');

COMMIT;

