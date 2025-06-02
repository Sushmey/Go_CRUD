package main


import(
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

type DBClient interface{
	Ready() bool
	RunMigration() error
}

type Client struct{ // this is a concrete implementation of the  above interface
	db *gorm.DB //Interprets the value of the gorm DB which points are gorm.DB instance
}


// c Client is a receiver for the function Ready, meaning its scope is only the elements of the struct Client
func (c Client) Ready() bool{
	var ready string
	result := c.db.Raw("SELECT 1 as ready").Scan(&ready)

	if result.Error != nil{
		return false
	}
	if ready=="1"{
		return true
	}

	return false
}

func (c Client) RunMigration() error{
	if(!c.Ready()){
		log.Fatal("Database is not ready") // DB is not ready for migration
	}
	
	err := c.db.AutoMigrate(&User{})
	if err != nil{
		return err
	}
	return nil
}


func NewDBClient() (Client, error){
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	// dbHost := "localhost"
	// dbUsername := "postgres"
	// dbPassword := "!"
	// dbName := "postgres"
	// dbPort := "5432"
	databasePort, err := strconv.Atoi(dbPort) //Convert int to string
	if err != nil{
		// log.Fatal("Invalid port")
		log.Fatalf("databasePort: %v Error: %s/n", databasePort, err)
	}

	fmt.Printf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",dbHost,dbUsername,dbPassword,dbName,databasePort,"disable") //Create a connection string

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",dbHost,dbUsername,dbPassword,dbName,databasePort,"disable") //Create a connection string

	db, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})

	if err!=nil{
		return Client{}, err
	}

	client := Client{db}
	return client, nil


}

