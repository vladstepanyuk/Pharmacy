package daos

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type GeneralDao struct {
	tableName string
	db        *sql.DB
}

func NewGeneralDao(tableName string, db *sql.DB) *GeneralDao {
	return &GeneralDao{tableName: tableName, db: db}
}

func getNamesFields(typ reflect.Type) []string {
	names := make([]string, typ.NumField())
	for i := range names {
		field := typ.Field(i)

		var nameField string
		tag := field.Tag.Get("sql")
		if tag == "" {
			nameField = strings.ToLower(field.Name)
		}

		nameField = tag
		names[i] = nameField
	}

	return names
}

func joinStringsToBuilder(b *strings.Builder, strs []string, sep string) {
	for i, str := range strs {
		b.WriteString(str)
		if i != len(strs)-1 {
			b.WriteString(sep)
		}
	}
}

func joinIndexesToBuilder(b *strings.Builder, from, to int, sep string) {

	for i := from; i < to; i++ {
		b.WriteString("$")
		b.WriteString(strconv.Itoa(i))

		if i != to-1 {
			b.WriteString(sep)
		}
	}
}

func prepareInsertQueryString(tableName string, typ reflect.Type) string {

	names := getNamesFields(typ)

	var b strings.Builder
	b.WriteString("INSERT INTO ")
	b.WriteString(tableName)
	b.WriteString("(")

	joinStringsToBuilder(&b, names, ", ")

	b.WriteString(") VALUES (")

	joinIndexesToBuilder(&b, 1, len(names)+1, ", ")

	b.WriteString(") RETURNING id")

	return b.String()
}

func prepareUpdateQueryString(tableName string, typ reflect.Type) string {

	names := getNamesFields(typ)

	var b strings.Builder
	b.WriteString("UPDATE ")
	b.WriteString(tableName)
	b.WriteString(" SET (")

	joinStringsToBuilder(&b, names, ", ")

	b.WriteString(") = ROW(")

	joinIndexesToBuilder(&b, 2, len(names)+2, ", ")

	b.WriteString(") WHERE id = $1")

	return b.String()
}

func prepareLookupQueryString(tableName string, typ reflect.Type) string {

	names := getNamesFields(typ)

	var b strings.Builder
	b.WriteString("SELECT ")

	joinStringsToBuilder(&b, names, ", ")
	b.WriteString(" FROM ")
	b.WriteString(tableName)
	b.WriteString(" WHERE id = $1")

	return b.String()
}

func prepareListQueryString(tableName string, typ reflect.Type) string {

	names := getNamesFields(typ)

	var b strings.Builder
	b.WriteString("SELECT id, ")

	joinStringsToBuilder(&b, names, ", ")
	b.WriteString(" FROM ")
	b.WriteString(tableName)

	return b.String()
}

func (g *GeneralDao) Insert(ctx context.Context, obj any) (ID, error) {

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Struct {
		panic("obj must be a struct")
	}

	args := make([]any, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		args = append(args, field.Interface())
	}

	typ := val.Type()

	query := prepareInsertQueryString(g.tableName, typ)

	fmt.Println(query)

	var id ID
	err := g.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (g *GeneralDao) Update(ctx context.Context, id ID, obj any) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Struct {
		panic("obj must be a struct")
	}

	args := make([]any, 0, val.NumField()+1)
	args = append(args, id)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		args = append(args, field.Interface())
	}

	typ := val.Type()

	query := prepareUpdateQueryString(g.tableName, typ)

	res, err := g.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if ra, err := res.RowsAffected(); err != nil || ra != 1 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (g *GeneralDao) Delete(ctx context.Context, id ID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", g.tableName)
	_, err := g.db.ExecContext(ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (g *GeneralDao) LookUp(ctx context.Context, id ID, obj any) error {

	pval := reflect.ValueOf(obj)
	if pval.Kind() != reflect.Pointer {
		panic("obj must be a pointer to a struct")
	}

	val := pval.Elem()
	typ := val.Type()

	query := prepareLookupQueryString(g.tableName, typ)

	pFields := make([]any, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		pfield := field.Addr()
		pFields = append(pFields, pfield.Interface())
	}

	err := g.db.QueryRowContext(ctx,
		query, id).Scan(pFields...)

	return err
}

func (g *GeneralDao) List(ctx context.Context, out any) error {
	pval := reflect.ValueOf(out)

	if pval.Kind() != reflect.Ptr || pval.Elem().Kind() != reflect.Map {
		panic("out must be a pointer to a map")
	}

	outVal := pval.Elem()

	elemType := outVal.Type().Elem()

	if elemType.Kind() != reflect.Struct || outVal.Type().Key().Kind() != reflect.Int {
		panic("out must be a map from int to struct")
	}

	query := prepareListQueryString(g.tableName, elemType)

	rows, err := g.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var id int
	tmpElem := reflect.New(elemType)

	args := make([]any, 0, elemType.NumField()+1)
	args = append(args, &id)
	for i := 0; i < tmpElem.Elem().NumField(); i++ {
		tmpField := tmpElem.Elem().Field(i)
		args = append(args, tmpField.Addr().Interface())
	}

	for rows.Next() {
		err := rows.Scan(args...)

		if err != nil {
			return err
		}

		outVal.SetMapIndex(reflect.ValueOf(id), tmpElem.Elem())
	}

	err = rows.Err()

	return err
}

func ExecQueryAndReadAll(ctx context.Context, db *sql.DB, query string, out any, args ...any) error {
	pval := reflect.ValueOf(out)

	if pval.Kind() != reflect.Ptr || pval.Elem().Kind() != reflect.Map {
		panic("out must be a pointer to a map")
	}

	outVal := pval.Elem()

	elemType := outVal.Type().Elem()

	if elemType.Kind() != reflect.Struct || outVal.Type().Key().Kind() != reflect.Int {
		panic("out must be a map from int to struct")
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var id ID
	tmpElem := reflect.New(elemType)

	fields := make([]any, 0, elemType.NumField()+1)
	fields = append(fields, &id)
	for i := 0; i < tmpElem.Elem().NumField(); i++ {
		tmpField := tmpElem.Elem().Field(i)
		fields = append(fields, tmpField.Addr().Interface())
	}

	for rows.Next() {
		err := rows.Scan(fields...)

		if err != nil {
			return err
		}

		outVal.SetMapIndex(reflect.ValueOf(id), tmpElem.Elem())
	}

	err = rows.Err()

	return err
}
