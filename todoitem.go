package main

import "database/sql"

type Item struct {
	ID      int
	Content string
	Keyword string
}

type TodoItem struct {
	db *sql.DB
}

func NewTodoItem(db *sql.DB) *TodoItem {
	return &TodoItem{db: db}
}

func (todoitem *TodoItem) CreateTable() error {
	const sql = `CREATE TABLE IF NOT EXISTS todoitems(
		id int auto_increment,
		content text,
		keyword text,
		primary key(id)
	);`
	_, err := todoitem.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (todolist *TodoItem) SaveItem(item *Item) error {
	const sql = `INSERT INTO todoitems(content, keyword) values (?,?)`
	_, err := todolist.db.Exec(sql, item.Content, item.Keyword)
	if err != nil {
		return err
	}
	return nil
}

func (todoitem *TodoItem) GetItems(keyword string) ([]*Item, error) {
	const sql = `SELECT * FROM todoitems WHERE keyword = ?`
	rows, err := todoitem.db.Query(sql, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Content, &item.Keyword)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, err
}

func (todoitem *TodoItem) DeleteItems(keyword string) error {
	const sql = `DELETE FROM todoitems where keyword = ?`
	_, err := todoitem.db.Exec(sql, keyword)
	if err != nil {
		return err
	}
	return nil
}
