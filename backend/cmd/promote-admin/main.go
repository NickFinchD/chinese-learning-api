// Command promote-admin grants (or revokes, with -revoke) admin access to a
// user by email. This is the only way to create the first admin — there is
// deliberately no self-service "become admin" endpoint — and it's also the
// intended way to add/remove admins after that, so admin access never
// depends on remembering one-off SQL.
//
//	go run ./cmd/promote-admin -email you@example.com
//	go run ./cmd/promote-admin -email you@example.com -revoke
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
)

func main() {
	email := flag.String("email", "", "email of the user to promote (required)")
	revoke := flag.Bool("revoke", false, "revoke admin access instead of granting it")
	flag.Parse()

	if *email == "" {
		log.Fatal("-email is required")
	}

	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	tag, err := db.Exec(ctx, `UPDATE users SET is_admin = $1 WHERE email = $2`, !*revoke, *email)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	if tag.RowsAffected() == 0 {
		log.Fatalf("no user found with email %q", *email)
	}

	if *revoke {
		fmt.Printf("Revoked admin access for %s\n", *email)
	} else {
		fmt.Printf("Granted admin access to %s\n", *email)
	}
}
