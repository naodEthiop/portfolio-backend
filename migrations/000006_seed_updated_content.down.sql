BEGIN;

-- Best-effort rollback for seed-only inserts. Does not revert profile text changes.
DELETE FROM teaching
WHERE (title = 'STEM Instructor' AND organization = 'Gondar STEM Center')
   OR (title = 'Embedded Systems Teaching Assistant' AND organization LIKE 'Addis Ababa Science%');

DELETE FROM awards
WHERE title IN (
  'Regional Technology Competition Winner',
  'National Robotics Competition Winner',
  'National Science Fair Participant'
);

COMMIT;

