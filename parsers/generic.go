package parsers

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/Popov-Dmitriy-Ivanovich/genmilk_backend/models"
)

type DefaultParser interface {
}

type ErrorNonPointerType struct{}

func (enpt ErrorNonPointerType) Error() string {
	return "ParseFromRecord got non pointer type"
}

var timeFormats = []string{
	time.DateOnly,
	"02.01.2006",
	"02.01.06",
	"02/01/2006",
	"02/01/06",
	"02-01-2006",
	"02-01-06",
}

func ParseTime(timeStr string) (time.Time, error) {
	for _, format := range timeFormats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("дата не соответсвует ни одному из доступных форматов")
}

func ParseFromRecord(dest any, record []string, headerIndexes map[string]int) error {

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Pointer {
		return ErrorNonPointerType{}
	}

	value := reflect.ValueOf(dest).Elem()

	for i := 0; i <= value.NumField(); i++ {

		fieldColumnTag := value.Type().Field(i).Tag.Get("csv_column")
		fieldColumnTypeTag := value.Type().Field(i).Tag.Get("csv_type")
		fieldValue := value.Field(i)

		switch fieldColumnTypeTag {
		case "*float64":
			strValue := record[headerIndexes[fieldColumnTag]]
			if strValue == "" {
				continue
			}
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			floatValue, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return err
			}
			fieldValue.Elem().SetFloat(floatValue)

		case "*int":
			strValue := record[headerIndexes[fieldColumnTag]]
			if strValue == "" {
				continue
			}
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			intValue, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.Elem().SetInt(intValue)

		case "*uint":
			strValue := record[headerIndexes[fieldColumnTag]]
			if strValue == "" {
				continue
			}
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			uintValue, err := strconv.ParseUint(strValue, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.Elem().SetUint(uintValue)

		case "DateOnly":
			dateAsTime, err := ParseTime(record[headerIndexes[fieldColumnTag]])
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(models.DateOnly{Time: dateAsTime}))

		case "bool":
			fieldValue.SetBool(record[headerIndexes[fieldColumnTag]] == "1")

		case "int":
			strValue := record[headerIndexes[fieldColumnTag]]
			intValue, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(intValue)

		case "uint":
			strValue := record[headerIndexes[fieldColumnTag]]
			uintValue, err := strconv.ParseUint(strValue, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.Elem().SetUint(uintValue)

		}
	}
	return nil
}
