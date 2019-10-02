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

var _createDbSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x31\x4f\xc3\x30\x10\x85\x77\xff\x8a\x1b\x5b\x89\x85\x96\xad\x62\x63\x61\x42\x62\x61\x8c\xae\xf6\x35\x3d\x88\xcf\xd1\xe5\x02\xf4\xdf\x23\x53\xa7\xa5\x6d\x12\x91\x25\x8a\xf2\xbd\xf3\xf3\x7b\xb7\xa5\x9a\x65\xe3\x9c\x57\x42\x23\x30\xdc\x36\x04\x11\xfd\x9e\x85\xdc\xc2\x01\x08\x46\x02\x0e\xd0\x99\xb2\xd4\x20\xc9\x40\xfa\xa6\x81\x56\x39\xa2\x1e\xe0\x83\x0e\x6e\x79\x3d\x40\xa9\x4d\x6a\xbf\x7a\x0e\x70\xf9\xcc\x0c\xba\x73\x30\x9c\x7d\xe6\x59\x2c\xbf\x06\x3e\x33\x35\x09\x29\x1a\x85\xca\x38\xd2\x38\x13\xc9\x94\x7d\x57\x7d\x92\x76\x9c\x64\x84\x71\x00\xad\xa6\xd0\x7b\x9b\x34\x97\x07\x6d\x7b\x6e\x42\xe5\xef\xe7\x0c\x15\x66\xf5\x0f\x66\x3d\xcd\x38\x80\xd0\x2b\x1a\x27\xa9\x8a\xfd\x31\x43\x2c\x9d\xa1\xd8\x09\x19\x63\x14\xbf\xaa\x63\x09\x53\x17\x73\x00\xbb\xa4\xc4\xb5\xe4\xe4\x61\x51\x72\x5f\x82\xd2\x8e\x94\xc4\x53\x77\xea\x62\xf1\xfa\xf2\xf6\xfc\xb4\x84\x24\x10\xa8\x21\xcb\x05\xe7\x81\xde\xfe\x56\xcf\x12\xe8\x7b\xd0\x54\xc7\xaf\x24\x65\x17\xce\x07\x6c\x2e\xf9\x52\xc1\x2d\x5f\x7e\x5c\xf3\xc7\x18\x23\xbe\x27\xbd\xd5\x0c\x5d\x4d\x88\x58\x66\x44\xab\x71\x51\x8b\xe6\xf7\x93\xa2\x75\xbe\x7f\xab\x58\x47\x84\xbe\x23\x1d\xd6\xed\xf1\x21\xe7\x92\x62\x64\xdb\xfc\x04\x00\x00\xff\xff\x00\xfd\x28\x5b\x63\x03\x00\x00")

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

	info := bindataFileInfo{name: "create-db.sql", size: 867, mode: os.FileMode(420), modTime: time.Unix(1570032475, 0)}
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

	info := bindataFileInfo{name: "insert-report.sql", size: 309, mode: os.FileMode(420), modTime: time.Unix(1570013904, 0)}
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
