package database

import (
	"GolangTD/session"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const ERROR_VALIDATE_DELETE_POST = "you don't have rights to delete this post"

type User struct { //структура данных пользователя
	UserID    string
	Username  string
	Firstname string
	Lastname  string
	Date      string
	TokenID   string
	Password  string
	Email     string
}

type Post struct { //структура данных поста
	Title          string
	Content        string
	DateOfCreation string
	DateOfPublish  string
	ContentID      string
	UserID         string
}

type Data struct {
	DataPosts []Post
	DataUser  *User
}

func InitDatabase(dataName string) (*sql.DB, error) { //функция-инициализатор работы с бд, возвращает объект типа *sql.DB
	return nil, nil
}

func NewUser(userID, username, firstname, lastname, date, tokenID, password, email string) *User {
	return &User{userID, username, firstname, lastname, date, tokenID, password, email}
}

func NewPost(title, content, dateOfCreation, dateOfPublish, contentID, userID string) *Post {
	return &Post{title, content, dateOfCreation, dateOfPublish, contentID, userID}
}

func (u *User) ShowDeb() {
	fmt.Println("UserID: " + u.UserID)
	fmt.Println("Username: " + u.Username)
	fmt.Println("Firstname: " + u.Firstname)
	fmt.Println("Lastname: " + u.Lastname)
	fmt.Println("Date: " + u.Date)
	fmt.Println("TokenID: " + u.TokenID)
	fmt.Println("Password: " + u.Password)
	fmt.Println("Email: " + u.Email)

}

func (p *Post) ShowDeb() {
	fmt.Println("Title: " + p.Title)
	fmt.Println("Content: " + p.Content)
	fmt.Println("DateOfCreation: " + p.DateOfCreation)
	fmt.Println("DateOfPublish: " + p.DateOfPublish)
	fmt.Println("ContentID: " + p.ContentID)
	fmt.Println("UserID: " + p.UserID)
}

func FindPostsByUserID(userID string) ([]Post, error) {
	query := `SELECT * FROM posts WHERE UserID = ?`
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Title, &p.Content, &p.DateOfCreation, &p.DateOfPublish, &p.ContentID, &p.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return posts, err
	}
	return posts, nil
}

func (p *Post) FindPostByContentID(contentID string) error {
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return err
	}
	row := db.QueryRow("SELECT * FROM posts WHERE ContentID = ?", contentID)
	err = row.Scan(&p.Title, &p.Content, &p.DateOfCreation, &p.DateOfPublish, &p.ContentID, &p.UserID)
	if err != nil {
		return err
	}
	return nil
}

func Authorisation(username, password string) (*User, error) { //возвращает указатель на структуру Users при правильном входе, а также код ошибки.
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT UserID FROM users WHERE Username=? AND Password=?", username, password)
	var id string
	err = row.Scan(&id)
	if err != nil {
		return nil, errors.New("указанного пользователя нет в базе")
	}
	u := &User{}
	u.Get(id)
	return u, nil
}

func (u *User) Get(id string) error { //функция для поиска данных пользователя по его идентификационному номеру
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return err
	}

	row := db.QueryRow("SELECT * FROM users WHERE UserID=?", id)
	if err != nil {
		return err
	}
	err = row.Scan(&u.UserID, &u.Username, &u.Firstname, &u.Lastname, &u.Date, &u.TokenID, &u.Password, &u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Create() error { //создать нового пользователя
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users VALUES(?,?,?,?,?,?,?,?)", u.UserID, u.Username, u.Firstname, u.Lastname, u.Date, u.TokenID, u.Password, u.Email)
	return err
}

func (p *Post) Create() error { //создать новый пост
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO posts VALUES(?,?,?,?,?,?)", p.Title, p.Content, p.DateOfCreation, p.DateOfPublish, p.ContentID, p.UserID)
	return err
}

func FormData(r *http.Request) (*Data, error) {
	cookie, err := session.CheckCookie(r)
	if err != nil {
		return nil, err
	}
	userID := session.SessionInMemory.Get(cookie.Value).UserID
	posts, err := FindPostsByUserID(userID)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = user.Get(userID)
	data := &Data{posts, user}
	return data, nil
}

func DeletePost(contentID string) error {
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM posts WHERE ContentID = ?", contentID)
	if err != nil {
		return err
	}
	return nil
}

func PostDeleteCheck(userID, contentID string) error {
	p := &Post{}
	p.FindPostByContentID(contentID)
	if p.UserID == userID {
		return nil
	} else {
		return errors.New(ERROR_VALIDATE_DELETE_POST)
	}
}

func EditPost(title, content, date, contentID string) error {
	db, err := sql.Open("sqlite3", "../data/database/data.db")
	if err != nil {
		return errors.New("ошибка в октрытии базы данных")
	}
	_, err = db.Exec("UPDATE posts SET Title=?, Content=?, DateOfPublish=? WHERE ContentID = ?", title, content, date, contentID)
	if err != nil {
		return errors.New("ошибка транзакции")
	}
	return nil
}
