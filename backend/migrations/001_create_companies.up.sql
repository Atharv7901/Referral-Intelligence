CREATE TABLE companies (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

    name VARCHAR(255) NOT NULL,

    slug VARCHAR(255) NOT NULL,

    website VARCHAR(500),

    linkedin_url VARCHAR(500),

    career_page_url VARCHAR(500),

    company_type VARCHAR(50) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    UNIQUE KEY uk_companies_slug (slug),

    INDEX idx_companies_name (name),

    INDEX idx_companies_active (is_active)
);