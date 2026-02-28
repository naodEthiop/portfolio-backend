CREATE TABLE IF NOT EXISTS teaching (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(180) NOT NULL,
    organization VARCHAR(180) NOT NULL,
    location VARCHAR(160),
    start_date DATE,
    end_date DATE,
    description TEXT,
    link_url VARCHAR(255),
    sort_order INT NOT NULL DEFAULT 0,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_teaching_visible_order
    ON teaching (visible DESC, sort_order, start_date DESC, created_at DESC);

CREATE TABLE IF NOT EXISTS awards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(180) NOT NULL,
    issuer VARCHAR(180) NOT NULL,
    award_date DATE,
    description TEXT,
    link_url VARCHAR(255),
    sort_order INT NOT NULL DEFAULT 0,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_awards_visible_order
    ON awards (visible DESC, sort_order, award_date DESC, created_at DESC);

DROP TRIGGER IF EXISTS trg_teaching_updated_at ON teaching;
CREATE TRIGGER trg_teaching_updated_at
BEFORE UPDATE ON teaching
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_awards_updated_at ON awards;
CREATE TRIGGER trg_awards_updated_at
BEFORE UPDATE ON awards
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

