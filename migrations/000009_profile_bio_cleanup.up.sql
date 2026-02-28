UPDATE profile
SET bio =
$$Profile: The Architect & The Guardian

I am a Software Engineer specializing in Go (Golang) and Cybersecurity. I bridge the gap between high-performance systems and secure system design, bringing experience from INSA and the Boeing STEM Innovation Program.

- Focus: Production-grade Go scaffolding and Zero-Trust APIs
- Security: Alumnus of the Ethiopian Cyber Army (Cohort 3)
- Winner: National Robotics Competition (Autonomous Systems)

Fkremariam Fentahun is a multidisciplinary software and hardware engineer specializing in systems that bridge software, embedded hardware, and cybersecurity. With a proven record from secondary-school science fairs to national robotics competitions and university-level teaching, Fkremariam delivers practical, secure, and user-centered solutions for clients and organizations.

He began his journey in Grade 8 by winning a regional technology competition and continued to excel: a national-level robotics competition winner in Grade 10 (robot forklift project), a contributor to an aviation learning initiative run in collaboration with ThinkYoung and Boeing, and a Grade 12 national science fair participant where he developed an exam control system. He also taught STEM courses at the Gondar STEM Center and later taught Arduino and embedded-systems topics to university students at Addis Ababa Science and Technology University (AAstu).

Fkremariam graduated from INSA and is recognized as part of the 3rd cohort of the Cyber Army of Ethiopia, focusing on cybersecurity. He holds multiple certificates across robotics, cybersecurity, embedded systems, and software development. He has strong programming skills in C++, Python, and Go, hands-on hardware experience (microcontrollers, sensors, robotics), and practical expertise in cybersecurity practices.

He is currently developing a modern, scalable Bingo game and a Duplicate Analyzer tool which demonstrates his strengths in system design, data processing, and user experience. Fkremariam is motivated by building reliable, secure systems that solve real-world problems and enjoys mentoring students and collaborating with multidisciplinary teams.
$$
WHERE id = (SELECT id FROM profile ORDER BY created_at ASC LIMIT 1)
  AND (
    bio LIKE '%ðŸ%'
    OR bio LIKE '%ï¸%'
    OR bio LIKE '%â€%'
    OR bio LIKE '%ð%'
  );

