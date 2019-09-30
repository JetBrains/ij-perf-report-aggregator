// Code generated for package analyzer by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/analyzer/sql/create-db.sql
// pkg/analyzer/sql/insert-report.sql
package analyzer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _createDbSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x91\x31\x4f\xc3\x30\x10\x85\x77\xff\x8a\x1b\x5b\x89\x85\x32\x46\x2c\x20\x84\x3a\x14\x16\x76\xeb\x6a\x9f\xc2\x41\x7c\x8e\x2e\x17\xa0\xff\x1e\xb5\x71\x04\x21\x6d\x84\x37\xfb\xbe\xa7\x7b\x7e\xef\xee\xe1\x71\xfb\x54\x39\x17\x94\xd0\x08\x0c\xf7\x0d\x81\x52\x9b\xd5\xdc\xca\x01\x00\x70\x84\xc9\xe9\x4c\x59\x6a\x90\x6c\x20\x7d\xd3\x40\xab\x9c\x50\x0f\xf0\x4e\x87\xab\x93\x20\x61\x78\x65\xa1\x4b\x82\x01\xaa\x49\x48\xd1\x28\x7a\xe3\x44\x00\x2c\x76\x7c\x9e\x42\x89\x4c\x39\x74\xfe\x83\xb4\xe3\x2c\x73\xe8\x44\xb5\x9a\x63\x1f\x6c\x79\xdf\xbe\xe7\x26\xfa\x70\x3d\x42\x67\xf7\x15\x68\xf3\x1f\xe8\xe6\x22\xf4\xdb\xfa\xb2\x29\xc5\x4f\x3f\x64\x7d\x16\x72\xeb\x9f\x62\x58\x22\x7d\x8d\xd1\xfa\xe1\x96\xa5\x34\x05\xab\x32\x58\x57\x53\xbe\x44\x33\xe7\xcb\xe0\x2f\x3f\xfc\x2d\xe1\x5b\xd6\xb9\x66\x8c\xf0\xe8\xaa\x55\xac\x13\x42\xdf\x91\x8e\xed\xdc\x6e\x2a\xe7\xee\x9f\x77\xbb\xed\x4b\xf5\x1d\x00\x00\xff\xff\x98\x14\xcb\x5c\x56\x02\x00\x00")

func createDbSqlBytes() ([]byte, error) {
	return bindataRead(
		_createDbSql,
		"create-db.sql",
	)
}

func createDbSql() (*asset, error) {
	bytes, err := createDbSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "create-db.sql", size: 598, mode: os.FileMode(420), modTime: time.Unix(1569831191, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _insertReportSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\x72\x0d\xf0\x71\x74\x76\x55\xf0\xf4\x0b\xf1\x57\x28\x4a\x2d\xc8\x2f\x2a\x51\xd0\xc8\x4c\xd1\x51\xc8\x4d\x4c\xce\xc8\xcc\x4b\xd5\x51\x48\x4f\xcd\x4b\x2d\x4a\x2c\x49\x4d\x89\x2f\xc9\xcc\x4d\xd5\xe1\x52\xc0\x06\x0a\x8a\xf2\x53\x4a\x93\x4b\x70\xc8\x26\x95\x66\xe6\xa4\xc4\x27\x1b\xea\xc0\x58\x46\x70\x96\x31\x0e\x2d\xb9\xa9\x25\x45\x99\xc9\xc5\xf1\x65\xa9\x45\xc5\x99\xf9\x79\x3a\x30\x01\x1c\xca\x8b\x12\xcb\xe3\x21\xae\xd7\xe4\x0a\x73\xf4\x09\x75\x0d\x56\xd0\xb0\xd7\x51\x00\x23\xb8\x0e\x14\x26\x86\x1c\x0a\x4f\x13\x10\x00\x00\xff\xff\xe7\x95\xe0\xe2\x18\x01\x00\x00")

func insertReportSqlBytes() ([]byte, error) {
	return bindataRead(
		_insertReportSql,
		"insert-report.sql",
	)
}

func insertReportSql() (*asset, error) {
	bytes, err := insertReportSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "insert-report.sql", size: 280, mode: os.FileMode(420), modTime: time.Unix(1569831191, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"create-db.sql":     createDbSql,
	"insert-report.sql": insertReportSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"create-db.sql":     &bintree{createDbSql, map[string]*bintree{}},
	"insert-report.sql": &bintree{insertReportSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
