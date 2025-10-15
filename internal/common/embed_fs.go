package common

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
)

// embedFileSystem 实现 static.ServeFileSystem 接口
type embedFileSystem struct {
	http.FileSystem
}

// Exists 检查文件是否存在
func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

// EmbedFolder 创建嵌入式文件系统
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(efs),
	}
}
