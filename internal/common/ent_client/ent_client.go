package ent_client

import (
	"context"
	"fmt"
	"go-backend/ent"
	"go-backend/internal/common/env"
	"log"

	_ "go-backend/ent/runtime"

	_ "github.com/go-sql-driver/mysql"
)

func New(env *env.Env) *ent.Client {
	client, err := ent.Open("mysql", env.DatabaseUrl)
	if err != nil {
		log.Fatalf(" [ENT] failed opening connection to mysql: %v", err)
	}

	ctx := context.Background()
	_, err = client.QueryContext(ctx, "SELECT 1 + 1")

	if err != nil {
		log.Fatalf(" XXX [ENT] failed opening connection to mysql: %v", err)
	}

	fmt.Println("OK [ENT] Connection To MySQL Successfully")

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf(" [ENT] failed creating schema resources: %v", err)
	}
	return client
}
