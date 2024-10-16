package main

import (
    "bytes"
    "encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
    "github.com/skip2/go-qrcode"
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

    // 创建一个字节缓冲区
	var buffer bytes.Buffer

	// 生成二维码并写入字节缓冲区
	qrCode, err := qrcode.New(qrCodeUrl, qrcode.Medium)
	if err != nil {
		fmt.Println("生成二维码失败:", err)
		return
	}

	// 将二维码写入 buffer（PNG 格式）
	err = qrCode.Write(256, &buffer)
	if err != nil {
		fmt.Println("二维码写入缓冲区失败:", err)
		return
	}

	// Base64 编码
	qrCodeBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())



	// Base64 编码
    c.JSON(http.StatusOK, gin.H{
        "qr_code_url": qrCodeUrl,
        "qu_key": secretKey,
        "base64Qr": qrCodeBase64,
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
