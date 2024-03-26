-- Filename: 202403261202_seed_books.sql

-- Seeder for inserting data into Books table
INSERT INTO books (title, author, price, category_id)
SELECT * FROM (
    VALUES
    ('to kill a mockingbird', 'harper lee', 12.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('1984', 'george orwell', 10.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('the great gatsby', 'f. scott fitzgerald', 11.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('the da vinci code', 'dan brown', 14.99, (SELECT id FROM categories WHERE name = 'mystery')),
    ('gone girl', 'gillian flynn', 13.99, (SELECT id FROM categories WHERE name = 'mystery')),
    ('the notebook', 'nicholas sparks', 10.99, (SELECT id FROM categories WHERE name = 'romance')),
    ('pride and prejudice', 'jane austen', 9.99, (SELECT id FROM categories WHERE name = 'romance')),
    ('the hobbit', 'j.r.r. tolkien', 14.99, (SELECT id FROM categories WHERE name = 'fantasy')),
    ('harry potter and the sorcerer''s stone', 'j.k. rowling', 15.99, (SELECT id FROM categories WHERE name = 'fantasy')),
    ('the hunger games', 'suzanne collins', 12.99, (SELECT id FROM categories WHERE name = 'science fiction')),
    ('a brief history of time', 'stephen hawking', 16.99, (SELECT id FROM categories WHERE name = 'science')),
    ('becoming', 'michelle obama', 19.99, (SELECT id FROM categories WHERE name = 'biography')),
    ('the subtle art of not giving a f*ck', 'mark manson', 17.99, (SELECT id FROM categories WHERE name = 'self-help')),
    ('sapiens: a brief history of humankind', 'yuval noah harari', 18.99, (SELECT id FROM categories WHERE name = 'history')),
    ('the catcher in the rye', 'j.d. salinger', 11.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('the girl with the dragon tattoo', 'stieg larsson', 13.99, (SELECT id FROM categories WHERE name = 'mystery')),
    ('the help', 'kathryn stockett', 12.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('the alchemist', 'paulo coelho', 10.99, (SELECT id FROM categories WHERE name = 'fiction')),
    ('the shining', 'stephen king', 14.99, (SELECT id FROM categories WHERE name = 'thriller')),
    ('educated', 'tara westover', 15.99, (SELECT id FROM categories WHERE name = 'biography'))
) AS book(title, author, price, category_id)
WHERE NOT EXISTS (
    SELECT 1 FROM books b WHERE b.title = book.title
);
