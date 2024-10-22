package utils

import (
	"Cloud/models"
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
	"time"
)

func SendExcelWithLogs(logs []models.RequestLog, w http.ResponseWriter) error {
	// Создайте новый файл Excel
	f := excelize.NewFile()

	// Установите заголовки
	headers := []string{"Method", "Endpoint", "UserID", "IP", "UserAgent", "Time", "StatusCode", "Duration"}
	for i, header := range headers {
		if err := f.SetCellValue("Sheet1", string('A'+i)+"1", header); err != nil {
			return fmt.Errorf("ошибка при установке заголовка: %v", err)
		}
	}

	// Заполняем данными
	for i, log := range logs {
		row := i + 2 // Начинаем с 2-й строки, так как первая строка — заголовки
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), log.Method); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), log.Endpoint); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), log.UserID); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), log.IP); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), log.UserAgent); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), log.Time.Format(time.RFC3339)); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), log.StatusCode); err != nil {
			return err
		}
		if err := f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), log.Duration.Seconds()); err != nil {
			return err
		}
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=request_logs.xlsx")

	// Записываем файл в ответ
	if err := f.Write(w); err != nil {
		return err
	}

	return nil
}
