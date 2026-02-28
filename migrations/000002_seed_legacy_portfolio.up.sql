BEGIN;

-- Profile from legacy profile.json
INSERT INTO profile (full_name, headline, location, summary, bio, avatar_url, resume_url)
SELECT
  'Fkremariam Fentahun',
  'Software, Robotics & Cybersecurity Engineer',
  'Addis Ababa, Ethiopia',
  'Multidisciplinary engineer with competition-winning robotics experience, cybersecurity training, and university-level teaching. Strong in C++, Python, Go, embedded systems, and secure system design.',
  $$Fkremariam Fentahun is a multidisciplinary software and hardware engineer specializing in systems that bridge software, embedded hardware, and cybersecurity. With a proven record from secondary-school science fairs to national robotics competitions and university-level teaching, Fkremariam delivers practical, secure, and user-centered solutions for clients and organizations.

He began his journey in Grade 8 by winning a regional technology competition and continued to excel: a national-level robotics competition winner in Grade 10 (robot forklift project), a contributor to an aviation learning initiative run in collaboration with ThinkYoung and Boeing, and a Grade 12 national science fair participant where he developed an exam control system. He also taught STEM courses at the Gondar STEM Center and later taught Arduino and embedded-systems topics to university students at Addis Ababa Science and Technology University (AAstu).

Fkremariam graduated from INSA and is recognized as part of the 3rd cohort of the Cyber Army of Ethiopia, focusing on cybersecurity. He holds multiple certificates across robotics, cybersecurity, embedded systems, and software development. He has strong programming skills in C++, Python, and Go, hands-on hardware experience (microcontrollers, sensors, robotics), and practical expertise in cybersecurity practices.

He is currently developing a modern, scalable Bingo game and a Duplicate Analyzer tool which demonstrates his strengths in system design, data processing, and user experience. Fkremariam is motivated by building reliable, secure systems that solve real-world problems and enjoys mentoring students and collaborating with multidisciplinary teams.$$,
  '/uploads/profile/profilePhoto.avif',
  ''
WHERE NOT EXISTS (SELECT 1 FROM profile);

UPDATE profile
SET
  full_name = 'Fkremariam Fentahun',
  headline = 'Software, Robotics & Cybersecurity Engineer',
  location = 'Addis Ababa, Ethiopia',
  summary = 'Multidisciplinary engineer with competition-winning robotics experience, cybersecurity training, and university-level teaching. Strong in C++, Python, Go, embedded systems, and secure system design.',
  bio = $$Fkremariam Fentahun is a multidisciplinary software and hardware engineer specializing in systems that bridge software, embedded hardware, and cybersecurity. With a proven record from secondary-school science fairs to national robotics competitions and university-level teaching, Fkremariam delivers practical, secure, and user-centered solutions for clients and organizations.

He began his journey in Grade 8 by winning a regional technology competition and continued to excel: a national-level robotics competition winner in Grade 10 (robot forklift project), a contributor to an aviation learning initiative run in collaboration with ThinkYoung and Boeing, and a Grade 12 national science fair participant where he developed an exam control system. He also taught STEM courses at the Gondar STEM Center and later taught Arduino and embedded-systems topics to university students at Addis Ababa Science and Technology University (AAstu).

Fkremariam graduated from INSA and is recognized as part of the 3rd cohort of the Cyber Army of Ethiopia, focusing on cybersecurity. He holds multiple certificates across robotics, cybersecurity, embedded systems, and software development. He has strong programming skills in C++, Python, and Go, hands-on hardware experience (microcontrollers, sensors, robotics), and practical expertise in cybersecurity practices.

He is currently developing a modern, scalable Bingo game and a Duplicate Analyzer tool which demonstrates his strengths in system design, data processing, and user experience. Fkremariam is motivated by building reliable, secure systems that solve real-world problems and enjoys mentoring students and collaborating with multidisciplinary teams.$$,
  avatar_url = '/uploads/profile/profilePhoto.avif',
  resume_url = ''
WHERE id = (SELECT id FROM profile ORDER BY created_at ASC LIMIT 1);

-- Replace seed content with legacy projects/skills/social links/certificates
DELETE FROM projects;
DELETE FROM certificates;
DELETE FROM skills;
DELETE FROM social_links;

