package parsers

import (
	"strconv"

	"github.com/Popov-Dmitriy-Ivanovich/genmilk_backend/models"
)

type Lactation struct {
	CowSelecs        uint            `csv_column:"CowSelecs" csv_type:"uint"`
	Number           uint            `csv_column:"Number" csv_type:"uint"`
	InsemenationNum  int             `csv_column:"InsemenationNum" csv_type:"int"`
	InsemenationDate models.DateOnly `csv_column:"InsemenationDate" csv_type:"DateOnly"`
	CalvingCount     int             `csv_column:"CalvingCount" csv_type:"int"`
	CalvingDate      models.DateOnly `csv_column:"CalvingDate" csv_type:"DateOnly"`
	Abort            bool            `csv_column:"Abort" csv_type:"bool"`
	MilkAll          *float64        `csv_column:"MilkAll" csv_type:"*float64"`
	FatAll           *float64        `csv_column:"FatAll" csv_type:"*float64"`
	ProteinAll       *float64        `csv_column:"ProteinAll" csv_type:"*float64"`
	Milk305          *float64        `csv_column:"Milk305" csv_type:"*float64"`
	Fat305           *float64        `csv_column:"Fat305" csv_type:"*float64"`
	Protein305       *float64        `csv_column:"Protein305" csv_type:"*float64"`
	Days             *int            `csv_column:"Days" csv_type:"*int"`
	ServicePeriod    *uint           `csv_column:"ServicePeriod" csv_type:"*uint"`
}

type NotFoundCowError struct {
	CowSelecs uint
}

func (nfce NotFoundCowError) Error() string {
	return "Не удалось найти корову с селексом " + strconv.FormatUint(uint64(nfce.CowSelecs), 10)
}

func (l Lactation) ToDbModel() (any, error) {
	lac := models.Lactation{
		Number:           l.Number,
		InsemenationNum:  l.InsemenationNum,
		InsemenationDate: l.InsemenationDate,
		CalvingCount:     l.CalvingCount,
		Abort:            l.Abort,
		MilkAll:          l.MilkAll,
		FatAll:           l.FatAll,
		ProteinAll:       l.ProteinAll,
		Milk305:          l.Milk305,
		Fat305:           l.Fat305,
		Protein305:       l.Protein305,
		Days:             l.Days,
		ServicePeriod:    l.ServicePeriod,
	}
	db := models.GetDb()
	cowId := uint(0)
	if err := db.Model(&models.Cow{}).Where(map[string]any{"selecs_number": l.CowSelecs}).Limit(1).Pluck("id", cowId).Error; err != nil || cowId == 0 {
		return nil, NotFoundCowError{CowSelecs: l.CowSelecs}
	}
	lac.CowId = cowId
	return lac, nil
}

func (l Lactation) ParseFromRecord(record []string, headerIndexes map[string]int) Lactation {
	return Lactation{}
}
