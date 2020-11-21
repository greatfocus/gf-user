CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	type VARCHAR(20) NOT NULL,
	email VARCHAR(100) NOT NULL,	
	password VARCHAR(100) NOT NULL,
	failedAttempts INTEGER default (0),
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expiredDate TIMESTAMP NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(email)	
);

DO $$ 
DECLARE
	country INTEGER := (select id from country where code='KE');

BEGIN
	INSERT INTO users (type, email, password, expiredDate, status, enabled)
	VALUES
		('admin', 'muthurimixphone@gmail.com', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', '2025-06-01 08:22:17.460493', 'USER.APPROVED', true) 
	ON CONFLICT
	DO NOTHING;
END $$;