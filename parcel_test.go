package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа).
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел.
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку.
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки.
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// проверяем добавление посылки
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// проверяем получение посылки
	parcelGet, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, parcel, parcelGet)

	// проверяем удаление посылки
	err = store.Delete(parcelGet.Number)
	require.NoError(t, err)
	_, err = store.Get(parcelGet.Number)
	require.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса.
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// проверяем добавление посылки
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// проверяем изменение адреса посылки
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	assert.NoError(t, err)

	parcleGet, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, newAddress, parcleGet.Address)
}

// TestSetStatus проверяет обновление статуса.
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// проверяем добавление посылки
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// проверяем изменение статуса посылки
	err = store.SetStatus(parcel.Number, ParcelStatusSent)
	assert.NoError(t, err)

	parcleGet, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, parcleGet.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента.
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		// проверяем добавление посылки
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// проверяем получение посылок и их кол-во
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcelMap))

	for _, parcel := range storedParcels {
		// проверяем соответстве данных посылок
		p, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, p, parcel)
	}
}
