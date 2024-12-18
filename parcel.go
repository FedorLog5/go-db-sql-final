package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT * FROM parcel WHERE number = ?`

	row := s.db.QueryRow(query, number)
	
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	query := `SELECT * FROM parcel WHERE client = ?`

	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	_, err := s.db.Exec(query, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := `UPDATE parcel SET address = ? WHERE number = ? AND status = 'registered'`
	_, err := s.db.Exec(query, address, number)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	query := `DELETE FROM parcel WHERE number = ? AND status = 'registered'`
	_, err := s.db.Exec(query, number)
	return err
}
