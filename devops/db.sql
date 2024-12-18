CREATE TABLE vulnerabilities (
    plugin_id VARCHAR(255) PRIMARY KEY,
    asset_name VARCHAR(255) NOT NULL,
    plugin_name VARCHAR(255) NOT NULL,
    severity VARCHAR(50),
    cvss3_base_score VARCHAR(50),
    vulnerability_priority_rating VARCHAR(50),
    plugin_output TEXT,
    synopsis TEXT,
    description TEXT,
    solution TEXT,
    vulnerability_first_discovered_date DATE,
    vulnerability_last_observed_date DATE,
    plugin_publication_date DATE,
    patch_publication_date DATE,
    exploitability VARCHAR(50)
);