package pharmacy

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"pharmacy/interanl/pharmacy/daos"
)

type Model struct {
	db        *sql.DB
	idToUses  map[int]daos.DrugUses
	idToTypes map[int]daos.DrugType
	usesToId  map[daos.DrugUses]int
	typesToId map[daos.DrugType]int
	typeToUse map[int]map[int]struct{}
}

//go:embed create_tables.sql
var createTablesQueres string

func (m *Model) initializeUses(ctx context.Context, tx *sql.Tx) error {
	tableName := daos.MapEnumToTableName[daos.DrugUsesE]
	rows, err := tx.QueryContext(ctx,
		fmt.Sprintf(
			"SELECT id, use_text FROM %s",
			tableName,
		),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id daos.ID
		var use daos.DrugUses
		if err := rows.Scan(&id, &use.UseText); err != nil {
			return err
		}

		m.idToUses[id] = use
		m.usesToId[use] = id
	}

	if err := rows.Err(); err != nil {
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (use_text) VALUES ($1) RETURNING id", tableName)

	for _, use := range standardUses {
		_, ok := m.usesToId[use]

		if !ok {
			var id daos.ID
			err := tx.QueryRowContext(ctx, insertQuery, use.UseText).Scan(&id)
			if err != nil {
				return err
			}

			m.usesToId[use] = id
			m.idToUses[id] = use
		}
	}

	return nil
}

func (m *Model) initializeTypes(ctx context.Context, tx *sql.Tx) error {
	tableName := daos.MapEnumToTableName[daos.DrugTypeE]
	rows, err := tx.QueryContext(ctx,
		fmt.Sprintf(
			"SELECT id, name_, produced FROM %s",
			tableName,
		),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id daos.ID
		var typ daos.DrugType
		if err := rows.Scan(&id, &typ.Name, &typ.Produced); err != nil {
			return err
		}

		m.idToTypes[id] = typ
		m.typesToId[typ] = id
	}

	if err := rows.Err(); err != nil {
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (name_, produced) VALUES ($1, $2) RETURNING id", tableName)

	for _, typ := range standardTypes {
		_, ok := m.typesToId[typ]

		if !ok {
			var id daos.ID
			err := tx.QueryRowContext(ctx, insertQuery, typ.Name, typ.Produced).Scan(&id)
			if err != nil {
				return err
			}

			m.typesToId[typ] = id
			m.idToTypes[id] = typ
		}
	}

	return nil

}

func (m *Model) initializeUsesForTypes(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx,
		"SELECT type_id, use_id FROM type_uses",
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var typeId daos.ID
		var useId daos.ID
		if err := rows.Scan(&typeId, &useId); err != nil {
			return err
		}

		set, ok := m.typeToUse[typeId]
		if !ok {
			set = make(map[int]struct{})
			m.typeToUse[typeId] = set
		}

		set[useId] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	insertQuery := "INSERT INTO type_uses (type_id, use_id) VALUES ($1, $2)"

	for typeI, uses := range standardUsesForTypes {
		tIdSql := m.typesToId[standardTypes[typeI]]
		usesHas, ok := m.typeToUse[tIdSql]

		if ok {
			for i := 0; i < len(uses); i++ {
				uIdSql := m.usesToId[standardUses[uses[i]]]
				_, okT := usesHas[uIdSql]
				if !okT {
					_, err := tx.ExecContext(ctx, insertQuery, tIdSql, uIdSql)
					if err != nil {
						return err
					}
					m.typeToUse[tIdSql][uIdSql] = struct{}{}
				}

			}
		} else {
			m.typeToUse[tIdSql] = make(map[int]struct{})
			for i := 0; i < len(uses); i++ {
				uIdSql := m.usesToId[standardUses[uses[i]]]
				_, err := tx.ExecContext(ctx, insertQuery, tIdSql, uIdSql)
				if err != nil {
					return err
				}

				m.typeToUse[tIdSql][uIdSql] = struct{}{}
			}

		}
	}

	return nil

}

func (m *Model) initializeUsesTypes(ctx context.Context, tx *sql.Tx) error {
	err := m.initializeUses(ctx, tx)
	if err != nil {
		return err
	}

	err = m.initializeTypes(ctx, tx)
	if err != nil {
		return err
	}

	return m.initializeUsesForTypes(ctx, tx)

}

func NewModel(ctx context.Context, db *sql.DB) (*Model, error) {

	_, err := db.ExecContext(ctx, createTablesQueres)
	if err != nil {
		return nil, err
	}

	m := &Model{db: db}
	m.typeToUse = make(map[int]map[int]struct{})
	m.idToTypes = make(map[int]daos.DrugType)
	m.typesToId = make(map[daos.DrugType]int)
	m.idToUses = make(map[int]daos.DrugUses)
	m.usesToId = make(map[daos.DrugUses]int)

	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	err = m.initializeUsesTypes(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Model) Close() {
	err := m.db.Close()
	if err != nil {
		panic(err)
	}
}
