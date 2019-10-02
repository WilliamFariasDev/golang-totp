package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/png"
	"os"

	"bitbucket.org/qg2/qg2-payment-totp/lib"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xlzd/gotp"
)

//MySQL a
var MySQL *gorm.DB

func main() {

	openConnection()

	user := defaultLogin()

	if user.TwoFactor == 1 {
		defaultTOTPUsage(user.UUID)
	}
	panic("login success")

}

func openConnection() {

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"user",
		"pass",
		"localhost",
		"3306",
		"payments")

	conn, err := gorm.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}

	MySQL = conn
}

func defaultLogin() lib.User {

	var email, pass string
	var item lib.User

	fmt.Println("Insira seu e-mail")
	fmt.Scanf("%s", &email)

	fmt.Println("Insira sua senha")
	fmt.Scanf("%s", &pass)

	pass = GetMD5Hash(pass)
	fmt.Println("hash md5 pass", pass)

	// query
	query := MySQL
	query = query.Table("payments_users.users")
	query = query.Where("payments_users.users.email = ?", email)
	query = query.Where("payments_users.users.password = ?", pass)

	if e := query.First(&item).Error; e != nil {
		panic(e.Error())
	}

	return item

}

func defaultTOTPUsage(uuid string) {

	var twoFactor lib.TwoFactor

	query := MySQL
	query = query.Select("payments_users.two_factor.key_secret")
	query = query.Table("payments_users.two_factor")
	query = query.Joins("JOIN payments_users.users AS us ON payments_users.two_factor.user_id = us.id")
	query = query.Where("us.uuid = ?", uuid)

	if e := query.First(&twoFactor).Error; e != nil {
		panic(e.Error())
	}

	otp := gotp.NewDefaultTOTP(twoFactor.KeySecret)

	fmt.Println("Digito o token")
	var passcode string
	fmt.Scanf("%s", &passcode)

	valid := otp.Now() == passcode
	if valid {
		println("Valid passcode!")
		// os.Exit(0)
	} else {
		println("Invalid passocde!")
		os.Exit(1)
	}

}

//GetMD5Hash as
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func register() {

	keySecret := gotp.RandomSecret(16)
	otp := gotp.NewDefaultTOTP(keySecret)

	qrcodeURL := otp.ProvisioningUri("email@servidor.com.br", "Servidor")

	qrCode, _ := qr.Encode(qrcodeURL, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	file, _ := os.Create("qr2.png")

	defer file.Close()
	png.Encode(file, qrCode)

	fmt.Println("Digito o token")
	var passcode string
	fmt.Scanf("%s", &passcode)

	valid := otp.Now() == passcode
	if valid {
		println("Valid passcode! Save in database")
		os.Exit(0)
	} else {
		println("Invalid passocde!")
		os.Exit(1)
	}

}
