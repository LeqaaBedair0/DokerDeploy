package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte
var db *sql.DB

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mysecret"
	}
	jwtKey = []byte(secret)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func connectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Println("Connected to DB ✅")
			return
		}
		fmt.Println("Waiting for DB...")
		time.Sleep(3 * time.Second)
	}

	panic("Could not connect to DB")
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE,
		password VARCHAR(255)
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Fixed: using jwtKey instead of jwtSecret
	return token.SignedString(jwtKey)
}

func main() {
	connectDB()
	createTable()

	r := gin.Default()

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Auth Service Running"})
	})

	// Register
	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		_, err := db.Exec(
			"INSERT INTO users (username, password) VALUES (?, ?)",
			user.Username,
			string(hashed),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or DB error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	// Login
	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var storedPassword string
		err := db.QueryRow(
			"SELECT password FROM users WHERE username = ?",
			user.Username,
		).Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Compare hashed password
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, _ := generateToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	r.Run(":8000")
}
