package pharmacy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pharmacy/interanl/pharmacy/daos"
	"slices"
)

func (m *Model) GetDaoByType(daoType int) (*daos.GeneralDao, error) {
	tableName, ok := daos.MapEnumToTableName[daoType]
	if !ok {
		return nil, errors.New("dao of this type not supported")
	}

	return daos.NewGeneralDao(tableName, m.db), nil
}

func diffCountInTableWithTransaction(ctx context.Context, tx *sql.Tx, count int, tablename, colName, condition string, args ...any) error {

	var hasCount int
	query1 := fmt.Sprintf("SELECT %s FROM %s WHERE %s FOR UPDATE", colName, tablename, condition)
	err := tx.QueryRowContext(
		ctx,
		query1,
		args...,
	).Scan(&hasCount)
	if err != nil {
		return err
	}
	if hasCount+count < 0 {
		return errors.New(colName + " in " + tablename + " not enough")
	}

	hasCount += count

	query2 := fmt.Sprintf("UPDATE %s SET %s = $%d WHERE %s", tablename, colName, len(args)+1, condition)

	myArgs := slices.Clone(args)
	myArgs = append(myArgs, hasCount)

	//fmt.Println(query2)
	_, err = tx.ExecContext(
		ctx,
		query2,
		myArgs...,
	)
	return err
}

func (m *Model) TakeOutOfStorage(ctx context.Context, drugId daos.ID, count int) error {
	if count <= 0 {
		return fmt.Errorf("invalid count: %d", count)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = diffCountInTableWithTransaction(ctx, tx, -count, "storage", "count_", "comp_id = $1", drugId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO inventasrization(drug_id, date_, count_) VALUES ($1::int, current_timestamp, $2::int)",
		drugId,
		count,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Model) AddDrugsToStorage(ctx context.Context, drugId daos.ID, count int) error {
	if count <= 0 {
		return fmt.Errorf("invalid count: %d", count)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = diffCountInTableWithTransaction(ctx, tx, count, "storage", "count_", "comp_id = $1", drugId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Model) AddDrugIdForTech(ctx context.Context, techId daos.ID, drugId daos.ID) error {

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ra, err := tx.ExecContext(
		ctx,
		"UPDATE technology_drug SET tech_id = $1 WHERE technology_drug.drug_id = $2",
		techId,
		drugId,
	)
	if err != nil {
		return err
	}

	if raff, err := ra.RowsAffected(); err != nil {
		return err
	} else if raff == 0 {

		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO technology_drug(tech_id, drug_id) VALUES ($1, $2)",
			techId,
			drugId,
		)
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}

func (m *Model) DiffCompIdForTech(ctx context.Context, techId daos.ID, compId daos.ID, count int) error {

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = diffCountInTableWithTransaction(ctx, tx, count, "technology_components", "count_", "comp_id = $1 and tech_id = $2", compId, techId)
	if errors.Is(err, sql.ErrNoRows) {
		if count < 0 {
			return fmt.Errorf("invalid count: %d", count)
		}

		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO technology_components(comp_id, tech_id, count_) Values ($1, $2, $3)",
			compId,
			techId,
			count,
		)
		return tx.Commit()
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}
