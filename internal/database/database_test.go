package database

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {
	godotenv.Load("../../.env")
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestCreateUserWithValidInputsAndDelete", func(t *testing.T) {
		// Create a user
		user, err := storage.CreateUser("TestFirstName", "TestLastName", "testmail@mail.com", "pass")
		if err != nil {
			t.Error(err)
		}

		// Check if the user is created
		if user == nil {
			t.Error("User is nil")
			return
		}

		// Check if the user has an ID
		if user.ID == 0 {
			t.Error("User ID is empty")
		}

		// Check if the user created at is not empty
		if user.CreatedAt.IsZero() {
			t.Error("User created at is empty")
		}

		// Delete the user for cleanup
		err = storage.DeleteUserById(user.ID)
		if err != nil {
			t.Error(err)
		}
	})
}
