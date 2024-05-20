package pharmacy

import (
	"context"
	"fmt"
	"pharmacy/interanl/pharmacy/daos"
	"strings"
	"time"
)

func (m *Model) GetLateConsumers(ctx context.Context) (map[daos.ID]daos.Consumer, error) {
	rv := make(map[daos.ID]daos.Consumer)
	err := daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		"with cur_date as (select cast(current_date as DATE) date)\n"+
			"select distinct c.id, c.name_, c.phone, c.addres, c.birth_date\n"+
			"from consumer_order co\n"+
			"         cross join cur_date\n"+
			"         inner join recipe r on recipe_id = r.id\n"+
			"         inner join consumers c on consumer_id = c.id\n"+
			"where complete_time < cur_date.date\n"+
			"  and co.status = 'WAITING';",
		&rv,
	)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (m *Model) GetWaitingConsumers(ctx context.Context, needTypes []int) (map[daos.ID]daos.Consumer, error) {

	var b strings.Builder

	b.WriteString("with ")

	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}
	b.WriteString("select distinct c.id, c.name_, c.phone, c.addres, c.birth_date\n" +
		"from consumer_order co\n" +
		"         inner join recipe r on co.recipe_id = r.id\n" +
		"         inner join consumers c on r.consumer_id = c.id\n" +
		"         inner join drug d on r.drug_id = d.id\n" +
		"where co.status = 'WAITING-COMPS'\n" +
		"  and d.type_id in (select * from need_types);",
	)

	rv := make(map[daos.ID]daos.Consumer)

	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

type DrugsWithCount struct {
	Name          string  `sql:"name_"`
	TypeId        daos.ID `sql:"type_id"`
	CriticalCount int     `sql:"critical_count"`
	Price         int     `sql:"price"`
	Count         int     `sql:"count"`
}

func (m *Model) GetTenMostPopularDrugs(ctx context.Context, needTypes []int) (map[daos.ID]DrugsWithCount, error) {
	var b strings.Builder

	b.WriteString("with ")

	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n" +
		"     help_table as (select comp_id, count(tech_id)\n" +
		"                    from technology_components tc\n" +
		"                             inner join drug d on tc.comp_id = d.id\n" +
		"                    where d.type_id in (select * from need_types)\n" +
		"                    group by comp_id)\n" +
		"select d.id, d.name_, d.type_id, d.critical_count, d.price, coalesce(ht.count, 0) count\n" +
		"from drug d\n" +
		"         left join help_table ht on ht.comp_id = d.id\n" +
		"where d.type_id in (select * from need_types)\n" +
		"order by count DESC\n" +
		"LIMIT 10;",
	)
	rv := make(map[daos.ID]DrugsWithCount)

	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetUsedAmmount(ctx context.Context, needDrugs []daos.ID, begin time.Time, end time.Time) (map[daos.ID]int, error) {
	var b strings.Builder

	b.WriteString("with ")
	err := formNeedDrugsSubQuery(&b, needDrugs)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n")

	err = formPeriodSubQuery(&b, begin, end)
	if err != nil {
		return nil, err
	}

	b.WriteString(
		",\n" +
			"     need_drugs as (select id\n" +
			"                    from drug d\n" +
			"                    where id in (select * from need_drug_ids)),\n" +
			"     drug_count as (select drug_id, count (*) count\n" +
			"                    from inventasrization\n" +
			"                        cross join period p\n" +
			"                    where p.start <= date_::date\n" +
			"                      and date_::date <= p.\"end\"\n" +
			"                      and drug_id in (select * from need_drug_ids)\n" +
			"                    group by drug_id)\n" +
			"select id, coalesce(dc.count, 0) count\n" +
			"from need_drugs\n" +
			"    left join drug_count dc\n" +
			"on id = dc.drug_id;",
	)

	myCount := make(map[int]Сount)

	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&myCount,
	)
	if err != nil {
		return nil, err
	}

	rv := make(map[daos.ID]int)

	for id, c := range myCount {
		rv[id] = c.Сount
	}

	return rv, nil
}

