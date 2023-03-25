package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var PORT = "8080"

const MySecret string = "bdf&1*~#^2^#s0^=)^^7%b59"
const validationTime = time.Minute * 30

// externalIP determines the external IP address
func ExternalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "localhost"
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	return "localhost" //errors.New("Are you connected to the network?")
}

func GetPort() (i int) {
	lPort := os.Getenv("PORT")
	if lPort != "" {
		i, _ = strconv.Atoi(lPort)
	} else {
		i, _ = strconv.Atoi(PORT)
	}
	return
}

// *****************************************************************//
type TToken struct {
	Random    string `json:"random"`
	Timestamp string `json:"timestamp"`
}

func (t *TToken) Code() string {
	lJsonToken, _ := json.Marshal(t)
	lCodeToken, _ := Encrypt(string(lJsonToken), MySecret)
	return lCodeToken
}

func (t *TToken) Decode(iCookie string) {
	lJsonToken, _ := Decrypt(iCookie, MySecret)
	json.Unmarshal([]byte(lJsonToken), &t)
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text, MySecret string) (string, error) {
	var _bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, _bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text, MySecret string) (string, error) {
	var _bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, _bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// ******************************************************************//
func main() {

	godotenv.Load(".env")

	// Creates a router without any middleware by default
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	//r.StaticFS("/", http.Dir("dist/spa")) //via NoRoute
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)

		//lToken := TToken{
		//	Random:    uuid.New().String(),
		//	Timestamp: time.Now().UTC().Format(time.RFC3339),
		//}

		c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		//c.SetCookie("movie_cat", lToken.Code(), 3600, "", "", false, true)
		if file == "" || ext == "" {
			c.File("./dist/pwa/index.html")
		} else {
			c.File("./dist/pwa/" + path.Join(dir, file))
		}
	})

	fmt.Println("server ip is ", ExternalIP())
	//router.Run("0.0.0.0:8080")
	if strings.HasPrefix(ExternalIP(), "192.168") {
		//wenn local, dann mit TLS
		panic(r.RunTLS(ExternalIP()+":"+strconv.Itoa(GetPort()), "./key/server.pem", "./key/server.key").Error())
	} else {
		//wenn globale IP ohne TLS
		panic(r.Run(ExternalIP() + ":" + strconv.Itoa(GetPort())).Error())
	}

}
