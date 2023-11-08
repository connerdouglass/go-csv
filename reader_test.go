package csv_test

import (
	"embed"
	"io"
	"testing"

	"github.com/connerdouglass/go-csv"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/*
	testdata embed.FS
)

func TestReadCsv(t *testing.T) {
	t.Run("read basic csv file", func(t *testing.T) {
		type Row struct {
			FirstName string `csv:"First Name"`
			LastName  string `csv:"Last Name"`
			Email     string `csv:"Email"`
			City      string `csv:"City"`
			State     string `csv:"State"`
		}
		file, _ := testdata.Open("testdata/strings.csv")
		defer file.Close()
		reader := csv.NewReader[Row](file)

		// Expected results
		expected := []Row{
			{"John", "Doe", "johndoe@email.com", "Los Angeles", "CA"},
		}

		for i, expectedRow := range expected {
			// Read the rows
			row, err := reader.Read()
			require.NoError(t, err, "error reading row %d", i)
			require.Equal(t, expectedRow, row, "row %d does not match", i)
		}

		_, err := reader.Read()
		require.Equal(t, io.EOF, err, "expected EOF error")
	})
}
