INSERT INTO News (title, body, datetime) VALUES
    ('Lorem ipsum', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.', NOW() - INTERVAL '5 days'),
    ('Wannacry: el ransomware que alertó a todo el mundo', 'Wannacry afectó a varios sistemas informáticos de un montón de empresa en septiembre de 2018.', NOW() - INTERVAL '2 days'),
    ('Go 1.25 Released', 'Most of its changes are in the implementation of the toolchain, runtime, and libraries.', NOW() - INTERVAL '1 day'),
    ('Listado de empresas afectadas por vulnerabilidades SQLi', 'Hoy vamos a ver un gran listado de casos reales de empresas que fueron afectada por esta vulnerabilidad tan popular.', NOW()),
    ('Docker Best Practices', 'Essential tips for writing efficient Dockerfiles and managing containers.', NOW() + INTERVAL '1 day');

INSERT INTO Users (email, name, password, accountId) VALUES
    ('juan@uoc.edu', 'Juan Nadie', 'b9c950640e1b3740e98acb93e669c65766f6670dd1609ba91ff41052ba48c6f3', 1001),
    ('fulano@uoc.edu', 'Fulano de Tal', '1532e76dbe9d43d0dea98c331ca5ae8a65c5e8e8b99d3e2a42ae989356f6242a', 1002),
    ('perengano@uoc.edu', 'Perengano', '29090710e3cc0fd7418cbe672b939eab7b5d1ad36de51fdc6dd93a999d2bfe7a', 1003);
