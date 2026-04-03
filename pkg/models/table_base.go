package models

import (
	"fmt"
	"time"
	
	"github.com/lanwenhong/planet_8583/pkg/utils"
)

type TableSuffix struct {
	Year  int
	Month int
}

func ParseTableSuffix(s string) *TableSuffix {
	utils.MustTrue(
		len(s) >= 6, fmt.Errorf("invalid table suffix format: %s", s),
	)
	year := utils.SafeToInt(s[:4])
	month := utils.SafeToInt(s[4:6])
	utils.MustTrue(
		year > 0 && month > 0 && month <= 12, fmt.Errorf("invalid date: %s", s[:6]),
	)
	return &TableSuffix{Year: year, Month: month}
}

func (ts *TableSuffix) GetTableName(baseTable string) string {
	return fmt.Sprintf("%s_%04d%02d", baseTable, ts.Year, ts.Month)
}

func (ts *TableSuffix) GetNextMonth() *TableSuffix {
	if ts.Month == 12 {
		return &TableSuffix{Year: ts.Year + 1, Month: 1}
	}
	return &TableSuffix{Year: ts.Year, Month: ts.Month + 1}
}

func (ts *TableSuffix) IsBeforeNow() bool {
	now := time.Now()
	return ts.Year < now.Year() || (ts.Year == now.Year() && ts.Month < int(now.Month()))
}
