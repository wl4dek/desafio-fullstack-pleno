CREATE TABLE IF NOT EXISTS health (
    id BIGSERIAL PRIMARY KEY,
    child_id TEXT NOT NULL UNIQUE,

    vaccinations_up_to_date BOOLEAN NOT NULL DEFAULT FALSE,
    last_consultation DATE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (child_id) REFERENCES children(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS education (
    id BIGSERIAL PRIMARY KEY,
    child_id TEXT NOT NULL UNIQUE,

    school_name TEXT,
    frequency_percent NUMERIC,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (child_id) REFERENCES children(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS social_assistance (
    id BIGSERIAL PRIMARY KEY,
    child_id TEXT NOT NULL UNIQUE,

    cad_unico boolean NOT NULL DEFAULT FALSE,
    active_benefit boolean NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (child_id) REFERENCES children(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS alert_health (
    id BIGSERIAL PRIMARY KEY,
    health_id BIGSERIAL NOT NULL,

    code TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (health_id) REFERENCES health(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS alert_education (
    id BIGSERIAL PRIMARY KEY,
    education_id BIGSERIAL NOT NULL,

    code TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (education_id) REFERENCES education(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS alert_social_assistance (
    id BIGSERIAL PRIMARY KEY,
    social_assistance_id BIGSERIAL NOT NULL,

    code TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (social_assistance_id) REFERENCES social_assistance(id) ON DELETE CASCADE
);