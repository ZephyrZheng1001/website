package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtSecret = []byte("zephyr-jwt-secret-key-2026")

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Tags      string    `json:"tags"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func authCheck(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false, ""
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid { return false, "" }
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { return false, "" }
	username, _ := claims["username"].(string)
	return true, username
}

// ====================== Public Handlers ======================

func getArticles(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	tag := r.URL.Query().Get("tag")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 { page = 1 }
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 { limit = 10 }
	offset := (page - 1) * limit

	var conditions []string
	var args []interface{}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}
	if tag != "" {
		conditions = append(conditions, "tags LIKE ?")
		args = append(args, "%"+tag+"%")
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	query := "SELECT id, title, summary, tags, category, created_at, updated_at FROM articles" + where + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	countQuery := "SELECT COUNT(*) FROM articles" + where

	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
	if err != nil { writeJSON(w, 500, APIResponse{Success: false, Message: err.Error()}); return }
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		rows.Scan(&a.ID, &a.Title, &a.Summary, &a.Tags, &a.Category, &a.CreatedAt, &a.UpdatedAt)
		articles = append(articles, a)
	}
	if articles == nil { articles = []Article{} }

	var total int
	countArgs := args[:len(args)-2]
	db.QueryRow(countQuery, countArgs...).Scan(&total)

	writeJSON(w, 200, APIResponse{Success: true, Data: map[string]interface{}{
		"articles": articles, "total": total, "page": page, "limit": limit,
	}})
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/articles/")
	var a Article
	err := db.QueryRow("SELECT id, title, content, summary, tags, category, created_at, updated_at FROM articles WHERE id=?", id).
		Scan(&a.ID, &a.Title, &a.Content, &a.Summary, &a.Tags, &a.Category, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows { writeJSON(w, 404, APIResponse{Success: false, Message: "文章不存在"}); return }
	if err != nil { writeJSON(w, 500, APIResponse{Success: false, Message: err.Error()}); return }
	// Get prev and next articles in same category
	var prevID, nextID int
	var prevTitle, nextTitle string
	db.QueryRow("SELECT id, title FROM articles WHERE category=? AND id < ? ORDER BY id DESC LIMIT 1", a.Category, id).Scan(&prevID, &prevTitle)
	db.QueryRow("SELECT id, title FROM articles WHERE category=? AND id > ? ORDER BY id ASC LIMIT 1", a.Category, id).Scan(&nextID, &nextTitle)

	writeJSON(w, 200, APIResponse{Success: true, Data: map[string]interface{}{
		"article": a,
		"prev": map[string]interface{}{"id": prevID, "title": prevTitle},
		"next": map[string]interface{}{"id": nextID, "title": nextTitle},
	}})
}

// ====================== Admin Handlers ======================

func adminLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username=?", input.Username).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		writeJSON(w, 401, APIResponse{Success: false, Message: "用户名或密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": input.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(jwtSecret)
	writeJSON(w, 200, APIResponse{Success: true, Data: map[string]string{"token": tokenStr}})
}

func adminChangePassword(w http.ResponseWriter, r *http.Request, username string) {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username=?", username).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.OldPassword)) != nil {
		writeJSON(w, 400, APIResponse{Success: false, Message: "旧密码错误"})
		return
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	db.Exec("UPDATE users SET password_hash=? WHERE username=?", string(newHash), username)
	writeJSON(w, 200, APIResponse{Success: true, Message: "密码已修改"})
}

func adminGetArticles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, summary, tags, category, created_at, updated_at FROM articles ORDER BY created_at DESC")
	if err != nil { writeJSON(w, 500, APIResponse{Success: false, Message: err.Error()}); return }
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		rows.Scan(&a.ID, &a.Title, &a.Content, &a.Summary, &a.Tags, &a.Category, &a.CreatedAt, &a.UpdatedAt)
		articles = append(articles, a)
	}
	if articles == nil { articles = []Article{} }
	writeJSON(w, 200, APIResponse{Success: true, Data: articles})
}

func adminCreateArticle(w http.ResponseWriter, r *http.Request) {
	var a Article
	json.NewDecoder(r.Body).Decode(&a)
	res, err := db.Exec("INSERT INTO articles (title, content, summary, tags, category) VALUES (?,?,?,?,?)",
		a.Title, a.Content, a.Summary, a.Tags, a.Category)
	if err != nil { writeJSON(w, 500, APIResponse{Success: false, Message: err.Error()}); return }
	id, _ := res.LastInsertId()
	writeJSON(w, 200, APIResponse{Success: true, Data: map[string]int64{"id": id}})
}

func adminUpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/articles/")
	var a Article
	json.NewDecoder(r.Body).Decode(&a)
	_, err := db.Exec("UPDATE articles SET title=?, content=?, summary=?, tags=?, category=? WHERE id=?",
		a.Title, a.Content, a.Summary, a.Tags, a.Category, id)
	if err != nil { writeJSON(w, 500, APIResponse{Success: false, Message: err.Error()}); return }
	writeJSON(w, 200, APIResponse{Success: true})
}

func adminDeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/articles/")
	db.Exec("DELETE FROM articles WHERE id=?", id)
	writeJSON(w, 200, APIResponse{Success: true})
}

// ====================== Router ======================

func router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" { w.WriteHeader(http.StatusOK); return }

	path := r.URL.Path

	// Public routes
	if path == "/api/articles" && r.Method == "GET" { getArticles(w, r); return }
	if strings.HasPrefix(path, "/api/articles/") && r.Method == "GET" { getArticle(w, r); return }

	// Admin login
	if path == "/api/admin/login" && r.Method == "POST" { adminLogin(w, r); return }

	// Admin routes
	ok, username := authCheck(r)
	if !ok {
		writeJSON(w, 401, APIResponse{Success: false, Message: "未授权"})
		return
	}

	if path == "/api/admin/password" && r.Method == "PUT" { adminChangePassword(w, r, username); return }

	if path == "/api/admin/articles" && r.Method == "GET" { adminGetArticles(w, r); return }
	if path == "/api/admin/articles" && r.Method == "POST" { adminCreateArticle(w, r); return }
	if strings.HasPrefix(path, "/api/admin/articles/") && r.Method == "PUT" { adminUpdateArticle(w, r); return }
	if strings.HasPrefix(path, "/api/admin/articles/") && r.Method == "DELETE" { adminDeleteArticle(w, r); return }

	writeJSON(w, 404, APIResponse{Success: false, Message: "接口不存在"})
}

// ====================== Init DB ======================

func initDB() {
	initDB, err := sql.Open("mysql", "root:zephyr123@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil { log.Fatal("DB init connect failed:", err) }
	initDB.Exec("CREATE DATABASE IF NOT EXISTS zephyr_website CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	initDB.Close()

	db, err = sql.Open("mysql", "root:zephyr123@tcp(127.0.0.1:3306)/zephyr_website?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil { log.Fatal("DB connect failed:", err) }
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	if err = db.Ping(); err != nil { log.Fatal("DB ping failed:", err) }

	db.Exec(`CREATE TABLE IF NOT EXISTS articles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		summary VARCHAR(500) DEFAULT '',
		tags VARCHAR(255) DEFAULT '',
		category VARCHAR(50) DEFAULT 'blog',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	db.Exec("ALTER TABLE articles ADD COLUMN IF NOT EXISTS category VARCHAR(50) DEFAULT 'blog' AFTER tags")

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Exec("INSERT INTO users (username, password_hash) VALUES (?,?)", "admin", string(hash))
		fmt.Println("[init] Default admin: admin / admin123")
	}
	fmt.Println("[init] Database ready")
}

func main() {
	initDB()
	http.HandleFunc("/api/", router)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
