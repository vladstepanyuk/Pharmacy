package daos

import "time"

type ID = int

const (
	DrugE = iota
	DrugTypeE
	DrugUsesE
	DoctorE
	ConsumerE
	TechnologyBookE
	RecipeE
	ConsumerOrderE
)

var MapEnumToTableName = map[int]string{
	DrugE:           "drug",
	DrugTypeE:       "drug_type",
	DrugUsesE:       "drug_uses",
	DoctorE:         "doctor",
	ConsumerE:       "consumers",
	TechnologyBookE: "technology_book",
	RecipeE:         "recipe",
	ConsumerOrderE:  "consumer_order",
}

const (
	SReady = iota
	SWaitComps
	SWait
	SInProgress
)

type DrugType struct {
	Name     string `sql:"name_"`
	Produced bool   `sql:"produced"`
}

type Drug struct {
	Name          string `sql:"name_"`
	TypeId        ID     `sql:"type_id"`
	CriticalCount int    `sql:"critical_count"`
	Price         int    `sql:"price"`
}

type DrugUses struct {
	UseText string `sql:"use_text"`
}

type Doctor struct {
	Name string `sql:"name_"`
}

type Consumer struct {
	Name      string    `sql:"name_"`
	Phone     string    `sql:"phone"`
	Address   string    `sql:"addres"`
	BirthDate time.Time `sql:"birth_date"`
}

type TechnologyBook struct {
	Description string `sql:"description"`
}

type Recipe struct {
	ConsumerId ID     `sql:"consumer_id"`
	DoctorId   ID     `sql:"doctor_id"`
	DrugId     ID     `sql:"drug_id"`
	UseId      ID     `sql:"use_id"`
	Disease    string `sql:"disease"`
	Count      int    `sql:"drug_count"`
}

type ConsumerOrder struct {
	ResipeId     ID     `sql:"recipe_id"`
	Status       string `sql:"status"`
	CompleteTime string `sql:"complete_time"`
	OrderTime    string `sql:"order_date"`
}
