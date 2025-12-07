package excel

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// XLSXReader XLSX 格式读取器 (Excel 2007+)
type XLSXReader struct {
	file *excelize.File
}

// NewXLSXReader 创建 XLSX 读取器
func NewXLSXReader(reader io.Reader) (*XLSXReader, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}

	return &XLSXReader{file: file}, nil
}

// Close 关闭 Excel 文件
func (r *XLSXReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// GetSheetNames 获取所有工作表名称
func (r *XLSXReader) GetSheetNames() []string {
	return r.file.GetSheetList()
}

// GetRows 获取指定工作表的所有行
// sheetName: 工作表名称,如果为空则使用第一个工作表
func (r *XLSXReader) GetRows(sheetName string) ([][]string, error) {
	if sheetName == "" {
		sheets := r.file.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("no sheets found in excel file")
		}
		sheetName = sheets[0]
	}

	rows, err := r.file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet %s: %w", sheetName, err)
	}

	return rows, nil
}

// GetCellValue 获取指定单元格的值
func (r *XLSXReader) GetCellValue(sheetName, axis string) (string, error) {
	if sheetName == "" {
		sheets := r.file.GetSheetList()
		if len(sheets) == 0 {
			return "", fmt.Errorf("no sheets found in excel file")
		}
		sheetName = sheets[0]
	}

	value, err := r.file.GetCellValue(sheetName, axis)
	if err != nil {
		return "", fmt.Errorf("failed to get cell value at %s: %w", axis, err)
	}

	return value, nil
}

// RowToMap 将行数据转换为 map (根据表头映射)
// headers: 表头行
// row: 数据行
func RowToMap(headers []string, row []string) map[string]string {
	result := make(map[string]string)
	for i, header := range headers {
		if i < len(row) {
			result[strings.TrimSpace(header)] = strings.TrimSpace(row[i])
		} else {
			result[strings.TrimSpace(header)] = ""
		}
	}
	return result
}

// ParseFloat 解析浮点数 (支持千分位逗号)
func ParseFloat(value string) (float64, error) {
	// 移除千分位逗号
	value = strings.ReplaceAll(value, ",", "")
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, nil
	}

	return strconv.ParseFloat(value, 64)
}

// ParseInt 解析整数
func ParseInt(value string) (int64, error) {
	value = strings.ReplaceAll(value, ",", "")
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, nil
	}

	return strconv.ParseInt(value, 10, 64)
}

// ParseDate 解析日期 (支持多种格式)
func ParseDate(value string, layouts ...string) (time.Time, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.Time{}, fmt.Errorf("empty date value")
	}

	// 默认支持的日期格式
	defaultLayouts := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02 15:04",
		"2006/01/02 15:04",
	}

	// 如果提供了自定义格式,优先使用
	if len(layouts) > 0 {
		defaultLayouts = append(layouts, defaultLayouts...)
	}

	for _, layout := range defaultLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse date: %s", value)
}

// ParseDateTime 解析日期时间 (将日期和时间字符串合并)
func ParseDateTime(dateStr, timeStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	timeStr = strings.TrimSpace(timeStr)

	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date value")
	}

	// 如果没有时间,默认为 00:00:00
	if timeStr == "" {
		timeStr = "00:00:00"
	}

	// 尝试多种组合格式
	combinations := []string{
		dateStr + " " + timeStr,
		dateStr + "T" + timeStr,
	}

	for _, combined := range combinations {
		if t, err := ParseDate(combined); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse datetime: date=%s, time=%s", dateStr, timeStr)
}

// FindHeaderRow 查找表头行 (包含指定关键字的行)
// rows: 所有行
// keywords: 表头关键字 (至少匹配一个)
// maxSearchRows: 最多搜索的行数
func FindHeaderRow(rows [][]string, keywords []string, maxSearchRows int) int {
	if maxSearchRows == 0 {
		maxSearchRows = 20 // 默认搜索前 20 行
	}

	searchLimit := maxSearchRows
	if searchLimit > len(rows) {
		searchLimit = len(rows)
	}

	for i := 0; i < searchLimit; i++ {
		row := rows[i]
		// 将行内容转为字符串
		rowStr := strings.Join(row, " ")

		// 检查是否包含任一关键字
		for _, keyword := range keywords {
			if strings.Contains(rowStr, keyword) {
				return i
			}
		}
	}

	return -1
}

// GetColumnIndex 获取列索引 (根据表头名称或别名)
// headers: 表头行
// columnNames: 列名或别名列表 (匹配任一即可)
func GetColumnIndex(headers []string, columnNames ...string) int {
	for i, header := range headers {
		headerClean := strings.TrimSpace(header)
		for _, name := range columnNames {
			if strings.Contains(headerClean, name) {
				return i
			}
		}
	}
	return -1
}

// ExtractValue 从行中提取指定列的值
func ExtractValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}
