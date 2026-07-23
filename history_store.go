package main

import "database/sql"

// AddHistory stores metadata only. Files are kept on disk.
func AddHistory(db *sql.DB, h History) error {
	_, err := db.Exec(`INSERT INTO history(type,direction,content,filename,filepath,mime,size) VALUES(?,?,?,?,?,?,?)`,
		h.Type, h.Direction, h.Content, h.Filename, h.Filepath, h.Mime, h.Size)
	return err
}

func ListHistory(db *sql.DB, page, size int) ([]History, error) {
	rows, err := db.Query(`SELECT id,type,direction,content,filename,filepath,mime,size,created_at FROM history ORDER BY id DESC LIMIT ? OFFSET ?`, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []History{}
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.Type, &h.Direction, &h.Content, &h.Filename, &h.Filepath, &h.Mime, &h.Size, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
