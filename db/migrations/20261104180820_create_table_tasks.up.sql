CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT
    , title TEXT
    , status TEXT
);

INSERT INTO tasks (title, status) 
VALUES ('Умыться и по чистить зубы.','в процессе')
    , ('Приготовить завтрак.', 'в планах')
    , ('По завтракать и выписть стаканчик валерьянки.','на рассмотрении')
;