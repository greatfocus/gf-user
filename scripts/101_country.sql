CREATE TABLE IF NOT EXISTS country (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	code VARCHAR(5) NOT NULL,
	prefix SMALLINT NOT NULL,
	UNIQUE(name),
	UNIQUE(code),
	UNIQUE(prefix)
);


INSERT INTO country (name, code, prefix)
VALUES
	('Rwanda', 'RW', 250),
	('Kenya', 'KE', 254),	
	('Tanzania', 'TZ', 255),
	('Uganda', 'UG', 256)	
ON CONFLICT (prefix) 
DO NOTHING;