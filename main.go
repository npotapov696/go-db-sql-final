package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Статусы посылки.
const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

// Parcel - структура, содержащая информацию о посылке, поля которой соответствуют
// атрибутам таблицы базы данных.
type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

// ParcelService является оберткой над структурой ParcelStore, работающей с базой данных.
type ParcelService struct {
	store ParcelStore
}

// NewParcelService создает новый объект типа ParcelService.
func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

// Register принимает на вход идентификатор клиента и адрес доставки, из которых генерирует
// новый элемент типа Parcel, который добавляет в базу данных и выводит информацию о нём на консоль.
// Возвращает сгенерированный элемент типа Parcel и возможную ошибку.
func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// добавляем новую посылку в таблицу базы данных
	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err
	}

	// заполняем поле идентификатора посылки, сгенерированное таблицой базы данных
	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

// PrintClientParcels формирует и выводит на консоль список посылок по указанному
// идентификатору клиента из базы данных. Возвращает возможную ошибку.
func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()

	return nil
}

// NextStatus меняет статус посылки в базе данных по указанному идентификатору.
// Выводит новый статус посылки на консоль. Возвращает возможную ошибку.
func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent // если статус "registered", поменять на "sent"
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered // если статус "sent", поменять на "delivered"
	case ParcelStatusDelivered:
		return nil // если статус "delivered", поменять на nil
	}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	return s.store.SetStatus(number, nextStatus)
}

// ChangeAddress меняет адрес посылки в базе данных на указанный по указанному идентификатору.
// Возвращает возможную ошибку.
func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

// Delete удаляет посылку из базы данных по указанному идентификатору.
// Возвращает возможную ошибку.
func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

func main() {
	// соединяемся с базой данных
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// создаем объекты для работы с базой данных
	store := NewParcelStore(db)
	service := NewParcelService(store)

	// регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
