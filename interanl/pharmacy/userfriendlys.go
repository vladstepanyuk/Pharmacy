package pharmacy

import (
	"fmt"
	"pharmacy/interanl/pharmacy/daos"
)

var standardUses = []daos.DrugUses{
	{"внутреннее"},
	{"наружное"},
	{"смешать"},
}

var standardTypes = []daos.DrugType{
	{"таблетки", false},
	{"мазь", false},
	{"настойка", false},
	{"микстура", true},
	{"мазь", true},
	{"раствор", true},
	{"настойка", true},
	{"порошок", true},
}

var nameToStdType = map[string][]int{
	"таблетки": {0},
	"мазь":     {1, 4},
	"настойка": {2, 6},
	"микстура": {3},
	"раствор":  {5},
	"порошок":  {7},
}

var prodToStdType = [][]int{
	0: {0, 1, 2},
	1: {3, 4, 5, 6, 7},
}

// микстуры, порошки -> внутреннее
// растворы, настойки -> наружное, внутреннее, смешать
// таблетки -> внутреннее
// мази -> наружное

var standardUsesForTypes = map[int][]int{
	1: {1},
	4: {1},
	0: {0},
	5: {0, 1, 2},
	2: {0, 1, 2},
	6: {0, 1, 2},
	3: {0},
	7: {0},
}

func StdTypesNames() []string {
	return []string{
		"таблетки",
		"мазь",
		"настойка",
		"микстура",
		"раствор",
		"порошок",
	}
}

func GetTypeByNameAndProduced(m *Model, name string, produced bool) (daos.ID, error) {
	types1, ok := nameToStdType[name]
	if !ok {
		return 0, fmt.Errorf("no type with name %s", name)
	}

	i := 0
	if produced {
		i = 1
	}

	types2 := prodToStdType[i]

	for _, t1 := range types1 {
		for _, t2 := range types2 {
			if t1 == t2 {
				return m.typesToId[standardTypes[t1]], nil
			}
		}
	}

	return 0, fmt.Errorf("no type with name %s and produced %v", name, produced)
}

func GetTypesIDByName(_ *Model, name string) ([]int, error) {
	types1, ok := nameToStdType[name]
	if !ok {
		return nil, fmt.Errorf("no type with name %s", name)
	}

	return types1, nil
}

func GetTypeNameByID(m *Model, id daos.ID) string {
	typ, ok := m.idToTypes[id]
	if !ok {
		panic(fmt.Errorf("no type with id %d", id))
	}

	return typ.Name
}

func GetUseId(m *Model, name string) (daos.ID, error) {
	for _, us := range standardUses {
		if us.UseText == name {
			return m.usesToId[us], nil
		}
	}

	return 0, fmt.Errorf("no use with name %s", name)
}

func GetStatusString(i int) string {
	switch i {
	case daos.SReady:
		return "READY"
	case daos.SWaitComps:
		return "WAITING-COMPS"
	case daos.SWait:
		return "WAITING"
	case daos.SInProgress:
		return "IN-PROGRESS"
	default:
		panic("aflwasef")
	}
}
