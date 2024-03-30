INSERT INTO crawls (url, title, status_code, timestamp)
VALUES ('https://millandcommons.com', 'Mill+Commons', 200, '2024-01-01 01:00:00');

INSERT INTO crawl_headings(crawl_id, text, tag)
VALUES (1, 'title text', 'h1'), (1, 'subtitle text', 'h2');

INSERT INTO crawl_links(crawl_id, text, url)
VALUES (1, 'view dining tables', 'https://https://millandcommons.com/products/ezra-dining-table'), (1, 'view catalog', 'https://millandcommons.com/collections/all');