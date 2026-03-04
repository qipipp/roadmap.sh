package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DSN string `json:"mysql_dsn"`
}

type PostInput struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func selectPostByID(db *sql.DB, id int64) (Post, error) {
	var p Post
	var tagsRaw []byte
	err := db.QueryRow(
		`SELECT id, title, content, category, tags, created_at, updated_at
		 FROM posts WHERE id=?`,
		id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Category, &tagsRaw, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return Post{}, err
	}
	if len(tagsRaw) > 0 {
		_ = json.Unmarshal(tagsRaw, &p.Tags)
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}
	return p, nil
}

func createPost(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var in PostInput
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Title) == "" ||
		strings.TrimSpace(in.Content) == "" ||
		strings.TrimSpace(in.Category) == "" {
		http.Error(w, "title/content/category required", http.StatusBadRequest)
		return
	}
	if in.Tags == nil {
		in.Tags = []string{}
	}
	tagsJSON, _ := json.Marshal(in.Tags)

	res, err := db.Exec(
		`INSERT INTO posts (title, content, category, tags) VALUES (?,?,?,?)`,
		in.Title, in.Content, in.Category, tagsJSON,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := selectPostByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func listPosts(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT id, title, content, category, tags, created_at, updated_at
		 FROM posts
		 ORDER BY id DESC`,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		var p Post
		var tagsRaw []byte
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &tagsRaw, &p.CreatedAt, &p.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(tagsRaw) > 0 {
			_ = json.Unmarshal(tagsRaw, &p.Tags)
		}
		if p.Tags == nil {
			p.Tags = []string{}
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func getPost(db *sql.DB, w http.ResponseWriter, r *http.Request, id int64) {
	p, err := selectPostByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func updatePost(db *sql.DB, w http.ResponseWriter, r *http.Request, id int64) {
	var in PostInput
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Title) == "" ||
		strings.TrimSpace(in.Content) == "" ||
		strings.TrimSpace(in.Category) == "" {
		http.Error(w, "title/content/category required", http.StatusBadRequest)
		return
	}
	if in.Tags == nil {
		in.Tags = []string{}
	}
	tagsJSON, err := json.Marshal(in.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := db.Exec(
		`UPDATE posts
		 SET title=?, content=?, category=?, tags=?, updated_at=NOW()
		 WHERE id=?`,
		in.Title, in.Content, in.Category, tagsJSON, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if n, err := res.RowsAffected(); err == nil && n == 0 {
		if _, err := selectPostByID(db, id); err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	p, err := selectPostByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func deletePost(db *sql.DB, w http.ResponseWriter, r *http.Request, id int64) {
	res, err := db.Exec(`DELETE FROM posts WHERE id=?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {

	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		log.Fatal(err)
	}
	if c.DSN == "" {
		log.Fatal("put mysql_dsn in config.json")
	}
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		createPost(db, w, r)
	})
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		listPosts(db, w, r)
	})

	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		getPost(db, w, r, id)
	})

	mux.HandleFunc("PUT /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		updatePost(db, w, r, id)
	})

	mux.HandleFunc("DELETE /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		deletePost(db, w, r, id)
	})

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
