-- Write your migrate up statements here

CREATE VIEW view_book_title_author as 
SELECT 
  authors.name, 
  books.title 
FROM 
  books
INNER JOIN authors ON authors.id = books.author_id; 

---- create above / drop below ----

DROP VIEW IF EXISTS view_book_title_author

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
