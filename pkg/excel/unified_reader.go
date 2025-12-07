package excel

import (
	"bytes"
	"fmt"
	"io"
)

// ExcelReader 统一的 Excel 读取接口 (支持 xls 和 xlsx)
type ExcelReader interface {
	GetRows(sheetName string) ([][]string, error)
	GetSheetNames() []string
	Close() error
}

// FileFormat Excel 文件格式
type FileFormat int

const (
	FormatUnknown FileFormat = iota
	FormatXLS                // Excel 97-2003 (.xls)
	FormatXLSX               // Excel 2007+ (.xlsx)
)

// 文件头魔数
var (
	// XLS 文件头 (OLE2 格式)
	xlsMagic = []byte{0xD0, 0xCF, 0x11, 0xE0}
	// XLSX 文件头 (ZIP 格式)
	xlsxMagic = []byte{0x50, 0x4B, 0x03, 0x04}
)

// NewUnifiedReader 创建统一的 Excel 读取器 (自动检测格式)
func NewUnifiedReader(reader io.ReadSeeker) (ExcelReader, error) {
	// 检测文件格式
	format, err := detectFormat(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file format: %w", err)
	}

	// 重置读取位置
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to reset reader: %w", err)
	}

	// 根据格式创建对应的读取器
	switch format {
	case FormatXLS:
		return NewXLSReader(reader)
	case FormatXLSX:
		return NewXLSXReader(reader)
	default:
		return nil, fmt.Errorf("unsupported file format (not xls or xlsx)")
	}
}

// detectFormat 检测 Excel 文件格式
func detectFormat(reader io.ReadSeeker) (FileFormat, error) {
	// 读取文件头 (前 4 字节)
	header := make([]byte, 4)
	n, err := io.ReadFull(reader, header)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return FormatUnknown, fmt.Errorf("file too small to detect format")
		}
		return FormatUnknown, fmt.Errorf("failed to read file header: %w", err)
	}

	if n < 4 {
		return FormatUnknown, fmt.Errorf("file too small to detect format")
	}

	// 检查 XLSX 格式 (ZIP 文件头)
	if bytes.Equal(header, xlsxMagic) {
		return FormatXLSX, nil
	}

	// 检查 XLS 格式 (OLE2 文件头)
	if bytes.Equal(header, xlsMagic) {
		return FormatXLS, nil
	}

	return FormatUnknown, nil
}
