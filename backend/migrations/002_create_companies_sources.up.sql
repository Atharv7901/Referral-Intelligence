CREATE TABLE company_sources (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    company_id BIGINT UNSIGNED NOT NULL,
    source VARCHAR(50) NOT NULL,
    external_id VARCHAR(255),
    source_url VARCHAR(500),
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    CONSTRAINT fk_company_sources_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,

    UNIQUE KEY uk_company_sources_source_external_id (
        source,
        external_id
    ),

    INDEX idx_company_sources_company_id (company_id),
    INDEX idx_company_sources_source (source)
);