INSERT INTO projects (
  title, slug, short_description, description, status,
  tech_stack, achievements, demo_url, repo_url, image_url, featured
)
VALUES
(
  'Web Recon Tool',
  'web-recon-tool',
  'Professional reconnaissance framework for cybersecurity researchers and ethical hackers.',
  $$A comprehensive reconnaissance tool designed for penetration testers and security researchers. Features port scanning, vulnerability assessment, subdomain enumeration, directory discovery, and Shodan API integration. Built with Go and Docker for professional security operations.$$,
  'complete',
  ARRAY['Go','Docker','Nmap','Shodan API','CLI','Security'],
  ARRAY['24 GitHub stars','Professional security tool','Docker containerization','Shodan API integration','Multiple reconnaissance modules'],
  'https://github.com/naodEthiop/web_recon_tool',
  'https://github.com/naodEthiop/web_recon_tool',
  NULL,
  TRUE
),
(
  'Bingo Game',
  'bingo-game',
  'Modern multiplayer bingo game with realtime features and gamification.',
  $$A scalable multiplayer bingo game built with real-time WebSocket communication, featuring user authentication, payment integration, and gamification elements. The system handles concurrent users with efficient data processing and provides a smooth gaming experience across devices.$$,
  'in_progress',
  ARRAY['Node.js','Websockets','React','PostgreSQL','Docker'],
  ARRAY['Real-time multiplayer functionality','Scalable architecture','Payment integration','Cross-platform compatibility'],
  NULL,
  NULL,
  NULL,
  FALSE
),
(
  'Duplicate Analyzer',
  'duplicate-analyzer',
  'Tool to detect and remove duplicate records from large datasets.',
  $$An intelligent data processing tool that uses vector search and embedding techniques to identify and remove duplicate records from large datasets. Features include fuzzy matching, performance optimization, and automated data cleanup workflows.$$,
  'in_progress',
  ARRAY['Python','Go','Vector Search','Embeddings','Data Processing'],
  ARRAY['High-performance data processing','Fuzzy matching algorithms','Automated ETL workflows','Scalable architecture'],
  NULL,
  NULL,
  NULL,
  FALSE
),
(
  'Robot Forklift',
  'robot-forklift',
  'Competition-winning robotics project built during secondary school competitions.',
  $$A fully autonomous robot forklift designed and built for national robotics competition. The robot features precise navigation, object detection, and lifting mechanisms, demonstrating advanced robotics engineering and programming skills.$$,
  'complete',
  ARRAY['Arduino','Motors','Sensors','C++','Robotics'],
  ARRAY['National competition winner','Autonomous navigation','Object detection and manipulation','Precision control systems'],
  NULL,
  NULL,
  NULL,
  TRUE
),
(
  'Exam Control System',
  'exam-control-system',
  'Digital examination management system for educational institutions.',
  $$A comprehensive exam control system developed for the Grade 12 national science fair. Features include secure exam delivery, real-time monitoring, and automated grading capabilities.$$,
  'complete',
  ARRAY['Python','Web Technologies','Database Design','Security'],
  ARRAY['National science fair participant','Secure exam delivery','Real-time monitoring','Automated grading'],
  NULL,
  NULL,
  NULL,
  FALSE
);

INSERT INTO skills (category, name, sort_order) VALUES
('Software','C++',1),
('Software','Python',2),
('Software','Golang',3),
('Software','Node.js',4),
('Software','React',5),
('Software','Websockets',6),
('Embedded/Hardware','Arduino',1),
('Embedded/Hardware','AVR',2),
('Embedded/Hardware','STM32',3),
('Embedded/Hardware','Microcontrollers',4),
('Embedded/Hardware','Sensors',5),
('Embedded/Hardware','Motors',6),
('Cybersecurity','Security Operations',1),
('Cybersecurity','Secure System Design',2),
('Cybersecurity','Penetration Testing',3),
('Cybersecurity','Risk Assessment',4),
('Tools & Platforms','Git',1),
('Tools & Platforms','Docker',2),
('Tools & Platforms','Linux',3),
('Tools & Platforms','PostgreSQL',4);

INSERT INTO social_links (platform, url, sort_order, visible) VALUES
('Email','mailto:fkremariamfentahun66@gmail.com',1,TRUE),
('Telegram','https://t.me/naod2i',2,TRUE),
('GitHub','https://github.com/naodEthiop',3,TRUE),
('LinkedIn','https://linkedin.com/in/fkremariam-fentahun',4,TRUE);

