package main

import(
  "os"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
)

func main(){
  db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
  if err != nil {
    panic(err)
  }
  defer db.Close()
  if _, err := db.Exec("CREATE TABLE IF NOT EXISTS role(id INTEGER NOT NULL AUTO_INCREMENT, type VARCHAR(20), PRIMARY KEY(`id`))"); err != nil{
    panic(err)
  }
  if _, err := db.Exec("CREATE TABLE IF NOT EXISTS user(id INT(10) NOT NULL AUTO_INCREMENT, name VARCHAR(64), mobile VARCHAR(15),is_active BOOL DEFAULT false,role INTEGER DEFAULT 4, PRIMARY KEY (`id`), UNIQUE(`mobile`), FOREIGN KEY(`role`) REFERENCES role(`id`))"); err != nil {
    panic(err)
  }
  if _, err := db.Exec("CREATE TABLE IF NOT EXISTS usermeta(access_token VARCHAR(128), user_id INT(10) UNIQUE, created_at DATETIME, FOREIGN KEY(`user_id`) REFERENCES user(`id`))"); err != nil {
    panic(err)
  }
  if _, err := db.Exec("CREATE TABLE IF NOT EXISTS product(id INT(10) NOT NULL AUTO_INCREMENT, name VARCHAR(64), mrp INTEGER, price INTEGER, retailer_price INTEGER, brand VARCHAR(20), category VARCHAR(20), description VARCHAR(64) NULL, image VARCHAR(512), PRIMARY KEY(`id`), UNIQUE(`name`))"); err != nil{
    panic(err)
  }
  if _, err := db.Exec("CREATE TABLE IF NOT EXISTS code(id VARCHAR(128) NOT NULL, code INTEGER, PRIMARY KEY(`id`))"); err != nil{
    panic(err)
  }
  stmt, err := db.Prepare("INSERT role SET id=?, type=?")
  if err != nil{
    panic(err)
  }
  if _, err := stmt.Exec(1, "OWNER"); err != nil{
    panic(err)
  }
  stmt, err = db.Prepare("INSERT role SET id=?, type=?")
  if err != nil{
    panic(err)
  }
  if _, err := stmt.Exec(2, "DISTRIBUTOR"); err != nil{
    panic(err)
  }
  stmt, err = db.Prepare("INSERT role SET id=?, type=?")
  if err != nil{
    panic(err)
  }
  if _, err := stmt.Exec(3, "RETAILER"); err != nil{
    panic(err)
  }
  stmt, err = db.Prepare("INSERT role SET id=?, type=?")
  if err != nil{
    panic(err)
  }
  if _, err := stmt.Exec(4, "CUSTOMER"); err != nil{
    panic(err)
  }
  stmt, err = db.Prepare("INSERT user SET mobile=?, role=?")
  if err != nil{
    panic(err)
  }
  if _, err := stmt.Exec("9845375411", 1); err != nil{
    panic(err)
  }
}
