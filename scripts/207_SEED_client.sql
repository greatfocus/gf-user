DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM client) THEN
        INSERT INTO client (email, clientId, secret, expiredDate)
        VALUES
            ('muthurimixphone@gmail.com', '$2a$04$u2w.uoMltUj/eq8rBBPiTew9EEHctusyhV4GE5vf9H2RvgUdlfeE.', '$2a$04$n8VdTg8GItB.q455rjv7XucbQXoT9xHYxza2NS2SVzFPosLd9dayO', '2020-11-21 10:48:53.513464') 
        ON CONFLICT
        DO NOTHING;

        -- '98590c398a254d2898838e1b17381575', 'ADRtjWLkttBbMQLpMADF'
    END IF;
END $$