INSERT INTO certificates (name, issuer, issue_date, description, credential_url, image_url) VALUES
('National Science Fair Participation','Ministry of Education',DATE '2022-01-01','Type: science. Participation certificate for Grade 12 National Science Fair',NULL,'/uploads/certificates/cert1.avif'),
('Robotics Competition 3rd Place','National Robotics Association',DATE '2020-01-01','Type: robotics. 3rd rank certificate of excellence in robotics competition',NULL,'/uploads/certificates/cert2.avif'),
('Grade 12 National Exam High Score','Ministry of Education',DATE '2022-01-01','Type: academic. Certificate for scoring high marks in Grade 12 National Exam',NULL,'/uploads/certificates/cert3.avif'),
('Regional Science Fair 3rd Place','Regional Education Bureau',DATE '2021-01-01','Type: science. 3rd rank in Regional Science Fair competition',NULL,'/uploads/certificates/cert4.avif'),
('Zone Science Fair 1st Place','Zone Education Office',DATE '2021-01-01','Type: science. 1st rank winner in Zone Science Fair',NULL,'/uploads/certificates/cert5.avif'),
('National Science Fair Recognition','Ministry of Education',DATE '2022-01-01','Type: science. Recognition certificate for National Science Fair participation',NULL,'/uploads/certificates/cert6.avif'),
('Regional Science Fair Participation','Regional Education Bureau',DATE '2021-01-01','Type: science. Recognition of participation in Regional Science Fair',NULL,'/uploads/certificates/cert7.avif'),
('Regional Science Fair Participation','Regional Education Bureau',DATE '2020-01-01','Type: science. Recognition of participation in Regional Science Fair',NULL,'/uploads/certificates/cert8.avif'),
('ThinkYoung STEM School Completion','ThinkYoung',DATE '2021-01-01','Type: education. Certification of completion from ThinkYoung STEM School',NULL,'/uploads/certificates/cert9.avif'),
('Regional Science Fair 3rd Place','Regional Education Bureau',DATE '2020-01-01','Type: science. Recognition of participation and 3rd rank in Regional Science Fair',NULL,'/uploads/certificates/cert10.avif'),
('7th National Science Fair Winner','Israel Embassy',DATE '2022-01-01','Type: science. Certificate from Israel Embassy for winning 7th National Science Fair',NULL,'/uploads/certificates/cert11.avif'),
('National Science Fair Participation','Ministry of Education',DATE '2021-01-01','Type: science. Recognition of participation in National Science Fair',NULL,'/uploads/certificates/cert12.avif'),
('ThinkYoung Alumni Training','ThinkYoung',DATE '2022-01-01','Type: education. Certification of completion from ThinkYoung STEM School Alumni Training Programme',NULL,'/uploads/certificates/cert13.avif'),
('Arduino Programming Trainer','Thea Ruino',DATE '2023-01-01','Type: teaching. Certificate as trainer in Arduino Programming Training',NULL,'/uploads/certificates/cert15.avif'),
('EthioCoder Programming Fundamentals','EthioCoder',DATE '2023-01-01','Type: programming. Programming Fundamentals certification from EthioCoder',NULL,'/uploads/certificates/cert16.avif'),
('EthioCoder Android Fundamentals','EthioCoder',DATE '2023-01-01','Type: programming. Android Development Fundamentals certification from EthioCoder',NULL,'/uploads/certificates/cert17.avif'),
('EthioCoder Data Analysis','EthioCoder',DATE '2023-01-01','Type: data. Data Analysis certification from EthioCoder',NULL,'/uploads/certificates/cert18.avif'),
('Gondar STEM School Trainer','Gondar STEM Center',NULL,'Type: teaching. Certification as trainer at Gondar STEM School (Coming Soon)',NULL,'/uploads/certificates/placeholder.avif'),
('INSA Summer Camp Completion - Cyber Army Cohort 3','INSA / Cyber Army of Ethiopia',DATE '2025-01-01','Type: education. Completed INSA Summer Camp and selected in 3rd round of the Cyber Army of Ethiopia',NULL,'/uploads/certificates/cert19.avif'),
('Cyber Army of Ethiopia - Cohort 3','Cyber Army of Ethiopia',DATE '2023-01-01','Advanced cybersecurity training and operations certification',NULL,'/uploads/certificates/placeholder.avif'),
('ThinkYoung + Boeing Aviation Program','ThinkYoung & Boeing',DATE '2022-01-01','Aviation-focused learning program in collaboration with industry leaders',NULL,'/uploads/certificates/placeholder.avif'),
('National Robotics Competition Winner','National Robotics Association',DATE '2020-01-01','First place in national robotics competition with robot forklift project',NULL,'/uploads/certificates/placeholder.avif');

COMMIT;
