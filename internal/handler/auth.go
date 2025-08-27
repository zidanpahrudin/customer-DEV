package handler

import (
	"net/http"
	"os"
	"time"

	"strconv"
	"fmt"
	"customer-api/internal/config"
	"customer-api/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	RoleID   string   `json:"role_id"` // Tambahkan field role_id (optional)
}

// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi bahwa salah satu dari username atau email harus diisi
	if input.Username == "" && input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau email harus diisi"})
		return
	}

	var user entity.User
	var usernameOrEmail string

	if input.Username != "" {
		usernameOrEmail = input.Username
	} else {
		usernameOrEmail = input.Email
	}

	
		// Always convert to string
		var uid string
		switch v := userID.(type) {
		case float64: // common case when ID was numeric
			uid = fmt.Sprintf("%.0f", v)
		case int:
			uid = strconv.Itoa(v)
		case int64:
			uid = strconv.FormatInt(v, 10)
		case uint:
			uid = strconv.FormatUint(uint64(v), 10)
		case string: // already string
			uid = v
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}


	// Cek apakah username/email terdaftar
	result := config.DB.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Username atau email tidak ditemukan
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau email tidak terdaftar"})
			return
		}
		// Error database lainnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat login"})
		return
	}

	// Cek password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// @Summary Register new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser entity.User
	if result := config.DB.Where("username = ?", input.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	// Check if email already exists
	if result := config.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Set default role_id jika tidak disediakan
	roleID := input.RoleID
	if roleID == "" { // ✅ string kosong, bukan 0
		// Cari role "User" dari database
		var userRole entity.Role
		if err := config.DB.Where("name = ?", "User").First(&userRole).Error; err != nil {
			// Jika role User tidak ditemukan, gunakan role pertama yang ada
			if err := config.DB.First(&userRole).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Tidak ada role yang tersedia"})
				return
			}
		}
		roleID = userRole.ID // ✅ sekarang string, bukan uint
	} else {
		// Validasi bahwa role_id yang diberikan ada di database
		var role entity.Role
		if err := config.DB.Where("id = ?", roleID).First(&role).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID tidak valid"})
			return
		}
	}

	// Buat user baru
	user := entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   roleID, // ✅ string ULID
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil didaftarkan"})
}

