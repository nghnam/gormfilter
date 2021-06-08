package gormfilter

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"
)


func TestExpression(t *testing.T) {
	db, _ := gorm.Open(tests.DummyDialector{}, nil)
	results := []struct {
		Clauses []clause.Interface
		Result string
		Vars []interface{}
	}{
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Or(Neq{Column: "name", Value: "nghnam"})},
			}},
			"WHERE `users`.`id` = ? OR `name` <> ?", []interface{}{"1", "nghnam"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{Gt{Column: "age", Value: 18}, Lt{Column: "age", Value: 40}},
			}},
			"WHERE `age` > ? AND `age` < ?", []interface{}{18, 40},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{Gte{Column: "age", Value: 18}, Lte{Column: "age", Value: 40}},
			}},
			"WHERE `age` >= ? AND `age` <= ?", []interface{}{18, 40},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{Contains{Column: "name", Value: "Nam"}},
			}},
			"WHERE `name` LIKE ?", []interface{}{"%Nam%"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{IContains{Column: "name", Value: "Nam"}},
			}},
			"WHERE `name` ILIKE ?", []interface{}{"%Nam%"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{NContains{Column: "name", Value: "Nam"}},
			}},
			"WHERE `name` NOT LIKE ?", []interface{}{"%Nam%"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{NIContains{Column: "name", Value: "Nam"}},
			}},
			"WHERE `name` NOT ILIKE ?", []interface{}{"%Nam%"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{In{Column: "name", Value: []string{"Nguyen", "Hoang", "Nam"}}},
			}},
			"WHERE `name` IN (?,?,?)", []interface{}{"Nguyen", "Hoang", "Nam"},
		},
		{
			[]clause.Interface{clause.Where{
				Exprs: []clause.Expression{Nin{Column: "name", Value: []string{"Nguyen", "Hoang", "Nam"}}},
			}},
			"WHERE `name` NOT IN (?,?,?)", []interface{}{"Nguyen", "Hoang", "Nam"},
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			var (
				buildNames    []string
				buildNamesMap = map[string]bool{}
				user, _       = schema.Parse(&tests.User{}, &sync.Map{}, db.NamingStrategy)
				stmt          = gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}
			)
			for _, c := range result.Clauses {
				if _, ok := buildNamesMap[c.Name()]; !ok {
					buildNames = append(buildNames, c.Name())
					buildNamesMap[c.Name()] = true
				}

				stmt.AddClause(c)
			}
			stmt.Build(buildNames...)
			if strings.TrimSpace(stmt.SQL.String()) != result.Result {
				t.Errorf("SQL expects %v got %v", result.Result, stmt.SQL.String())
			}

			if !reflect.DeepEqual(stmt.Vars, result.Vars) {
				t.Errorf("Vars expects %+v got %v", stmt.Vars, result.Vars)
			}
		})
	}
}