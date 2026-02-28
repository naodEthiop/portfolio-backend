UPDATE profile
SET bio =
  $$### 🧠 Profile: The Architect & The Guardian

I am a **Software Engineer** specializing in **Go (Golang)** and **Cybersecurity**. I bridge the gap between high-performance systems and ironclad security, bringing experience from **INSA** and the **Boeing STEM Innovation Program**.

- 🛠️ **Focus:** Production-grade Go scaffolding & Zero-Trust APIs
- 🔐 **Security:** Alumnus of the Ethiopian Cyber Army (3rd Cohort)
- 🤖 **Winner:** National Robotics Competition (Autonomous Systems)

---

$$ || COALESCE(bio, '')
WHERE id = (SELECT id FROM profile ORDER BY created_at ASC LIMIT 1);

