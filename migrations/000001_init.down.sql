DROP TRIGGER IF EXISTS trg_social_links_updated_at ON social_links;
DROP TRIGGER IF EXISTS trg_skills_updated_at ON skills;
DROP TRIGGER IF EXISTS trg_profile_updated_at ON profile;
DROP TRIGGER IF EXISTS trg_certificates_updated_at ON certificates;
DROP TRIGGER IF EXISTS trg_projects_updated_at ON projects;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

DROP FUNCTION IF EXISTS set_updated_at;

DROP TABLE IF EXISTS social_links;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS profile;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
