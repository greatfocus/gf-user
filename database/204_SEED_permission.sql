DO $$
DECLARE
    usrId INTEGER := (SELECT id FROM users WHERE identifier = 'muthurimixphone@gmail.com');
	register INTEGER := (SELECT id FROM action WHERE name='register');
    otpVerify INTEGER := (SELECT id FROM action WHERE name='otpVerify');
	lgin INTEGER := (SELECT id FROM action WHERE name='login');
	forgetPassword INTEGER := (SELECT id FROM action WHERE name='forgetPassword');
	manageUser INTEGER := (SELECT id FROM action WHERE name='manageUser');
BEGIN 
    -- Seed System User	
    IF NOT EXISTS (SELECT 1 FROM permission WHERE userId=usrId AND actionId=register) THEN
       INSERT INTO permission (userId, actionId) VALUES(usrId, register);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM permission WHERE userId=usrId AND actionId=otpVerify) THEN
       INSERT INTO permission (userId, actionId) VALUES(usrId, otpVerify);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM permission WHERE userId=usrId AND actionId=lgin) THEN
       INSERT INTO permission (userId, actionId) VALUES(usrId, lgin);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM permission WHERE userId=usrId AND actionId=forgetPassword) THEN
       INSERT INTO permission (userId, actionId) VALUES(usrId, forgetPassword);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM permission WHERE userId=usrId AND actionId=manageUser) THEN
       INSERT INTO permission (userId, actionId) VALUES(usrId, manageUser);
    END IF;
END $$