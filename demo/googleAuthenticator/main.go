package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// user account password
var username = "admin"
var password = "8lab123"

// Authenticator TOTP 
var userTOTPSecret string

func initTOTP() string {
    secretKey, err := totp.Generate(totp.GenerateOpts{
        Issuer: "myApp",
        AccountName: username,
    })
    if err != nil { 
        fmt.Println(err)
    }

    return secretKey.Secret()
}

// generater QRCode
func getQRCode(c *gin.Context) {
    // qrCodeUrl := totp.URL(userTOTPSecret, "MyApp", username)
    username := "admin"
    secretKey := initTOTP()

    qrCodeUrl := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issure=%s", "MyApp", username, secretKey, "MyApp")

    c.JSON(http.StatusOK, gin.H{
        "qr_code_url": qrCodeUrl,
        "qu_key": secretKey,
        "message": "Scan the QR code or add the key manually in google Authenticator",

    })
    return
}

func login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
        TOTPCode string `json:"toty_code"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

        return
    }

    // check user password
    if loginData.Username != username || loginData.Password != password {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
        return
    }

    if valida := totp.Validate(loginData.TOTPCode, userTOTPSecret); !valida {
        c.JSON(http.StatusOK, gin.H{"message": "login successful"})
        return
    }

}


func main() {
    // init TOTP
    userTOTPSecret = initTOTP()
    fmt.Println("TOTP Secret:", userTOTPSecret)

    r := gin.Default()
    r.POST("/login", login)

    r.GET("/qr_code", getQRCode)

    r.Run("10.200.1.82:9000")


}
