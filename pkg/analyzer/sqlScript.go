// Package analyzer Code generated by go-bindata. (@generated) DO NOT EDIT.
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
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

var _createDbSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xd1\xb1\x4e\xc3\x40\x0c\x06\xe0\xfd\x9e\xc2\x63\x2b\xb1\xd0\x8e\x11\x0b\x08\xa1\x0e\x85\x85\xfd\xe4\xe6\xac\x62\xc8\xf9\x22\xc7\x01\xfa\xf6\x28\xc9\xa5\x28\x34\x01\xb2\x25\xf7\xfd\xba\xd8\xff\xed\xfd\xc3\xee\xb1\x70\xae\x54\x42\x23\x30\x3c\x54\x04\x4a\x75\x52\x73\x2b\x07\x00\xc0\x01\xa6\x4f\x63\xca\x72\x04\x49\x06\xd2\x56\x15\xd4\xca\x11\xf5\x04\x6f\x74\xba\xea\x13\x11\xcb\x17\x16\x5a\x4c\x0c\xea\x48\x42\x8a\x46\xc1\x1b\xc7\x0e\xb3\x58\xf7\x7d\xaa\x22\x99\x72\xd9\xf8\x77\xd2\x86\x93\xcc\xa8\x9e\xd5\x9a\x42\x5b\xda\x1f\x57\x1e\x5a\xae\x82\x2f\xaf\xcf\x6a\xf6\xca\xac\x36\xff\x52\xdb\x65\xd5\xb3\xd0\x2a\x1a\x27\xf1\x79\x90\xf9\x1f\x63\x69\x0c\xc5\xce\x68\x5e\x29\x7e\xf8\xa1\x9a\xf9\x21\xdd\xfa\xbb\x48\x96\x40\x9f\x63\x13\x7e\x78\x4b\x92\x9b\x85\x55\x3e\x58\x17\x53\x9f\xd7\x78\xe9\xf3\xc1\x4f\x3f\x2c\x21\xe2\x6b\xd2\xcb\xcc\xb8\xed\x85\x10\xcb\x2f\xa1\x4d\x37\x4a\xad\x78\x8c\x08\x6d\x43\x3a\xf6\x7f\xb3\x2d\x9c\xbb\x7b\xda\xef\x77\xcf\xc5\x57\x00\x00\x00\xff\xff\x3a\xa9\x0d\xc7\xbb\x02\x00\x00")

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

	info := bindataFileInfo{name: "create-db.sql", size: 699, mode: os.FileMode(420), modTime: time.Unix(1569950724, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _insertReportSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8f\xcb\x0a\xc2\x40\x0c\x45\xf7\xfd\x8a\x2c\x5b\xc8\x46\xfd\x00\x29\xd2\x85\x50\x54\x7c\x6d\x87\x71\x26\x68\xc0\xce\x94\x34\xd5\xdf\x17\x7c\x0c\xa8\x6d\xc8\xe2\x70\x0f\x37\x90\x6d\xb5\xa9\xcb\x45\x05\xcb\xd5\x7e\x0d\x42\x6d\x14\x85\x9c\x3d\x42\x63\xdd\x85\x03\x21\x9c\x29\x90\x58\x25\x6f\x94\x1b\xc2\x0c\x86\xa6\x95\xe8\x7b\xa7\x23\xf6\xd4\xf3\xd5\x1b\x37\xc1\x0f\x4d\x13\xcd\x46\x2a\x0d\xa9\xb0\xeb\xcc\x8d\xa4\xe3\x18\x10\x7c\x2f\x56\x39\x06\xf3\x36\x08\x1c\x3a\xb5\x41\x53\x30\x7c\x48\xec\xdd\xbc\xfe\x2a\xb2\x63\x59\x1f\xaa\x1d\xe4\x73\x84\xe7\xa6\xc6\x17\xfe\xb9\xdf\xa0\x78\x04\x00\x00\xff\xff\x69\x08\xf9\x15\x35\x01\x00\x00")

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

	info := bindataFileInfo{name: "insert-report.sql", size: 309, mode: os.FileMode(420), modTime: time.Unix(1569950724, 0)}
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
