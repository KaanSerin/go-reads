package database

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kaanserin/go-reads/internal/utils"
	_ "github.com/lib/pq"
)

type contextKey string

const DBContextKey contextKey = "db"

func (c contextKey) String() string {
	return string(c)
}

type User struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	RoleId    int       `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Book struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	Genre           string    `json:"genre" db:"genre"`
	PublicationDate time.Time `json:"publicationDate" db:"publication_date"`
	Publisher       string    `json:"publisher" db:"publisher"`
	ISBN            string    `json:"isbn" db:"isbn"`
	PageCount       string    `json:"pageCount" db:"page_count"`
	Language        string    `json:"language" db:"language"`
	Format          string    `json:"format" db:"format"`
}

type BookReview struct {
	ID        int       `json:"id" db:"id"`
	BookID    int       `json:"book_id" db:"book_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Score     int       `json:"score" db:"score"`
	Review    string    `json:"review" db:"review"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PostgresqlStorage struct {
	db *sqlx.DB
}

func GetLazyPaginatedResponsePG[V any](storage *PostgresqlStorage, r *http.Request, query string) ([]*V, error) {
	results := make([]*V, 0)

	page := r.URL.Query().Get("page")
	if page == "" || page == "0" {
		page = "1"
	}

	pageNum, err := strconv.ParseInt(page, 10, 0)
	if err != nil {
		return nil, err
	}

	pageLength := r.URL.Query().Get("pageLength")
	if pageLength == "" {
		pageLength = "15"
	}

	pageLengthNum, err := strconv.ParseInt(pageLength, 10, 0)
	if err != nil {
		return nil, err
	}

	offset := (pageNum - 1) * pageLengthNum

	queryWithLimit := fmt.Sprintf("%s order by id desc offset %d limit %s", query, offset, pageLength)

	err = storage.db.Select(&results, queryWithLimit)
	if err != nil {
		return nil, err
	}

	return results, nil
}

type Storage interface {
	GetUsers(r *http.Request) ([]*User, error)
	GetUserById(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	CreateUser(firstName, lastName, email, password string) (*User, error)
	DeleteUserById(int) error
	GetRoleById(int) (*Role, error)

	// Books
	GetBooks(r *http.Request) ([]*Book, error)

	// Book Reviews
	GetBookReviews(r *http.Request) ([]*BookReview, error)
	GetBookReviewById(id int) (*BookReview, error)
	DeleteBookReviewById(id int) error
	UpdateBookReview(id int, updateBookReviewDto UpdateBookReviewDto) (*BookReview, error)
}

func (storage *PostgresqlStorage) GetUserById(id int) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, role_id, created_at from users where id = $1", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.RoleId,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *PostgresqlStorage) GetUsers(r *http.Request) ([]*User, error) {
	query := "select id, first_name, last_name, email, role_id, created_at from users"
	users, err := GetLazyPaginatedResponsePG[User](storage, r, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (storage *PostgresqlStorage) GetUserByEmail(email string) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, role_id, password, created_at from users where email = $1", email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.RoleId,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *PostgresqlStorage) CreateUser(first_name, last_name, email, password string) (*User, error) {
	_, err := storage.db.Exec(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
		first_name,
		last_name,
		email,
		password)

	if err != nil {
		return nil, err
	}

	user, err := storage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (storage *PostgresqlStorage) DeleteUserById(id int) error {
	_, err := storage.db.Exec(fmt.Sprintf("DELETE FROM users where id = %d", id))
	return err
}

type UpdateUserDto struct {
	ID        int    `json:"id" validate:"nonzero"`
	FirstName string `json:"firstName" validate:"nonzero"`
	LastName  string `json:"lastName" validate:"nonzero"`
	Email     string `json:"email" validate:"nonzero"`
}

func (storage *PostgresqlStorage) UpdateUserById(id int, payload *UpdateUserDto) (*User, error) {
	user, err := storage.GetUserById(id)
	if err != nil {
		return nil, err
	}

	_, err = storage.db.Exec("UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4", payload.FirstName, payload.LastName, payload.Email, id)
	if err != nil {
		return nil, err
	}

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Email = payload.Email
	return user, nil
}

func (storage *PostgresqlStorage) GetRoleById(id int) (*Role, error) {
	var role *Role = &Role{}
	err := storage.db.QueryRow("SELECT id, name, created_at FROM roles WHERE id = $1", id).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (storage *PostgresqlStorage) GetBooks(r *http.Request) ([]*Book, error) {
	return GetLazyPaginatedResponsePG[Book](storage, r, "SELECT id, title, author, genre, publication_date, publisher, isbn, page_count, language, format FROM books")
}

func (storage *PostgresqlStorage) GetBookById(id int) (*Book, error) {
	var book *Book = &Book{}
	err := storage.db.Get(book, "SELECT * from books where id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

type UpdateBookDto struct {
	Title           string    `json:"title" validate:"nonzero" db:"title"`
	Author          string    `json:"author" validate:"nonzero" db:"author"`
	Genre           string    `json:"genre" validate:"nonzero" db:"genre"`
	PublicationDate time.Time `json:"publicationDate" validate:"nonzero" db:"publication_date"`
	Publisher       string    `json:"publisher" validate:"nonzero" db:"publisher"`
	ISBN            string    `json:"isbn" validate:"nonzero" db:"isbn"`
	PageCount       string    `json:"pageCount" validate:"nonzero" db:"page_count"`
	Language        string    `json:"language" validate:"nonzero" db:"language"`
	Format          string    `json:"format" validate:"nonzero" db:"format"`
}

func (storage *PostgresqlStorage) UpdateBookById(id int, payload *UpdateBookDto) (*Book, error) {
	result, err := storage.db.NamedExec(fmt.Sprintf(`UPDATE books
	SET title = :title, author = :author, genre = :genre, publication_date = :publication_date,
	publisher = :publisher, isbn = :isbn, page_count = :page_count, language = :language, format = :format
	WHERE id = %d`, id), payload)

	if err != nil {
		return nil, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, &utils.CustomError{
			Message: fmt.Sprintf("No book found for the given id %d", id),
		}
	}

	var book *Book = &Book{}
	err = storage.db.Get(book, "SELECT * from books WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (storage *PostgresqlStorage) DeleteBookById(id int) error {
	_, err := storage.GetBookById(id)
	if err != nil {
		return err
	}

	result, err := storage.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAff == 0 {
		return &utils.CustomError{
			Message: fmt.Sprintf("No book found for the given id %d", id),
		}
	}

	return nil
}

func (storage *PostgresqlStorage) GetBookReviews(r *http.Request) ([]*BookReview, error) {
	return GetLazyPaginatedResponsePG[BookReview](storage, r, "SELECT * from book_reviews")
}

func (storage *PostgresqlStorage) GetBookReviewById(id int) (*BookReview, error) {
	var bookReview *BookReview = &BookReview{}

	err := storage.db.Get(bookReview, "SELECT * FROM book_reviews WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}

	return bookReview, nil
}

type CreateBookReviewDto struct {
	BookID int    `json:"bookId" db:"book_id" validate:"nonzero"`
	UserID int    `json:"userId" db:"user_id"`
	Score  int    `json:"score" db:"score" validate:"nonzero"`
	Review string `json:"review" db:"review" validate:"nonzero"`
}

func (storage *PostgresqlStorage) CreateBookReview(createUserDto *CreateBookReviewDto) (*BookReview, error) {
	var id int
	err := storage.db.QueryRow("INSERT INTO book_reviews (book_id, user_id, score, review) VALUES ($1, $2, $3, $4) RETURNING id",
		createUserDto.BookID, createUserDto.UserID, createUserDto.Score, createUserDto.Review).Scan(&id)
	if err != nil {
		return nil, err
	}

	return storage.GetBookReviewById(id)
}

func (storage *PostgresqlStorage) DeleteBookReviewById(id int) error {
	result, err := storage.db.Exec("DELETE FROM book_reviews WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return &utils.CustomError{
			Message: "Book review not found",
		}
	}

	return nil
}

type UpdateBookReviewDto struct {
	Score     int       `json:"score" db:"score" validate:"nonzero"`
	Review    string    `json:"review" db:"review" validate:"nonzero"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (storage *PostgresqlStorage) UpdateBookReview(id int, updateBookReviewDto UpdateBookReviewDto) (*BookReview, error) {
	updateBookReviewDto.UpdatedAt = time.Now()
	result, err := storage.db.NamedExec(fmt.Sprintf("UPDATE book_reviews SET score = :score, review = :review, updated_at = :updated_at WHERE id = %d", id), updateBookReviewDto)
	if err != nil {
		return nil, err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAff == 0 {
		return nil, &utils.CustomError{
			Message: "Book review not found",
		}
	}

	var bookReview *BookReview = &BookReview{}
	if err := storage.db.Get(bookReview, "SELECT * FROM book_reviews WHERE id = $1", id); err != nil {
		return nil, err
	}

	return bookReview, nil
}

func NewPostgresStorage() (*PostgresqlStorage, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &PostgresqlStorage{
		db: db,
	}, nil
}

func GetPgStorageFromRequest(r *http.Request) (*PostgresqlStorage, error) {
	db := r.Context().Value(DBContextKey).(*PostgresqlStorage)
	return db, nil
}
