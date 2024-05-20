package pharmacy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Сount struct {
	Сount int `sql:"count"`
}

func formNeedDrugsSubQuery(b *strings.Builder, needDrugs []int) error {
	if len(needDrugs) == 0 {
		b.WriteString("need_drug_ids as (select distinct * from (VALUES (-1)))")
		return nil
	}

	b.WriteString("need_drug_ids as (select distinct * from(VALUES ")
	for i, needDrug := range needDrugs {
		b.WriteString("(")
		b.WriteString(strconv.Itoa(needDrug))
		b.WriteString(")")
		if i != len(needDrugs)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString("))")

	return nil
}

func formPeriodSubQuery(b *strings.Builder, begin time.Time, end time.Time) error {

	b.WriteString("period as (select date ('")
	b.WriteString(begin.Format("2006-01-02"))
	b.WriteString("') \"start\", date ('")
	b.WriteString(end.Format("2006-01-02"))
	b.WriteString("') \"end\")")

	return nil
}

func (m *Model) formNeedTypesSubQuery(b *strings.Builder, needTypes []int) error {

	if len(needTypes) == 0 {
		b.WriteString("need_types as (select distinct * from (VALUES (-1)))")
		return nil
	}

	b.WriteString("need_types as (select distinct * from(VALUES ")
	for i, needType := range needTypes {
		if needType < 0 || needType >= len(standardTypes) {
			return fmt.Errorf("no type with id: %d", needType)
		}
		idSql := m.typesToId[standardTypes[needType]]
		b.WriteString("(")
		b.WriteString(strconv.Itoa(idSql))
		b.WriteString(")")
		if i != len(needTypes)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString("))")

	return nil
}
