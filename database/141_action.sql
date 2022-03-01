CREATE TABLE IF NOT EXISTS action (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(30) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);

INSERT INTO action (name, description)
VALUES
	('register', 'Action to allow register user'),
	('otpVerify', 'Action to allow verify OTP'),
	('login', 'Action to allow Login user'),
	('forgetPassword', 'Action to allow change of password'),
	('manageUser', 'Action to allow create, fetch, update and remove user')
ON CONFLICT (name) 
DO NOTHING;