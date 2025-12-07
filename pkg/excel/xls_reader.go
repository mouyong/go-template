package excel

import (
	"fmt"
	"io"

	"github.com/extrame/xls"
)

// XLSReader XLS 格式读取器 (Excel 97-2003)
type XLSReader struct {
	workbook *xls.WorkBook
}

// NewXLSReader 创建 XLS 读取器
func NewXLSReader(reader io.ReadSeeker) (*XLSReader, error) {
	workbook, err := xls.OpenReader(reader, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("failed to open xls file: %w", err)
	}

	return &XLSReader{workbook: workbook}, nil
}

// Close 关闭 XLS 文件 (XLS 库不需要显式关闭)
func (r *XLSReader) Close() error {
	return nil
}

// GetSheetNames 获取所有工作表名称
func (r *XLSReader) GetSheetNames() []string {
	count := r.workbook.NumSheets()
	names := make([]string, 0, count)

	for i := 0; i < count; i++ {
		if sheet := r.workbook.GetSheet(i); sheet != nil {
			names = append(names, sheet.Name)
		}
	}

	return names
}

// GetRows 获取指定工作表的所有行
// sheetName: 工作表名称, 如果为空则使用第一个工作表
func (r *XLSReader) GetRows(sheetName string) ([][]string, error) {
	var sheet *xls.WorkSheet

	// 如果未指定工作表名称, 使用第一个工作表
	if sheetName == "" {
		if r.workbook.NumSheets() == 0 {
			return nil, fmt.Errorf("no sheets found in xls file")
		}
		sheet = r.workbook.GetSheet(0)
	} else {
		// 根据名称查找工作表
		for i := 0; i < r.workbook.NumSheets(); i++ {
			s := r.workbook.GetSheet(i)
			if s != nil && s.Name == sheetName {
				sheet = s
				break
			}
		}
	}

	if sheet == nil {
		return nil, fmt.Errorf("sheet not found: %s", sheetName)
	}

	// 转换为 [][]string 格式
	rows := make([][]string, 0, int(sheet.MaxRow)+1)

	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			// 空行也要保留 (保持与 xlsx 一致)
			rows = append(rows, []string{})
			continue
		}

		// 获取该行的列数据
		cols := make([]string, 0)
		lastCol := row.LastCol()

		for j := row.FirstCol(); j < lastCol; j++ {
			cellValue := row.Col(j)
			cols = append(cols, cellValue)
		}

		rows = append(rows, cols)
	}

	return rows, nil
}

// GetCellValue 获取指定单元格的值
func (r *XLSReader) GetCellValue(sheetName, axis string) (string, error) {
	var sheet *xls.WorkSheet

	// 如果未指定工作表名称, 使用第一个工作表
	if sheetName == "" {
		if r.workbook.NumSheets() == 0 {
			return "", fmt.Errorf("no sheets found in xls file")
		}
		sheet = r.workbook.GetSheet(0)
	} else {
		// 根据名称查找工作表
		for i := 0; i < r.workbook.NumSheets(); i++ {
			s := r.workbook.GetSheet(i)
			if s != nil && s.Name == sheetName {
				sheet = s
				break
			}
		}
	}

	if sheet == nil {
		return "", fmt.Errorf("sheet not found: %s", sheetName)
	}

	// 解析坐标 (如 A1, B2)
	col, row, err := parseAxis(axis)
	if err != nil {
		return "", fmt.Errorf("invalid axis %s: %w", axis, err)
	}

	// 检查行范围
	if row > int(sheet.MaxRow) {
		return "", nil
	}

	rowData := sheet.Row(row)
	if rowData == nil {
		return "", nil
	}

	// 检查列范围
	if col < rowData.FirstCol() || col >= rowData.LastCol() {
		return "", nil
	}

	return rowData.Col(col), nil
}

// parseAxis 解析单元格坐标 (如 A1 -> col=0, row=0)
func parseAxis(axis string) (col int, row int, err error) {
	if len(axis) < 2 {
		return 0, 0, fmt.Errorf("invalid axis format")
	}

	// 解析列 (A-Z, AA-ZZ)
	colStr := ""
	rowStr := ""
	for _, ch := range axis {
		if ch >= 'A' && ch <= 'Z' {
			colStr += string(ch)
		} else if ch >= 'a' && ch <= 'z' {
			colStr += string(ch - 32) // 转大写
		} else if ch >= '0' && ch <= '9' {
			rowStr += string(ch)
		}
	}

	if colStr == "" || rowStr == "" {
		return 0, 0, fmt.Errorf("invalid axis format")
	}

	// 列名转索引 (A=0, B=1, ..., Z=25, AA=26)
	col = 0
	for i := 0; i < len(colStr); i++ {
		col = col*26 + int(colStr[i]-'A') + 1
	}
	col--

	// 行号转索引 (1-based -> 0-based)
	_, err = fmt.Sscanf(rowStr, "%d", &row)
	if err != nil {
		return 0, 0, err
	}
	row--

	return col, row, nil
}
