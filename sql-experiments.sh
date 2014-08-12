#! /bin/bash

sqlite3 .redo.db <<EOF
.mode column

UPDATE files SET generation=2 WHERE id=1;

SELECT p.id, p.path, c.id, c.path, p.generation, p.step, c.generation, c.step
FROM files AS p JOIN deps ON p.id = deps.to_make JOIN files AS c ON deps.you_need = c.id
WHERE p.generation = 2
AND c.generation != p.generation;

UPDATE files SET generation=2, step=0
WHERE id IN (
  SELECT c.id
  FROM files AS p JOIN deps ON p.id = deps.to_make JOIN files AS c ON deps.you_need = c.id
  WHERE p.generation = 2
  AND c.generation != p.generation);

SELECT * FROM files;

.exit
EOF

: <<EOF

PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA temp_store = MEMORY;

CREATE TEMPORARY TABLE demands (
   file INTEGER REFERENCES files(id) ON UPDATE CASCADE ON DELETE CASCADE,
   idx INTEGER NOT NULL);

SELECT id, idx, generation, step
FROM files LEFT JOIN demands ON file = id
WHERE path="t/all";


EOF
