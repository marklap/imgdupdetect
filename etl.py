import os
import sys
import sqlite3

db_file = './imgdd.sqlite'
data_file = './photo_md5sums.txt'
photo_path = '/volume0/photo'

create_table_sql = 'CREATE TABLE IF NOT EXISTS images (filePath TEXT NOT NULL PRIMARY KEY, md5sum TEXT NOT NULL)'
insert_row_sql = 'INSERT OR REPLACE INTO images (filePath, md5sum) VALUES (?, ?)'
duplicate_sql = '''\
SELECT count(filePath) AS fileCount, md5sum, GROUP_CONCAT(filePath, char(10)) AS files FROM images
GROUP BY md5sum
HAVING fileCount > 1
ORDER BY fileCount DESC'''


def file_path_depth(path):
    depth = 0
    while (path != '/'):
        path = os.path.dirname(path)
        depth += 1
    return depth


def main():
    db = sqlite3.connect(db_file)
    c = db.cursor()
    c.execute(create_table_sql)

    with open(data_file) as fp:
        for line in fp:
            md5_sum, file_path = line.strip().split(None, 1)
            file_ext = os.path.splitext(file_path)[1].lower()

            if file_ext in ('.jpg', '.jpeg', '.gif', '.png'):
                c.execute(insert_row_sql, [
                    os.path.realpath(os.path.join(photo_path, file_path)),
                    md5_sum])

    db.commit()

    c = db.cursor()
    for row in c.execute(duplicate_sql):
        count, md5sum, paths = row
        for path in paths.splitlines():
            print(md5sum, ' ---- ', path)
    db.commit()


if __name__ == "__main__":
    sys.exit(main())
