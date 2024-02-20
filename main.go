package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/htahreem/sample/database"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Name      string `json:"name"`
	RollNo    int    `json:"roll_no"`
	ContactNo int    `json:"contact_no"`
	Email     string `json:"email"`
	ID        string `json:"id"`
}

func DeleteStudent(ctx *gin.Context) {
	ID := ctx.Param("ID")

	tx, err := database.Db.Begin()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	_, err = database.Db.Exec("DELETE FROM students WHERE id = $1", ID)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't delete user.")
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully deleted.")
}

func GetAllStudents(ctx *gin.Context) {
	rows, err := database.Db.Query("SELECT name, roll_no, contact_no, email, id FROM students")
	// rows, err := database.Db.Query("SELECT * FROM students")
	if err != nil {
		fmt.Print(err)
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var stu Student
		err := rows.Scan(&stu.Name, &stu.RollNo, &stu.ContactNo, &stu.Email, &stu.ID)
		if err != nil {
			fmt.Print(err)
			ctx.AbortWithStatusJSON(500, "Internal Server Error")
			return
		}

		// Printing table values here
		// fmt.Printf("Raw values: %+v\n", stu)

		students = append(students, stu)
	}

	ctx.JSON(http.StatusOK, students)
}

func AddUser(ctx *gin.Context) {
	stu := Student{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &stu)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	_, err = database.Db.Exec(`INSERT INTO students (name, roll_no, contact_no, email, id) VALUES ($1, $2, $3, $4, $5)`, stu.Name, stu.RollNo, stu.ContactNo, stu.Email, stu.ID)
	// _, err = database.Db.Exec(`INSERT INTO students VALUES ($1, $2, $3, $4, $5)`, stu.Name, stu.RollNo, stu.ContactNo, stu.Email, stu.ID)

	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}

func UpdateUser(ctx *gin.Context) {
	ID := ctx.Param("ID")
	stu := Student{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User data is not defined")
		return
	}
	err = json.Unmarshal(data, &stu)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	tx, err := database.Db.Begin()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	_, err = tx.Exec(`
        UPDATE students
        SET name = $1,
            roll_no = $2,
            contact_no = $3,
            email = $4
        WHERE id = $5`,
		stu.Name, stu.RollNo, stu.ContactNo, stu.Email, ID)

	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't update user.")
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.AbortWithStatusJSON(500, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully updated.")
}

func main() {
	route := gin.Default()
	database.ConnectDatabase()
	database.RunMigration()

	route.GET("/getAllStudents", GetAllStudents)
	route.POST("/addStudent", AddUser)
	route.PUT("/updateStudent/:ID", UpdateUser)
	route.DELETE("/deleteStudent/:ID", DeleteStudent)

	// port := os.Getenv("PORT")
	port := os.Getenv("8001")
	if port == "" {
		port = "3005"
	}

	err := route.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
