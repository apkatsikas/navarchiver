ATTACH DATABASE 'archive_run.db' AS archive_run;

CREATE TABLE IF NOT EXISTS archive_run (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_run DATE NOT NULL,
    CONSTRAINT id_unique UNIQUE (id)
);
