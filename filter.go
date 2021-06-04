package gormfilter

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// FilterQuery ...
func FilterQuery(db *gorm.DB, s interface{}) (*gorm.DB, error) {
	joinTables, exprs, err := parseFilterStruct(s)
	if err != nil {
		return nil, err
	}
	tx := db
	for _, joinTable := range joinTables {
		tx = tx.Joins(joinTable)
	}
	tx = tx.Clauses(exprs...)
	return tx, nil
}

func parseFilterStruct(s interface{}) ([]string, []clause.Expression, error) {
	joinTables := []string{}
	exprs := []clause.Expression{}
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		rvalue := val.Field(i)

		if rvalue.Kind() == reflect.Ptr {
			rvalue = reflect.Indirect(rvalue)
		}
		if rvalue.Kind() == reflect.Slice && rvalue.Len() == 0 {
			continue
		}
		if !rvalue.IsValid() {
			continue
		}
		value := rvalue.Interface()
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if filter, ok := tag.Lookup("filter"); ok {
			setting := schema.ParseTagSetting(filter, ";")
			column := setting["COLUMN"]
			joinTable := getJoinTable(column)
			if joinTable != "" {
				joinTables = append(joinTables, joinTable)
			}
			expr, err := makeFilterClause(setting["COLUMN"], setting["OP"], value)
			if err != nil {
				return nil, nil, err
			}
			exprs = append(exprs, expr)
		}
	}
	return joinTables, exprs, nil
}

func makeFilterClause(column string, op string, value interface{}) (clause.Expression, error) {
	switch op {
	case "eq":
		return Eq{Column: column, Value: value}, nil
	case "neq":
		return Neq{Column: column, Value: value}, nil
	case "gt":
		return Gt{Column: column, Value: value}, nil
	case "gte":
		return Gte{Column: column, Value: value}, nil
	case "lt":
		return Lt{Column: column, Value: value}, nil
	case "lte":
		return Lte{Column: column, Value: value}, nil
	case "in":
		return In{Column: column, Value: value}, nil
	case "!in":
		return Nin{Column: column, Value: value}, nil
	case "like", "contains":
		return Contains{Column: column, Value: value}, nil
	case "!like", "!contains":
		return NContains{Column: column, Value: value}, nil
	case "ilike", "icontains":
		return IContains{Column: column, Value: value}, nil
	case "!ilike", "!icontains":
		return NIContains{Column: column, Value: value}, nil
	}
	return nil, fmt.Errorf(fmt.Sprintf("%v op is not supported", op))
}

func getJoinTable(column string) string {
	parts := strings.Split(column, ".")
	if len(parts) == 1 {
		return ""
	}
	return parts[0]
}
