package gormfilter

import (
	"fmt"
	"reflect"

	"gorm.io/gorm/clause"
)

// Eq Equal
type Eq = clause.Eq

// Neq Not Equal
type Neq = clause.Neq

// Gt Greater than
type Gt = clause.Gt

// Lt Less than
type Lt = clause.Lt

// Gte Greater than or equal
type Gte = clause.Gte

// Lte Less than or equal
type Lte = clause.Lte

// Contains LIKE
type Contains clause.Eq

// Build ...
func (c Contains) Build(builder clause.Builder) {
	builder.WriteQuoted(c.Column)
	builder.WriteString(" LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", c.Value))
}

// NContains ...
type NContains clause.Eq

// Build ...
func (nc NContains) Build(builder clause.Builder) {
	builder.WriteQuoted(nc.Column)
	builder.WriteString(" NOT LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", nc.Value))
}

// IContains ILIKE
type IContains clause.Eq

// Build ...
func (ic IContains) Build(builder clause.Builder) {
	builder.WriteQuoted(ic.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", ic.Value))
}

// NIContains ILIKE
type NIContains clause.Eq

// Build ...
func (nic NIContains) Build(builder clause.Builder) {
	builder.WriteQuoted(nic.Column)
	builder.WriteString(" NOT ILIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", nic.Value))
}

// In ...
type In clause.Eq

// Build ...
func (in In) Build(builder clause.Builder) {
	var values []interface{}
	rv := reflect.ValueOf(in.Value)
	for i := 0; i < rv.Len(); i++ {
		values = append(values, rv.Index(i).Interface())
	}
	inClause := clause.IN{
		Column: in.Column,
		Values: values,
	}
	inClause.Build(builder)
}

// Nin ...
type Nin clause.Eq

// Build ...
func (nin Nin) Build(builder clause.Builder) {
	var values []interface{}
	rv := reflect.ValueOf(nin.Value)
	for i := 0; i < rv.Len(); i++ {
		values = append(values, rv.Index(i).Interface())
	}
	inClause := clause.IN{
		Column: nin.Column,
		Values: values,
	}
	inClause.NegationBuild(builder)
}
