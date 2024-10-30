package PatientCare

import (
	"bytes"
	"clinic_server/database"
	"net/http"

	"clinic_server/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GetPatientByParams(c *gin.Context) {
	var searchCondition SearchResult
	var searchResult []SearchResult

	logger.Init(zerolog.InfoLevel)

	err := c.ShouldBindJSON(&searchCondition)
	if err != nil {
		logger.Error("Invalid Search Condition", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if searchCondition.ID.Int64 != 0 {
		GetPatientById(c)
	} else {
		var query = getWhereClausenBasedOnSearch(searchCondition)
		var db *gorm.DB = database.GetDBContext()
		err = db.Raw(query).Scan(&searchResult).Error
		if err == nil {
			c.IndentedJSON(http.StatusOK, searchResult)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
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
	log.Info().Msg(sqlQuery.String())
	return sqlQuery.String()
}
