package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"sync"
)

// Reader is an interface for reading structured CSV files.
type Reader[T any] interface {
	Read() (T, error)
}

// NewReader creates a new CSV reader.
func NewReader[T any](r io.Reader) Reader[T] {
	return &csvReader[T]{
		r: csv.NewReader(r),
	}
}

type csvReader[T any] struct {
	r                  *csv.Reader
	columnsOnce        sync.Once
	columnIndexToField map[int]*structTag
}

func (r *csvReader[T]) parseHeaders() error {
	// Read the headers row
	columns, err := r.r.Read()
	if err != nil {
		if err == io.EOF {
			return err
		}
		return fmt.Errorf("reading header row: %w", err)
	}

	// Map column names to indices
	columnToIndex := make(map[string]int)
	for i, column := range columns {
		columnToIndex[column] = i
	}

	// Create a mapping of columns to struct fields
	var row T
	r.columnIndexToField = make(map[int]*structTag)
	typ := reflect.TypeOf(row)
	for i := 0; i < typ.NumField(); i++ {
		// Get the struct tag "csv" on the field
		csvTag := parseStructTag(typ.Field(i))
		if csvTag == nil {
			continue
		}

		// Get the index of the column
		columnIndex, ok := columnToIndex[csvTag.ColumnName]
		if !ok {
			continue
		}

		// Map the column index to the struct field name
		r.columnIndexToField[columnIndex] = csvTag
	}
	return nil
}

func (r *csvReader[T]) parseHeadersOnce() error {
	var err error
	r.columnsOnce.Do(func() {
		err = r.parseHeaders()
	})
	return err
}

func (r *csvReader[T]) Read() (T, error) {
	var row T

	// Parse the headers if we haven't already
	if err := r.parseHeadersOnce(); err != nil {
		return row, err
	}

	// Read a row from the CSV
	columns, err := r.r.Read()
	if err != nil {
		return row, err
	}

	// Map the columns to the struct fields
	rowReflectVal := reflect.Indirect(reflect.ValueOf(&row))
	for i, valueStr := range columns {
		// Get the struct field that corresponds to the column
		csvTag, ok := r.columnIndexToField[i]
		if !ok {
			continue
		}

		// Set the value of the struct field
		field := rowReflectVal.FieldByName(csvTag.FieldName)
		switch field.Kind() {
		case reflect.String:
			field.SetString(valueStr)
		}
	}
	return row, nil
}

func ReadAll[T any](r Reader[T]) ([]T, error) {
	var rows []T
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return rows, nil
			}
			return rows, err
		}
		rows = append(rows, row)
	}
}
