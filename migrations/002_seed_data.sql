INSERT INTO News (title, body, datetime) VALUES
    ('Lorem ipsum', '"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.', NOW() - INTERVAL '5 days'),
    ('Wannacry: el ransomware que alertó a todo el mundo', 'Wannacry afectó a varios sistemas informáticos de un montón de empresa en septiembre de 2018.', NOW() - INTERVAL '2 days'),
    ('Go 1.25 Released', 'Most of its changes are in the implementation of the toolchain, runtime, and libraries.', NOW() - INTERVAL '1 day'),
    ('Listado de empresas afectadas por vulnerabilidades SQLi', 'Hoy vamos a ver un gran listado de casos reales de empresas que fueron afectada por esta vulnerabilidad tan popular.', NOW()),
    ('Docker Best Practices', 'Essential tips for writing efficient Dockerfiles and managing containers.', NOW() + INTERVAL '1 day');

INSERT INTO Users (email, name, password, accountId) VALUES
    ('juan@uoc.edu', 'Juan Nadie', '', 1001),
    ('fulano@uoc.edu', 'Fulano de Tal', '', 1002),
    ('perengano@uoc.edu', 'Perengano', '', 1003);
