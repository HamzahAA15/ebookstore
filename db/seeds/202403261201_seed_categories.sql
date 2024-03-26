-- Seeder for inserting data into Categories table
INSERT INTO Categories (name)
SELECT name FROM (
    VALUES
    ('fiction'),
    ('non-fiction'),
    ('mystery'),
    ('romance'),
    ('science fiction'),
    ('fantasy'),
    ('biography'),
    ('self-help'),
    ('history'),
    ('thriller'),
    ('science')
) AS category(name)
WHERE NOT EXISTS (
    SELECT 1 FROM Categories c WHERE c.name = category.name
);