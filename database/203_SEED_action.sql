DO $$
BEGIN 
    INSERT INTO action (name, description)
    VALUES
        ('register', 'Action to allow register user'),
        ('otpVerify', 'Action to allow verify OTP'),
        ('login', 'Action to allow Login user'),
        ('forgetPassword', 'Action to allow change of password'),
        ('manageUser', 'Action to allow create, fetch, update and remove user')
    ON CONFLICT (name) 
    DO NOTHING;
END $$