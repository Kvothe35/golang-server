package main

import (
	"context"
	"http-server/server"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
)

func main() {

	projID := "golang-integration"

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projID)

	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	srv := server.NewServer(client)
	http.ListenAndServe(":8080", srv)

	// type Task struct {
	// 	Desc    string    `datastore:"description"`
	// 	Created time.Time `datastore:"created"`
	// 	Done    bool      `datastore:"done"`
	// 	id      int64     // The integer ID used in the datastore.
	// }

	// task := &Task{
	// 	Desc:    "adding from GO application",
	// 	Created: time.Now(),
	// }
	// key := datastore.IncompleteKey("Task", nil)

	// outKey, err := client.Put(ctx, key, task)
	// if err != nil {
	// 	log.Printf("Failed to create task: %v", err)
	// 	return
	// }
	// fmt.Printf("Created new task with ID %d\n", outKey.ID)
}
