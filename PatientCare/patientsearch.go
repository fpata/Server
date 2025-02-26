package PatientCare

import (
	"clinic_server/database"
	"net/http"

	"clinic_server/logger"

	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetPatientByParams(c *gin.Context) {
	var searchParams SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid search parameters",
			"details": err.Error(),
		})
		return
	}

	// Handle direct ID lookup
	if searchParams.ID.Valid && searchParams.ID.Int64 != 0 {
		GetPatientById(c)
		return
	}

	// Build query using query builder
	query, args := buildSearchQuery(searchParams)

	db := database.GetDBContext()
	if db == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	var searchResult []SearchResult
	if err := db.Raw(query, args...).Scan(&searchResult).Error; err != nil {
		logger.Error("Database query failed", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patient data"})
		return
	}

	if len(searchResult) == 0 {
		c.JSON(http.StatusOK, []SearchResult{})
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

// buildSearchQuery constructs the SQL query and arguments safely
func buildSearchQuery(params SearchParams) (string, []interface{}) {
	const baseQuery = `
		SELECT 
			Id,
			FirstName,
			LastName,
			PrimaryPhone,
			PrimaryEmail,
			PermCity 
		FROM Patient 
		WHERE 1=1`

	var (
		conditions []string
		args       []interface{}
		query      strings.Builder
	)

	query.WriteString(baseQuery)

	// Map of field conditions
	fieldConditions := map[*sql.NullString]string{
		&params.FirstName:    "FirstName LIKE ?",
		&params.LastName:     "LastName LIKE ?",
		&params.PrimaryEmail: "PrimaryEmail LIKE ?",
		&params.PrimaryPhone: "PrimaryPhone LIKE ?",
		&params.PermCity:     "PermCity LIKE ?",
	}

	// Build conditions and args
	for field, condition := range fieldConditions {
		if field.Valid && field.String != "" {
			conditions = append(conditions, condition)
			args = append(args, "%"+field.String+"%")
		}
	}

	// Add conditions to query
	for _, condition := range conditions {
		query.WriteString(" AND ")
		query.WriteString(condition)
	}

	// Add order by for consistent results
	query.WriteString(" ORDER BY Id")

	// Add limit to prevent excessive results
	query.WriteString(" LIMIT 100")

	return query.String(), args
}

/*func getWhereClausenBasedOnSearch(searchCondition SearchResult) string {
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
}*/
