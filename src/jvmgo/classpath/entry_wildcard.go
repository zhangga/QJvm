package classpath

import (
	"os"
	"path/filepath"
	"strings"
)

// 通配符类路径不能递归匹配子目录下的JAR文件
func newWildcardEntry(path string) CompositeEntry {
	// remove *
	baseDir := path[:len(path)-1]
	compositeEntry := []Entry{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			// 跳过子目录
			return filepath.SkipDir
		}
		// 匹配JAR文件
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	filepath.Walk(baseDir, walkFn)
	return compositeEntry
}
