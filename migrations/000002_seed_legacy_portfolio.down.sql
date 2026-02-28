BEGIN;

DELETE FROM projects
WHERE slug IN (
  'web-recon-tool',
  'bingo-game',
  'duplicate-analyzer',
  'robot-forklift',
  'exam-control-system'
);

DELETE FROM certificates
WHERE name IN (
  'National Science Fair Participation',
  'Robotics Competition 3rd Place',
  'Grade 12 National Exam High Score',
  'Regional Science Fair 3rd Place',
  'Zone Science Fair 1st Place',
  'National Science Fair Recognition',
  'Regional Science Fair Participation',
  'ThinkYoung STEM School Completion',
  '7th National Science Fair Winner',
  'ThinkYoung Alumni Training',
  'Arduino Programming Trainer',
  'EthioCoder Programming Fundamentals',
  'EthioCoder Android Fundamentals',
  'EthioCoder Data Analysis',
  'Gondar STEM School Trainer',
  'INSA Summer Camp Completion - Cyber Army Cohort 3',
  'Cyber Army of Ethiopia - Cohort 3',
  'ThinkYoung + Boeing Aviation Program',
  'National Robotics Competition Winner'
);

DELETE FROM skills
WHERE category IN ('Software', 'Embedded/Hardware', 'Cybersecurity', 'Tools & Platforms');

DELETE FROM social_links
WHERE platform IN ('Email', 'Telegram', 'GitHub', 'LinkedIn');

UPDATE profile
SET
  full_name = 'Portfolio Owner',
  headline = 'Profile is being configured',
  location = 'Unknown',
  summary = 'Set profile details from admin dashboard',
  bio = '',
  avatar_url = '',
  resume_url = ''
WHERE id = (SELECT id FROM profile ORDER BY created_at ASC LIMIT 1);

COMMIT;
