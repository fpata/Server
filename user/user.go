package user

import (
	"clinic_server/database"
	"clinic_server/types"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	Id                       types.NullInt64  `gorm:"PrimaryKey"`
	FirstName                types.NullString `json:"FirstName"`
	LastName                 types.NullString `json:"LastName"`
	MiddleName               types.NullString `json:"MiddleName"`
	Age                      types.NullInt64  `json:"Age"`
	Gender                   types.NullInt64  `json:"Gender"`
	UserName                 types.NullString `json:"UserName"`
	Password                 types.NullString `json:"Password"`
	PermAddress1             types.NullString `json:"PermAddress1"`
	PermAddress2             types.NullString `json:"PermAddress2"`
	PermCity                 types.NullString `json:"PermCity"`
	PermState                types.NullString `json:"PermState"`
	PermCountry              types.NullString `json:"PermCountry"`
	PermPostalCode           types.NullString `json:"PermPostalCode"`
	CorrAddress1             types.NullString `json:"CorrAddress1"`
	CorrAddress2             types.NullString `json:"CorrAddress2"`
	CorrCity                 types.NullString `json:"CorrCity"`
	CorrState                types.NullString `json:"CorrState"`
	CorrCountry              types.NullString `json:"CorrCountry"`
	CorrPostalCode           types.NullString `json:"CorrPostalCode"`
	PrimaryPhone             types.NullString `json:"PrimaryPhone"`
	PrimaryEmail             types.NullString `json:"PrimaryEmail"`
	SecondaryPhone           types.NullString `json:"SecondaryPhone"`
	SecondaryEmail           types.NullString `json:"SecondaryEmail"`
	EmergencyContactName     types.NullString `json:"EmergencyContactName"`
	EmergencyContactEmail    types.NullString `json:"EmergencyContactEmail"`
	EmergencyContactPhone    types.NullString `json:"EmergencyContactPhone"`
	EmergencyContactRelation types.NullString `json:"EmergencyContactRelation"`   
	ExistingDieases 		 types.NullString `json:"ExistingDieases"`
	Medications 		 	 types.NullString `json:"Medications"`
	Allergies 		 		 types.NullString `json:"Allergies"`
	FatherMedicalHistory 	 types.NullString `json:"FatherMedicalHistory"`
	MotherMedicalHistory 	 types.NullString `json:"MotherMedicalHistory"`
}

type SearchResult struct {
	Id                       types.NullInt64  `gorm:"PrimaryKey"`
	FirstName                types.NullString `json:"FirstName"`
	LastName                 types.NullString `json:"LastName"`
	PrimaryPhone             types.NullString `json:"PrimaryPhone"`
	PrimaryEmail             types.NullString `json:"PrimaryEmail"`
	PermCity                 types.NullString `json:"PermCity"`
}
var users = []User{}

func GetAllUsersWithPaging(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	var db *gorm.DB = database.GetDBContext()
	db.Table("user").Limit(limit).Offset(offset).Find(&users)
	c.IndentedJSON(http.StatusOK, users)
}

func GetAllUsers(c *gin.Context) {
	var db *gorm.DB = database.GetDBContext()
	db.Table("User").Find(&users)
	c.IndentedJSON(http.StatusOK, users)
}

func GetUserByParams(c *gin.Context) {
	var searchCondition SearchResult
	var searchResult []SearchResult
	if err := c.ShouldBindJSON(&searchCondition); err != nil {
        fmt.Println(err);
		fmt.Println(c.Error)
		c.Error(err)

        c.Abort()
        return
    }
	var query = getWhereClausenBasedOnSearch(searchCondition)
	var db *gorm.DB = database.GetDBContext()
	db.Raw(query).Scan(&searchResult)
	c.IndentedJSON(http.StatusOK, searchResult)
}

func getWhereClausenBasedOnSearch(searchCondition SearchResult) string {
	var putAndCondition bool=false
	var sqlQuery bytes.Buffer
	sqlQuery.WriteString("Select Id,FirstName,LastName,PrimaryPhone,PrimaryEmail,PermCity from user Where ")
	if(len(searchCondition.FirstName.String) != 0){
		sqlQuery.WriteString("FirstName like '%")
		sqlQuery.WriteString(searchCondition.FirstName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if(len(searchCondition.LastName.String) != 0){
		if(putAndCondition) { sqlQuery.WriteString(" And ") }
		sqlQuery.WriteString(" LastName like '%")
		sqlQuery.WriteString(searchCondition.LastName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if(len(searchCondition.PrimaryEmail.String) != 0){
		if(putAndCondition) { sqlQuery.WriteString(" And ") }
		sqlQuery.WriteString(" PrimaryEmail like '%")
		sqlQuery.WriteString(searchCondition.PrimaryEmail.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if(len(searchCondition.PrimaryPhone.String) != 0){
		if(putAndCondition) {sqlQuery.WriteString(" And ")}
		sqlQuery.WriteString(" PrimaryPhone like '%")
		sqlQuery.WriteString(searchCondition.PrimaryPhone.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if(len(searchCondition.PermCity.String) != 0){
		if(putAndCondition) {sqlQuery.WriteString(" And ")}
		sqlQuery.WriteString(" PermCity like '%")
		sqlQuery.WriteString(searchCondition.PermCity.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	fmt.Println(sqlQuery.String())
return sqlQuery.String()
}

func GetUserById(c *gin.Context) {
	userId := c.Param("Id")
	var objUser User = User{}
	var db *gorm.DB = database.GetDBContext()
	db.Table("user").First(&objUser, userId)
	c.IndentedJSON(http.StatusOK, objUser)
}

func CreateUser(c *gin.Context) {
	var objUser User
	c.ShouldBind(&objUser)
	var db *gorm.DB = database.GetDBContext()
	createResult := db.Create(&objUser)

	if createResult.Error != nil {
		fmt.Println(createResult.Error)
	} else {
		fmt.Println(createResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, objUser)
}

func UpdateUser(c *gin.Context) {
	var objUser User
	c.ShouldBind(&objUser)
	var db *gorm.DB = database.GetDBContext()
	updateResult := db.Save(&objUser)
	if updateResult.Error != nil {
		fmt.Println(updateResult.Error)
	} else {
		fmt.Println(updateResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, objUser)
}

func PatchUser(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}

	var db *gorm.DB = database.GetDBContext()
	updateResult := db.Updates(&jsonData)
	if updateResult.Error != nil {
		fmt.Println(updateResult.Error)
	} else {
		fmt.Println(updateResult.RowsAffected)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	userId := c.Query("Id")
	var db *gorm.DB = database.GetDBContext()
	db.Delete(&User{}, userId)
	c.IndentedJSON(http.StatusOK, users)
}
