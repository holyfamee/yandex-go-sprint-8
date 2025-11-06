package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())

	randRange = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcel.Number = id

	stored, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, parcel.Number, stored.Number)
	require.Equal(t, parcel.Client, stored.Client)
	require.Equal(t, parcel.Status, stored.Status)
	require.Equal(t, parcel.Address, stored.Address)
	require.Equal(t, parcel.CreatedAt, stored.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, stored.Address)
}

func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, stored.Status)
}

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

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		stored, ok := parcelMap[parcel.Number]
		require.True(t, ok)

		require.Equal(t, stored.Number, parcel.Number)
		require.Equal(t, stored.Client, parcel.Client)
		require.Equal(t, stored.Status, parcel.Status)
		require.Equal(t, stored.Address, parcel.Address)
		require.Equal(t, stored.CreatedAt, parcel.CreatedAt)
	}
}
