CREATE TABLE IF NOT EXISTS client (
	id BIGSERIAL PRIMARY KEY,
	email VARCHAR(100) NOT NULL,
	clientId VARCHAR(100) NOT NULL,
	secret VARCHAR(100) NOT NULL,
	failedAttempts INTEGER default (0),
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expiredDate TIMESTAMP NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(true),
	UNIQUE(email),
	UNIQUE(clientId),
	UNIQUE(secret)	
);

DO $$ 
BEGIN
	INSERT INTO client (email, clientId, secret, expiredDate)
	VALUES
		('muthurimixphone@gmail.com', '$2a$04$u2w.uoMltUj/eq8rBBPiTew9EEHctusyhV4GE5vf9H2RvgUdlfeE.', '$2a$04$n8VdTg8GItB.q455rjv7XucbQXoT9xHYxza2NS2SVzFPosLd9dayO', '2020-11-21 10:48:53.513464') 
	ON CONFLICT
	DO NOTHING;
END $$;

-- '98590c398a254d2898838e1b17381575', 'ADRtjWLkttBbMQLpMADF'