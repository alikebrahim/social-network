package main

import "log"

func main() {
	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	// test accounts and data for development uncomment to add
	//store.AddTestAccount()
	server := NewAPIServer(":3000", store)
	server.Run()
}
