package main

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-jwt/domain"
	"go-jwt/handler"
	"go-jwt/repository"
	"go-jwt/service"
	"log"
	"net/http"
	"time"
)

//go:embed migrations
var migrations embed.FS

func init() {
	viper.SetConfigFile(`../config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// jwt
	jwt := domain.JwtCredentials{
		ApplicationName: viper.GetString(`jwt.applicationName`),
		SecretKey:       viper.GetString(`jwt.secretKey`),
		Expired:         viper.GetInt(`jwt.expired`),
	}

	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.name`)
	dbPort := viper.GetString(`database.port`)
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbName, dbPort)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test ping
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Successfully ping db")

	// migration + seeeder
	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	// repos
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	// services
	productService := service.NewProductService(productRepo)
	userService := service.NewUserService(userRepo, jwt)

	// handlers
	timeout := viper.GetInt(`context.timeout`)
	portAddress := viper.GetString(`server.address`)
	timeOutContext := time.Duration(timeout) * time.Second
	productHandler := handler.NewProductHandler(productService, timeOutContext)
	userHandler := handler.NewUserHandler(userService, timeOutContext)

	// server
	router := httprouter.New()
	router.GET("/product", productHandler.GetProductByName)
	router.POST("/login", userHandler.Authentication)
	logrus.Infof("Server run on localhost%v", portAddress)
	log.Fatal(http.ListenAndServe(portAddress, router))
}
