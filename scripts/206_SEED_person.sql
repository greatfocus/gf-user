DO $$
DECLARE
	country INTEGER := (select id from country where code='KE');
	usrId INTEGER := (select id from users where email='muthurimixphone@gmail.com');
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM person WHERE userId=usrId) THEN
        INSERT INTO person (userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber, createdBy, updatedBy)
        VALUES
            (usrId, country, 'John', 'Peter', 'Mucunga', '0780904371', '27496388', usrId, usrId) 
        ON CONFLICT
        DO NOTHING;
    END IF;
END $$