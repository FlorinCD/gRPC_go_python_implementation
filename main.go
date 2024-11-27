package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	userpb "my_project/proto"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	userpb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) RegisterUser(ctx context.Context, user *userpb.User) (*userpb.UserID, error) {
	// Insert user into the database
	_, err := s.db.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)", user.Id, user.Name, user.Email)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, err
	}
	log.Printf("User registered: %v", user)
	return &userpb.UserID{Id: user.Id}, nil
}

func (s *server) GetUser(ctx context.Context, id *userpb.UserID) (*userpb.User, error) {
	// Query user by ID
	row := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id.Id)
	user := &userpb.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		log.Printf("User not found: %v", id.Id)
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *server) DeleteUser(ctx context.Context, id *userpb.UserID) (*userpb.Empty, error) {
	// Delete user by ID
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return nil, err
	}
	log.Printf("User deleted: %v", id.Id)
	return &userpb.Empty{}, nil
}

func (s *server) ListUsers(empty *userpb.Empty, stream userpb.UserService_ListUsersServer) error {
	// Query all users
	rows, err := s.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		user := &userpb.User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Printf("Error scanning user: %v", err)
			return err
		}
		if err := stream.Send(user); err != nil {
			log.Printf("Error streaming user: %v", err)
			return err
		}
	}
	return nil
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Initialize server
	s := &server{db: db}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, s)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}
	log.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
