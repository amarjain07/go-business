package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
  "database/sql"
	"math/rand"
	"crypto/md5"
  "encoding/hex"
	"time"
	// _ "github.com/heroku/x/hmetrics/onload"
)

type User struct{
	Name string `json:"name, omitempty"`
	Mobile string `json:"mobile"`
	Role string `json: role, omitempty`
}

type Product struct{
	Id int `json:"id", omitempty`
	Name string `json:"name"`
	Mrp int `json:"mrp"`
	Price int `json:"price"`
	RetailerPrice int `json:"retailer_price, omitempty"`
	Brand string `json:"brand"`
	Category string `json:"category"`
	Description string `json:"description, omitempty"`
	Image string `json:"image, omitempty"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	//router.Static("/static", "static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		panic(err)
	}

	router.GET("/link/:mobile", func(c *gin.Context){
		var mobile = c.Param("mobile")
		stmt, err := db.Prepare("insert into code(id, code) values(?, ?)")
		defer stmt.Close()
		if err != nil{
			c.JSON(http.StatusOK, gin.H{"success": true, "code": 4003, "message": "Error sending sms"})
			return
		}
		var random = 1000 + rand.Intn(8999)
		hash := getMD5Hash(mobile + time.Now().String())
		_, err = stmt.Exec(hash, random)
		if err != nil{
			c.JSON(http.StatusOK, gin.H{"success": true, "code": 4003, "message": "Error sending sms"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "code": 2000, "message": "SMS sent successfully!", "sms_code": random, "ref_id": hash})
	})

	router.GET("/verify/:id/:code", func(c *gin.Context){
		var hash = c.Param("id")
		paramCode, _ := strconv.Atoi(c.Param("code"))
		var code int
		row := db.QueryRow("select code from code where id = ?", hash)
		err = row.Scan(&code)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4004, "message": err.Error()})
			return
		}
		if paramCode == code{
			c.JSON(http.StatusOK, gin.H{"success": true, "code": 2000, "message": "User verified successfully!"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4005, "message": "Invalid otp code"})
		}
	})

	router.GET("/products", func(c *gin.Context){
		var (
			product  Product
			products []Product
		)
		rows, err := db.Query("select id, name, mrp, price, retailer_price, brand, category, description, image from product")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4002, "message": err.Error()})
			return
		}
		for rows.Next() {
			err = rows.Scan(&product.Id, &product.Name, &product.Mrp, &product.Price, &product.RetailerPrice, &product.Brand, &product.Category, &product.Description, &product.Image)
			if err != nil {
				continue
			}
			products = append(products, product)
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code": 2000,
			"products":  products,
		})
	})

	router.POST("/product", func(c *gin.Context){
		var product Product
		if err := c.ShouldBindJSON(&product); err == nil {
			log.Print("Name" + product.Name)
			stmt, err := db.Prepare("insert into product(name, mrp, price, retailer_price, brand, category, description, image) values(?, ?, ?, ?, ?, ?, ?, ?)")
			_, err = stmt.Exec(product.Name, product.Mrp, product.Price, product.RetailerPrice, product.Brand, product.Category, product.Description, product.Image)
			if(err != nil){
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4001, "message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "code": 2000, "message": "Product added successfully!"})
		}else{
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4001, "message": err.Error()})
		}
	})

	router.POST("/user", func(c *gin.Context){
		var (
			user User
			isLogin bool
		)
		queryIsLogin := c.Query("is_login")
		if(queryIsLogin == ""){
			isLogin = false
		}else{
			isLogin, _= strconv.ParseBool(queryIsLogin)
		}
		if err := c.ShouldBindJSON(&user); err == nil {
			tx, er := db.Begin()
			if er != nil {
					c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4000, "message": err.Error()})
					return
			}
			defer tx.Rollback()
			var roleId int
			row := tx.QueryRow("select id from role where type = ?", getRole(user.Role))
			er = row.Scan(&roleId)
			_, er = tx.Query("insert into user(name, mobile, is_active, role) values(?, ?, ?, ?) on duplicate key update name = ?, role = ?, is_active=true", user.Name, user.Mobile, isLogin, roleId, user.Name, roleId)
			tx.Commit()
			if er != nil{
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4000, "message": er.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "code": 2000, "message": "User added successfully!"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 4000, "message": err.Error()})
		}
	})

	router.Run(":" + port)
}

func getMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func getRole(role string) string{
	if role != ""{
		return role
	}
	return "CUSTOMER"
}
