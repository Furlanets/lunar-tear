package masterdata

import (
	"fmt"

	"lunar-tear/server/internal/utils"
)

type ConsumableItemRow struct {
	ConsumableItemId int32 `json:"ConsumableItemId"`
	SellPrice        int32 `json:"SellPrice"`
}

type ConsumableItemCatalog struct {
	All map[int32]ConsumableItemRow
}

func LoadConsumableItemCatalog() (*ConsumableItemCatalog, error) {
	rows, err := utils.ReadJSON[ConsumableItemRow]("EntityMConsumableItemTable.json")
	if err != nil {
		return nil, fmt.Errorf("load consumable item table: %w", err)
	}

	catalog := &ConsumableItemCatalog{
		All: make(map[int32]ConsumableItemRow, len(rows)),
	}
	for _, row := range rows {
		catalog.All[row.ConsumableItemId] = row
	}
	return catalog, nil
}