func (m *Model) GetConsumersOfDrugs(ctx context.Context, needDrugs []daos.ID, needTypes []int, begin time.Time, end time.Time) (map[daos.ID]daos.Consumer, error) {
	var b strings.Builder
	b.WriteString("with ")
	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}

	b.WriteString(",\n")
	err = formPeriodSubQuery(&b, begin, end)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n")
	err = formNeedDrugsSubQuery(&b, needDrugs)
	if err != nil {
		return nil, err
	}
	b.WriteString("select distinct c.id, c.name_, c.phone, c.addres, c.birth_date\n" +
		"from consumer_order co\n" +
		"         cross join period\n" +
		"         inner join recipe r on co.recipe_id = r.id\n" +
		"         inner join consumers c on c.id = r.consumer_id\n" +
		"         inner join drug d on r.drug_id = d.id\n" +
		"where (drug_id in (select * from need_drug_ids) or d.type_id in (select * from need_types))\n" +
		"  and period.start <= co.order_date\n" +
		"  and co.order_date <= period.\"end\";",
	)

	fmt.Println(b.String())

	rv := make(map[daos.ID]daos.Consumer)

	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)

	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetCriticalDrugs(ctx context.Context) (map[daos.ID]daos.Drug, error) {

	query := "select d.id, d.name_, d.type_id, d.critical_count, d.price\n" +
		"from drug d\n" +
		"         left join storage s on d.id = s.comp_id\n" +
		"where coalesce(s.count_, 0) <= d.critical_count;"

	rv := make(map[daos.ID]daos.Drug)
	err := daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		query,
		&rv,
	)

	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetDrugsWithMinStorageCount(ctx context.Context, needTypes []int) (map[daos.ID]daos.Drug, error) {

	var b strings.Builder

	b.WriteString("with ")

	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n" +
		"     drug_count as (select id, coalesce(count_, 0) count\n" +
		"from drug\n" +
		"    left join storage s\n" +
		"on drug.id = s.comp_id\n" +
		"where drug.type_id in (select * from need_types))\n" +
		"select d.id, d.name_, d.type_id, d.critical_count, d.price\n" +
		"from drug_count\n" +
		"         inner join drug d on drug_count.id = d.id\n" +
		"where count = (select min(count) from drug_count);")

	rv := make(map[daos.ID]daos.Drug)
	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)

	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetInProgressOrders(ctx context.Context) (map[daos.ID]daos.ConsumerOrder, error) {

	query := "select co.id, co.recipe_id, co.status, co.complete_time, co.order_date\n" +
		"from consumer_order co\n" +
		"         inner join recipe r on recipe_id = r.id\n" +
		"where co.status = 'IN-PROGRESS';"

	rv := make(map[daos.ID]daos.ConsumerOrder)
	err := daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		query,
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetComponentsForInProgressOrders(ctx context.Context) (map[daos.ID]daos.Drug, error) {

	query := "with drugs_comp as (select td.drug_id, tc.comp_id, tc.count_\n" +
		"                    from technology_components tc\n" +
		"                             inner join technology_drug td on tc.tech_id = td.tech_id)\n" +
		"select distinct d.id, d.name_, d.type_id, d.critical_count, d.price\n" +
		"from consumer_order co\n" +
		"         inner join recipe r on r.id = co.recipe_id\n" +
		"         inner join drugs_comp dc on dc.drug_id = r.drug_id\n" +
		"         inner join drug d on dc.drug_id = d.id\n" +
		"where co.status = 'IN-PROGRESS';"

	rv := make(map[daos.ID]daos.Drug)
	err := daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		query,
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (m *Model) GetTechsForDrugs(ctx context.Context, needDrugs []daos.ID, needTypes []int) (map[daos.ID]daos.TechnologyBook, error) {
	var b strings.Builder
	b.WriteString("with ")
	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n")
	err = formNeedDrugsSubQuery(&b, needDrugs)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n" +
		"drugs_in_progress as (select distinct r.drug_id\n" +
		"                           from consumer_order co\n" +
		"                                    inner join recipe r on r.id = co.recipe_id\n" +
		"                           where co.status = 'IN-PROGRESS')\n" +
		"select tb.id, tb.description\n" +
		"from technology_drug td\n" +
		"         inner join technology_book tb on td.tech_id = tb.id\n" +
		"where td.drug_id in (select * from drugs_in_progress)\n" +
		"   or td.drug_id in (select * from need_types)\n" +
		"   or td.drug_id in (select * from need_drug_ids);")

	rv := make(map[daos.ID]daos.TechnologyBook)
	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

type ConsumersWithCount struct {
	Name      string    `sql:"name_"`
	Phone     string    `sql:"phone"`
	Address   string    `sql:"addres"`
	BirthDate time.Time `sql:"birth_date"`
	Count     int
}

func (m *Model) GetMostOftenConsumers(ctx context.Context, needDrugs []daos.ID, needTypes []int) (map[daos.ID]ConsumersWithCount, error) {
	var b strings.Builder
	b.WriteString("with ")
	err := m.formNeedTypesSubQuery(&b, needTypes)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n")
	err = formNeedDrugsSubQuery(&b, needDrugs)
	if err != nil {
		return nil, err
	}
	b.WriteString(",\n" +
		"     help_table_1 as (select r.consumer_id, count(*) count\n" +
		"from consumer_order co\n" +
		"    inner join recipe r\n" +
		"on co.recipe_id = r.id\n" +
		"    inner join drug d on r.drug_id = d.id\n" +
		"where d.id in (select * from need_drug_ids)\n" +
		"   or d.type_id in (select * from need_types)\n" +
		"group by r.consumer_id),\n" +
		"    consumer_count as (\n" +
		"select id, coalesce (ht.count, 0) count\n" +
		"from consumers\n" +
		"    left join help_table_1 ht\n" +
		"on id = ht.consumer_id)\n" +
		"select c.id, c.name_, c.phone, c.addres, c.birth_date, cc.count\n" +
		"from consumer_count cc\n" +
		"         inner join consumers c on c.id = cc.id\n" +
		"where count = (select max(count) from consumer_count);",
	)
	//fmt.Println(b.String())

	rv := make(map[daos.ID]ConsumersWithCount)
	err = daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		b.String(),
		&rv,
	)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

type Component struct {
	Name  string
	Count int
	Price int
}
type DrugDetail struct {
	daos.Drug
	TechDescription string
	Components      map[daos.ID]Component
}

func (m *Model) GetDrugDetailInfo(ctx context.Context, needDrug daos.ID) (DrugDetail, error) {
	//fmt.Println(1)
	query1 := "select d.name_, d.type_id, d.critical_count, d.price, coalesce(tb.description, 'nil') \n" +
		"from drug d\n" +
		"         left join technology_drug td on d.id = td.drug_id\n" +
		"         left join technology_book tb on td.tech_id = tb.id\n" +
		"where d.id in (select  $1::int);"

	var dd DrugDetail

	err := m.db.QueryRowContext(ctx, query1, needDrug).Scan(&dd.Name, &dd.TypeId, &dd.CriticalCount, &dd.Price, &dd.TechDescription)
	if err != nil {
		return DrugDetail{}, err
	}

	query2 := "with need_drug_ids as (select $1::int)\n" +
		"select d.id, d.name_, count_, count_ * price price\n" +
		"from technology_drug td\n" +
		"         inner join technology_components tc on td.tech_id = tc.tech_id\n" +
		"         inner join drug d on tc.comp_id = d.id\n" +
		"where td.drug_id in (select * from need_drug_ids);"

	rows, err := m.db.QueryContext(
		ctx,
		query2,
		needDrug,
	)
	if err != nil {
		return DrugDetail{}, err
	}
	defer rows.Close()

	dd.Components = make(map[daos.ID]Component)
	for rows.Next() {
		var c Component
		var id daos.ID
		err = rows.Scan(&id, &c.Name, &c.Count, &c.Price)
		if err != nil {
			return DrugDetail{}, err
		}

		dd.Components[id] = c
	}

	if err := rows.Err(); err != nil {
		return DrugDetail{}, err
	}

	return dd, nil
}

func (m *Model) ShowStorage(ctx context.Context) (map[daos.ID]DrugsWithCount, error) {

	rv := make(map[daos.ID]DrugsWithCount)

	query := "SELECT d.id, d.name_, d.type_id, d.critical_count, d.price, s.count_ from storage s inner join drug d on d.id = s.comp_id"

	err := daos.ExecQueryAndReadAll(
		ctx,
		m.db,
		query,
		&rv,
	)

	if err != nil {
		return nil, err
	}

	return rv, nil
}
