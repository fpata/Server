package Patient

import (
	"bytes"
	"clinic_server/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPatientByParams(c *gin.Context) {
	var searchCondition SearchResult
	var searchResult []SearchResult
	err := c.ShouldBindJSON(&searchCondition)
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		c.Abort()
		return
	}
	if searchCondition.ID.Int64 != 0 {
		GetPatientById(c)
	} else {
		var query = getWhereClausenBasedOnSearch(searchCondition)
		var db *gorm.DB = database.GetDBContext()
		db.Raw(query).Scan(&searchResult)
		c.IndentedJSON(http.StatusOK, searchResult)
	}
}

func getWhereClausenBasedOnSearch(searchCondition SearchResult) string {
	var putAndCondition bool = false
	var sqlQuery bytes.Buffer
	sqlQuery.WriteString("Select Id,FirstName,LastName,PrimaryPhone,PrimaryEmail,PermCity from Patient Where ")
	if len(searchCondition.FirstName.String) != 0 {
		sqlQuery.WriteString("FirstName like '%")
		sqlQuery.WriteString(searchCondition.FirstName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.LastName.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" LastName like '%")
		sqlQuery.WriteString(searchCondition.LastName.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PrimaryEmail.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PrimaryEmail like '%")
		sqlQuery.WriteString(searchCondition.PrimaryEmail.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PrimaryPhone.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PrimaryPhone like '%")
		sqlQuery.WriteString(searchCondition.PrimaryPhone.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	if len(searchCondition.PermCity.String) != 0 {
		if putAndCondition {
			sqlQuery.WriteString(" And ")
		}
		sqlQuery.WriteString(" PermCity like '%")
		sqlQuery.WriteString(searchCondition.PermCity.String)
		sqlQuery.WriteString("%'")
		putAndCondition = true
	}
	fmt.Println(sqlQuery.String())
	return sqlQuery.String()
}
