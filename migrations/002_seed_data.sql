INSERT INTO News (title, body, datetime) VALUES
    ('Lorem ipsum', '"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.', NOW() - INTERVAL '5 days'),
    ('Tech Conference 2025', 'Join us for the biggest tech conference of the year.', NOW() - INTERVAL '2 days'),
    ('Go 1.25 Released', 'Most of its changes are in the implementation of the toolchain, runtime, and libraries.', NOW() - INTERVAL '1 day'),
    ('Database Optimization Tips', 'Learn how to optimize your PostgreSQL queries for better performance.', NOW()),
    ('Docker Best Practices', 'Essential tips for writing efficient Dockerfiles and managing containers.', NOW() + INTERVAL '1 day');

INSERT INTO Users (email, name, password, accountId) VALUES
    ('juan@uoc.edu', 'Juan Nadie', '', 1001),
    ('fulano@uoc.edu', 'Fulano de Tal', '', 1002),
    ('perengano@uoc.edu', 'Perengano', '', 1003);
