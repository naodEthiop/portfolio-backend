BEGIN;

-- Profile: updated positioning + hero banner + CTAs
UPDATE profile
SET
  handle = 'NaodEthiop',
  headline = 'Software Engineer — Go (Golang) & Cybersecurity',
  summary = 'I build production-grade Go systems and defensive-first APIs with clean architecture and secure defaults. INSA alumnus • Cyber Army (Cohort 3) • National robotics winner.',
  banner_url = 'https://capsule-render.vercel.app/api?type=waving&color=0d1117&height=250&section=header&text=NaodEthiop&fontSize=70&fontColor=00ADD8&fontAlignY=35&animation=twinkling',
  cta_primary = 'Hire me for embedded or cybersecurity work',
  cta_secondary = 'Request a project quote',
  cta_tertiary = 'Schedule a 15-minute call'
WHERE id = (SELECT id FROM profile ORDER BY created_at ASC LIMIT 1);

-- Skills (ensure categories and updated names exist; avoid duplicates)
UPDATE skills SET name = 'Go (Golang)' WHERE category = 'Software' AND name = 'Go';

INSERT INTO skills (category, name, sort_order)
SELECT 'Software', v.name, v.sort_order
FROM (VALUES
  ('C++', 1),
  ('Python', 2),
  ('Go (Golang)', 3),
  ('Node.js', 4),
  ('React', 5),
  ('Websockets', 6)
) AS v(name, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM skills s WHERE s.category = 'Software' AND s.name = v.name);

INSERT INTO skills (category, name, sort_order)
SELECT 'Embedded/Hardware', v.name, v.sort_order
FROM (VALUES
  ('Arduino', 1),
  ('AVR', 2),
  ('STM32', 3),
  ('Microcontrollers', 4),
  ('Sensors', 5),
  ('Motors', 6)
) AS v(name, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM skills s WHERE s.category = 'Embedded/Hardware' AND s.name = v.name);

INSERT INTO skills (category, name, sort_order)
SELECT 'Cybersecurity', v.name, v.sort_order
FROM (VALUES
  ('Security Operations', 1),
  ('Secure System Design', 2),
  ('Penetration Testing', 3),
  ('Risk Assessment', 4)
) AS v(name, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM skills s WHERE s.category = 'Cybersecurity' AND s.name = v.name);

INSERT INTO skills (category, name, sort_order)
SELECT 'Tools & Platforms', v.name, v.sort_order
FROM (VALUES
  ('Git', 1),
  ('Docker', 2),
  ('Linux', 3),
  ('PostgreSQL', 4)
) AS v(name, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM skills s WHERE s.category = 'Tools & Platforms' AND s.name = v.name);

-- Teaching
INSERT INTO teaching (title, organization, start_date, end_date, description, sort_order, visible)
SELECT 'STEM Instructor', 'Gondar STEM Center', DATE '2021-01-01', DATE '2022-12-31',
  'Delivered comprehensive STEM courses to high-school students, focusing on hands-on learning and practical applications.',
  1, TRUE
WHERE NOT EXISTS (SELECT 1 FROM teaching WHERE title = 'STEM Instructor' AND organization = 'Gondar STEM Center');

INSERT INTO teaching (title, organization, start_date, end_date, description, sort_order, visible)
SELECT 'Embedded Systems Teaching Assistant', 'Addis Ababa Science and Technology University (AAstu)', DATE '2022-01-01', DATE '2023-12-31',
  'Taught Arduino and embedded-systems topics to university students, including lab sessions and project guidance.',
  2, TRUE
WHERE NOT EXISTS (SELECT 1 FROM teaching WHERE title = 'Embedded Systems Teaching Assistant' AND organization LIKE 'Addis Ababa Science%');

-- Awards
INSERT INTO awards (title, issuer, award_date, description, sort_order, visible)
SELECT 'Regional Technology Competition Winner', 'Competition', DATE '2018-01-01',
  'Grade 8 regional technology competition winner',
  1, TRUE
WHERE NOT EXISTS (SELECT 1 FROM awards WHERE title = 'Regional Technology Competition Winner');

INSERT INTO awards (title, issuer, award_date, description, sort_order, visible)
SELECT 'National Robotics Competition Winner', 'Competition', DATE '2020-01-01',
  'Grade 10 national robotics competition winner with robot forklift project',
  2, TRUE
WHERE NOT EXISTS (SELECT 1 FROM awards WHERE title = 'National Robotics Competition Winner');

INSERT INTO awards (title, issuer, award_date, description, sort_order, visible)
SELECT 'National Science Fair Participant', 'Competition', DATE '2022-01-01',
  'Grade 12 national science fair participant with exam control system',
  3, TRUE
WHERE NOT EXISTS (SELECT 1 FROM awards WHERE title = 'National Science Fair Participant');

-- Certificates (normalize dates to yyyy-01-01)
UPDATE certificates
SET issue_date = DATE '2023-01-01', description = 'Advanced cybersecurity training and operations certification'
WHERE name = 'Cyber Army of Ethiopia - Cohort 3' AND issuer = 'Cyber Army of Ethiopia';

INSERT INTO certificates (name, issuer, issue_date, description)
SELECT 'Cyber Army of Ethiopia - Cohort 3', 'Cyber Army of Ethiopia', DATE '2023-01-01',
  'Advanced cybersecurity training and operations certification'
WHERE NOT EXISTS (SELECT 1 FROM certificates WHERE name = 'Cyber Army of Ethiopia - Cohort 3' AND issuer = 'Cyber Army of Ethiopia');

INSERT INTO certificates (name, issuer, issue_date, description)
SELECT 'ThinkYoung + Boeing Aviation Program', 'ThinkYoung & Boeing', DATE '2022-01-01',
  'Aviation-focused learning program in collaboration with industry leaders'
WHERE NOT EXISTS (SELECT 1 FROM certificates WHERE name = 'ThinkYoung + Boeing Aviation Program' AND issuer = 'ThinkYoung & Boeing');

INSERT INTO certificates (name, issuer, issue_date, description)
SELECT 'National Robotics Competition Winner', 'National Robotics Association', DATE '2020-01-01',
  'First place in national robotics competition with robot forklift project'
WHERE NOT EXISTS (SELECT 1 FROM certificates WHERE name = 'National Robotics Competition Winner' AND issuer = 'National Robotics Association');

INSERT INTO certificates (name, issuer, issue_date, description)
SELECT 'INSA Summer Camp Completion - Cyber Army Cohort 3', 'INSA / Cyber Army of Ethiopia', DATE '2025-01-01',
  'Completed INSA Summer Camp and selected in the 3rd round of the Cyber Army of Ethiopia'
WHERE NOT EXISTS (SELECT 1 FROM certificates WHERE name = 'INSA Summer Camp Completion - Cyber Army Cohort 3' AND issuer LIKE 'INSA%');

COMMIT;

