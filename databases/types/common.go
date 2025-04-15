package types

import "fmt"

type Data map[string]interface{}

func ConvertMapToTypeArray(columnTypeMap map[string]string) []string {
	types := []string{}
	for _, v := range columnTypeMap {
		types = append(types, v)
	}
	return types
}

func ExtractHeaders(columnTypeMap map[string]string) []string {
	headers := []string{}
	for k := range columnTypeMap {
		headers = append(headers, k)
	}
	return headers
}

func ConvertDataToRows(data []Data) [][]string {
	rows := [][]string{}
	for _, item := range data {
		row := []string{}
		for _, v := range item {
			row = append(row, fmt.Sprintf("%v", v)) // Converte valores para strings
		}
		rows = append(rows, row)
	}
	return rows
}
