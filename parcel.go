package main

import (
	"database/sql"
)

// ParcelStore - структура для рабоыт с БД.
type ParcelStore struct {
	db *sql.DB
}

// NewParcelStorece создает новый объект типа ParcelStoree.
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в БД.
// Возвращает сгенерированный таблицой идентификатор и возможную ошибку.
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

// Get считывает из тБД данные на посылку по указанному идентификатору.
// Возвращает полученные данные в объекте структуры типа Parcel и возможную ошибку.
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

// GetByClient считывает из БД данные на все посылки по указанному идентификатору клиента.
// Возвращает полученные данные в виде массива объектов структуры типа Parcel и возможную ошибку.
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))

	var res []Parcel
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{
			Client: client,
		}
		err = rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SetStatus меняет статус посылки в БД на указанный по указанному идентификатору.
// Возвращает возможную ошибку.
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

// SetAddress меняет адрес посылки в БД на указанный по указанному идентификатору.
// Менять адрес можно только если статус посылки "registered".
// Возвращает возможную ошибку.
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE (number = :number) AND (status = :status)",
		sql.Named("number", number),
		sql.Named("address", address),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

// Delete удалет посылку из БД по указанному идентификатору.
// Удалять можно только если статус посылки "registered".
// Возвращает возможную ошибку.
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE (number = :number) AND (status = :status)",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
