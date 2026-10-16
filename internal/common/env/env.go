package env

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	IsProduction bool
	Port         string
	Host         string
	DatabaseUrl  string

	// access token
	ExpiresAtAccessToken time.Duration
	SecretAccessToken    string
	// refresh token
	ExpiresAtRefreshToken time.Duration
	SecretRefreshToken    string

	// Google Login feature
	GoogleClientId     string
	GoogleClientSecret string
	GoogleRedirectUrl  string

	DomainFe      string
	DomainBe      string
	CloudinaryURL string

	RedisAddr string
	RedisPass string

	// elastic search
	ElasticAddrs           string
	ElasticUser            string
	ElasticPassword        string
	ElasticCertFingerprint string

	RabbitMQURL string
}

func New() *Env {
	// load biến rồi đưa vào object
	godotenv.Load()

	// lấy ra
	isProduction := os.Getenv("IS_PRODUCTION") == "true"
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	databaseUrl := os.Getenv("DATABASE_URL")

	// Get access Token
	expiresAtAccessTokenString := os.Getenv("EXPIRES_AT_ACCESS_TOKEN")
	expiresAtAccessToken := getDuration(expiresAtAccessTokenString)
	secretAccessToken := os.Getenv("SECRET_ACCESS_TOKEN")

	// Get refresh token
	expiresAtRefreshTokenString := os.Getenv("EXPIRES_AT_REFRESH_TOKEN")
	expiresAtRefreshToken := getDuration(expiresAtRefreshTokenString)
	secretRefreshToken := os.Getenv("SECRET_REFRESH_TOKEN")

	// Google Login
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")

	domainFe := os.Getenv("DOMAIN_FE")
	domainBe := os.Getenv("DOMAIN_BE")

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")

	// redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")

	// Elastic search
	elasticAddrs := os.Getenv("ELASTIC_ADDRS")
	elasticUser := os.Getenv("ELASTIC_USER")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	elasticCertFingerprint := os.Getenv("ELASTIC_CERT_FINGERPRINT")

	// Rabbit MQ
	rabbitMQURL := os.Getenv("RABBIT_MQ_URL")

	// server infor
	fmt.Println("port", port)
	fmt.Println("isProduction", isProduction)
	fmt.Println("host", host)
	fmt.Println("databaseUrl", databaseUrl)

	// print access token
	fmt.Println("expiresAtAccessTokenString", expiresAtAccessTokenString)
	fmt.Println("secretAccessToken", secretAccessToken)

	// print refresh token
	fmt.Println("expiresAtRefreshTokenString", expiresAtRefreshTokenString)
	fmt.Println("secretRefreshToken", secretRefreshToken)

	// print ra trong lúc dev để dễ follow
	fmt.Println("googleClientId", googleClientId)
	fmt.Println("googleClientSecret", googleClientSecret)
	fmt.Println("googleRedirectUrl", googleRedirectUrl)

	fmt.Println("domainFe", domainFe)
	fmt.Println("cloudinaryUrl", cloudinaryURL)

	fmt.Println("RedisAddr", redisAddr)
	fmt.Println("RedisPass", redisPass)

	// elastic search
	fmt.Println("elasticAddrs", elasticAddrs)
	fmt.Println("elasticUser", elasticUser)
	fmt.Println("elasticPassword", elasticPassword)
	fmt.Println("elasticCertFingerprint", elasticCertFingerprint)

	fmt.Println("rabbitMQURL", rabbitMQURL)

	return &Env{
		// config db
		IsProduction: isProduction,
		Port:         port,
		Host:         host,
		DatabaseUrl:  databaseUrl,

		// Refresh + Access Token
		ExpiresAtAccessToken:  expiresAtAccessToken,
		SecretAccessToken:     secretAccessToken,
		ExpiresAtRefreshToken: expiresAtRefreshToken,
		SecretRefreshToken:    secretRefreshToken,

		// Google login
		GoogleClientId:     googleClientId,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectUrl:  googleRedirectUrl,

		DomainFe: domainFe,
		DomainBe: domainBe,

		CloudinaryURL: cloudinaryURL,

		// Redis
		RedisAddr: redisAddr,
		RedisPass: redisPass,

		// elastic
		ElasticAddrs:           elasticAddrs,
		ElasticUser:            elasticUser,
		ElasticPassword:        elasticPassword,
		ElasticCertFingerprint: elasticCertFingerprint,
		RabbitMQURL:            rabbitMQURL,
	}
}

// hàm duration
func getDuration(durationString string) time.Duration {

	durationTime, err := time.ParseDuration(durationString)
	if err != nil {
		log.Fatal("Parser expiresAtAccessToken error")
	}

	return durationTime
}
