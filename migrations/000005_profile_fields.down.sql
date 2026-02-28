ALTER TABLE profile
    DROP COLUMN IF EXISTS cta_tertiary,
    DROP COLUMN IF EXISTS cta_secondary,
    DROP COLUMN IF EXISTS cta_primary,
    DROP COLUMN IF EXISTS banner_url,
    DROP COLUMN IF EXISTS handle;

