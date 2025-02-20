package main

import (
	"database/sql"
	"errors"
)

// ParcelStore - структура для рабоыт с базой данных
type ParcelStore struct {
	db *sql.DB
}

// NewParcelStorece создает новый объект типа ParcelStoree.
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в таблицу базы данных.
// Возвращает сгенерированный таблицой идентификатор и возможную ошибку.
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
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

// Get считывает из таблицы базы данных данные на посылку по указанному идентификатору.
// Возвращает полученные данные в объекте структуры типа Parcel и возможную ошибку.
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("select client, status, address, created_at from parcel where number = :number",
		sql.Named("number", number))

	p := Parcel{
		Number: number,
	}
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

// GetByClient считывает из таблицы базы данных данные на все посылки по указанному идентификатору клиента.
// Возвращает полученные данные в виде массива объектов структуры типа Parcel и возможную ошибку.
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("select number, status, address, created_at from parcel where client = :client",
		sql.Named("client", client))

	var res []Parcel
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{
			Client: client,
		}
		err = rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	return res, err
}

// SetStatus меняет статус посылки на указанный по указанному идентификатору.
// Возвращает возможную ошибку.
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("update parcel set status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	row := s.db.QueryRow("select status from parcel where number = :number",
		sql.Named("number", number))
	var status string
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New(`cant update the address, status should be "registered"`)
	}
	_, err = s.db.Exec("update parcel set address = :address where number = :number",
		sql.Named("number", number),
		sql.Named("address", address))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	row := s.db.QueryRow("select status from parcel where number = :number",
		sql.Named("number", number))
	var status string
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New(`cant delete the parcel, status should be "registered"`)
	}
	_, err = s.db.Exec("delete from parcel where number = :number",
		sql.Named("number", number))

	return err
}
