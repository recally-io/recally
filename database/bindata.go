// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// database/migrations/000001_create_cache_table.down.sql
// database/migrations/000001_create_cache_table.up.sql
// database/migrations/000002_create_users_database.down.sql
// database/migrations/000002_create_users_database.up.sql
// database/migrations/000003_create_assistant_threads.down.sql
// database/migrations/000003_create_assistant_threads.up.sql
package migrations

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

var __000001_create_cache_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4e\x4c\xce\x48\xb5\xe6\x02\x04\x00\x00\xff\xff\xf4\x50\x95\xa6\x1d\x00\x00\x00")

func _000001_create_cache_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_create_cache_tableDownSql,
		"000001_create_cache_table.down.sql",
	)
}

func _000001_create_cache_tableDownSql() (*asset, error) {
	bytes, err := _000001_create_cache_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_create_cache_table.down.sql", size: 29, mode: os.FileMode(420), modTime: time.Unix(1721380493, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000001_create_cache_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x90\xbd\x4e\xc3\x30\x14\x85\xf7\x3c\xc5\x19\x1b\xa9\x5d\x90\x3a\x75\x72\xd3\x5b\x61\x70\x9c\xe2\x1f\xd4\x4e\x96\x89\x2d\x11\x51\x9a\x2a\x75\x10\xbc\x3d\x6a\x42\x07\x32\x20\xb1\xde\xef\x9c\x2b\x9d\x6f\xb1\x40\xd1\x45\x9f\x22\x3c\x92\x7f\x39\x46\xa4\x16\x97\xd4\x76\x11\xb5\xaf\x5f\x23\x82\x4f\x3e\x2b\x14\x31\x43\x30\x6c\x2d\x08\x7c\x0b\x59\x19\xd0\x9e\x6b\xa3\x7f\x52\xb3\x0c\x00\x9a\x00\x4d\x8a\x33\x81\x9d\xe2\x25\x53\x07\x3c\xd2\x61\x3e\xa0\xd0\xbe\xfb\xe6\x84\x67\xa6\x8a\x7b\xa6\x66\x77\xcb\x65\x3e\x7c\x91\x56\x88\x31\xf1\x16\xbf\xfe\xc2\x1f\xfe\xd8\x47\x3c\xe8\x4a\xae\x27\x24\x7e\x9e\x9b\x2e\x5e\x9c\x4f\x30\xbc\x24\x6d\x58\xb9\x1b\x49\x3d\x2c\x0b\xbf\x08\x36\xb4\x65\x56\x18\x14\x56\x29\x92\xc6\x4d\x3a\xfd\x39\xfc\xa3\x93\xe5\xab\xec\x66\xc7\x4a\xfe\x64\x09\x5c\x6e\x68\x3f\x91\xd4\x9f\x1a\x37\x88\x72\xa3\x07\x77\x1d\x5b\xc9\x9b\xbc\xf1\x38\xbf\x2a\xc8\x57\xd9\x77\x00\x00\x00\xff\xff\xcc\xaf\x06\x66\x93\x01\x00\x00")

func _000001_create_cache_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_create_cache_tableUpSql,
		"000001_create_cache_table.up.sql",
	)
}

func _000001_create_cache_tableUpSql() (*asset, error) {
	bytes, err := _000001_create_cache_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_create_cache_table.up.sql", size: 403, mode: os.FileMode(420), modTime: time.Unix(1725002409, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_create_users_databaseDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x2d\x4e\x2d\x2a\xb6\xe6\x02\x04\x00\x00\xff\xff\x2c\x02\x3d\xa7\x1c\x00\x00\x00")

func _000002_create_users_databaseDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_create_users_databaseDownSql,
		"000002_create_users_database.down.sql",
	)
}

func _000002_create_users_databaseDownSql() (*asset, error) {
	bytes, err := _000002_create_users_databaseDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_create_users_database.down.sql", size: 28, mode: os.FileMode(420), modTime: time.Unix(1725002009, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_create_users_databaseUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\x31\x4f\xc3\x30\x10\x85\xf7\xfe\x8a\xdb\xda\x4a\x4c\x48\x9d\x98\x4c\xeb\xaa\x16\x69\x5a\x1c\x1b\x5a\x16\xeb\xa8\x4f\x89\xa5\xc6\xa9\x6c\x07\xfe\x3e\x22\x89\x90\x02\x94\x89\xf1\xee\x7b\xef\x9d\xee\x2d\x25\x67\x8a\x83\x62\xf7\x19\x07\xb1\x86\x7c\xa7\x80\x1f\x44\xa1\x0a\x68\x23\x85\x08\xb3\x09\x00\x80\xb3\x50\x70\x29\x58\x06\x7b\x29\xb6\x4c\x1e\xe1\x81\x1f\x6f\x3a\xd4\xb6\xce\x82\xd6\x62\x05\x3a\x17\x8f\x9a\x77\x11\xb9\xce\x32\x58\xf1\x35\xd3\x99\x82\x92\xbc\x09\xe8\x6d\x53\x9b\x4f\xf1\x6c\x3e\x18\x23\x05\x8f\x35\xc1\x13\x93\xcb\x0d\x93\xb3\xdb\xc5\x62\x40\x17\x8c\xf1\xbd\x09\xd6\x54\x18\x2b\x50\xfc\xa0\xfa\x3d\xd5\xe8\xce\x23\xfd\x70\xb4\xc7\xa5\x4b\x55\xfb\xfa\x07\x6f\x9a\xf2\x4c\xd7\x79\xa2\x33\x95\x01\xeb\xeb\x0a\x3c\x25\xf7\x86\x89\x0c\xc6\xe8\x62\x42\x9f\xcc\xf0\xfd\x37\x9e\xaa\x40\x68\xc7\x30\x26\x4c\x6d\x1c\x87\xff\x28\x6b\x7a\x21\x6f\x9d\x2f\xa7\xbd\xe7\x14\x08\x13\x59\x83\x09\x94\xd8\xf2\x42\xb1\xed\x1e\x9e\x85\xda\x74\x23\xbc\xec\xf2\x5f\x0a\x5f\x6a\x29\x79\xae\xcc\x97\x63\x28\xfc\x62\xff\x21\x6b\x32\xbf\x9b\x7c\x04\x00\x00\xff\xff\xfc\x87\x1c\xb6\x36\x02\x00\x00")

func _000002_create_users_databaseUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_create_users_databaseUpSql,
		"000002_create_users_database.up.sql",
	)
}

func _000002_create_users_databaseUpSql() (*asset, error) {
	bytes, err := _000002_create_users_databaseUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_create_users_database.up.sql", size: 566, mode: os.FileMode(420), modTime: time.Unix(1725002011, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_create_assistant_threadsDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\x2c\x2e\xce\x2c\x2e\x49\xcc\x2b\x89\x4f\xcd\x4d\x4a\x4d\x49\x49\xc9\xcc\x4b\x2f\xb6\xe6\xc2\xaa\x22\xb1\xa4\x24\x31\x39\x23\x37\x35\xaf\x04\x97\x8a\xdc\xd4\xe2\xe2\xc4\xf4\x54\x5c\xd2\x25\x19\x45\xa9\x89\x29\x38\x64\x61\xc2\xae\x11\x21\xae\x7e\xc1\x9e\xfe\x7e\x0a\x9e\x6e\x0a\xae\x11\x9e\xc1\x21\xc1\x0a\x65\xa9\xc9\x25\xf9\x45\xd6\x5c\x80\x00\x00\x00\xff\xff\x68\x55\x1c\x88\xb9\x00\x00\x00")

func _000003_create_assistant_threadsDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_create_assistant_threadsDownSql,
		"000003_create_assistant_threads.down.sql",
	)
}

func _000003_create_assistant_threadsDownSql() (*asset, error) {
	bytes, err := _000003_create_assistant_threadsDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_create_assistant_threads.down.sql", size: 185, mode: os.FileMode(420), modTime: time.Unix(1725003095, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_create_assistant_threadsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x97\x4b\x53\xdb\x30\x10\xc7\xef\xfe\x14\x7b\x23\x9e\x81\x03\x30\xe9\x81\x9c\x4c\x10\xe0\x36\x28\xa9\x2c\xd3\xd0\x8b\x47\x44\x3b\xc1\xd3\xf8\x31\x96\xd2\x09\xed\xf4\xbb\x77\xe2\xb7\x4d\x5e\xa5\x4d\x87\x4e\x7a\xf4\x6a\x1f\xd6\xee\xff\x27\xd9\x27\x27\xd0\x4f\x50\x68\x84\xd1\xcd\x3d\x4e\x74\x94\x00\x59\x68\x0c\x95\x1f\x85\x46\x9f\x11\x8b\x13\x20\x63\x4e\xa8\x63\x0f\x29\xd8\xd7\x40\x87\x1c\xc8\xd8\x76\xb8\x03\x5f\x53\xff\x9e\x61\x54\x49\xb8\x78\x9c\xa1\x02\x4b\x29\x5f\x69\x11\xea\x22\x05\xb7\x2e\x07\xa4\x15\x2e\x0a\x27\x05\x1d\x03\x00\xc0\x97\xe0\x10\x66\x5b\x03\x18\x31\xfb\xce\x62\x0f\xf0\x81\x3c\x1c\xa7\x4b\xf3\xb9\x2f\xc1\x75\xed\x2b\x70\xa9\xfd\xd1\x25\x69\x1e\xea\x0e\x06\x70\x45\xae\x2d\x77\xc0\x61\x8a\xa1\x97\x88\x50\x46\x81\xb7\x74\xee\x98\x79\xa0\xc2\xc4\x2b\x62\x19\xb9\x26\x8c\xd0\x3e\x71\x52\xbb\xea\x2c\x3d\x73\xc7\x50\x04\x08\xf7\x16\xeb\xdf\x5a\xac\x73\xd6\xed\x9a\x65\x89\x6c\x5d\xa2\x9a\x24\x7e\xac\xfd\x28\x04\x4e\xc6\x3c\xb3\xaa\x67\xa5\x31\xf0\xe2\x24\x0a\x62\x5d\xb3\x07\x91\xc4\x59\x99\xee\xfc\xac\x9d\x2d\x40\x2d\xa4\xd0\x02\xde\x3b\x43\x7a\x59\x6e\xe2\xe8\xfb\x8f\xa3\x8b\x8b\xd4\x96\xf9\x4d\xd2\xb6\x4a\x4f\x68\xe0\xf6\x1d\x71\xb8\x75\x37\x82\x4f\x36\xbf\x4d\x1f\xe1\xf3\x90\xae\x68\x45\xdf\x65\x8c\x50\xee\x95\x11\x79\x2b\x62\xf9\x07\x72\x19\x66\xaf\x18\xaa\x4d\xaf\xc8\xb8\x35\x54\x5f\x2e\xbc\xb4\xe7\xb5\x37\x1f\xd2\xc6\xac\xf3\x91\x1c\xd7\x36\x67\xf6\x0c\x63\x27\xa9\x78\xfa\x29\x41\x21\xdf\x84\x62\xaa\x97\x5a\xe1\x5d\x6d\xf8\xbf\xc8\xf6\x27\xb2\x6a\x04\x6b\xe4\x56\xe9\xa5\x54\x5d\x7d\x6c\x2d\x0d\x6e\x29\xf9\x0b\xd5\x36\x14\xd9\x51\xe8\x01\x2a\x25\xa6\xf8\xaf\x2a\x3d\xeb\xc4\x46\xff\xa2\x5b\xf5\xb0\x17\x92\xce\xcc\x49\x34\xdb\xc8\x8d\xc6\x45\x1d\x8c\x8c\x14\x4f\x47\x5f\x30\x04\x9b\x72\x72\x43\x58\x2e\xf5\x28\x88\x67\xb8\xe4\x6b\xd5\x2a\x06\x8f\x28\xa5\x1f\x4e\x55\x7e\xb7\x41\xe7\xb4\x7b\xfe\xce\x3c\x14\x9c\xd6\xc9\xba\xd2\xe2\xea\xb3\xfb\xb7\xb9\xa9\x0a\xbc\x9e\xce\x5c\x71\xdb\x4b\x94\xd2\x7c\x1d\x98\x42\x6b\x31\x79\x0a\xf0\x8d\x7c\xb7\xfc\x0d\x36\x6b\xa4\xbd\xb8\xc1\xf2\x94\xcf\xf1\x2a\xf3\x3c\xa9\x60\xee\x9e\x16\x34\x2b\xff\x1b\x36\xc1\x3b\x5c\xb4\x1a\x6a\xda\x13\x5d\x8d\x1a\x7b\x03\xac\x51\x65\x3d\x63\xb5\x7f\x85\xe5\xa1\x4d\xaa\x23\x37\xfd\x77\xd8\x8d\xc1\xec\xa0\xce\xc2\xb6\x33\xb8\x33\x4a\xe5\x0e\x0a\xf7\xd6\xed\xd2\xba\x75\xd6\x5d\x17\x87\xf6\x19\xd6\xd4\x41\x63\x36\x79\xef\x77\x50\x71\xd9\xfb\x0d\xd9\x1a\x03\x32\x7b\xc6\xcf\x00\x00\x00\xff\xff\xf4\x4f\xd6\x1d\xbd\x0e\x00\x00")

func _000003_create_assistant_threadsUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_create_assistant_threadsUpSql,
		"000003_create_assistant_threads.up.sql",
	)
}

func _000003_create_assistant_threadsUpSql() (*asset, error) {
	bytes, err := _000003_create_assistant_threadsUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_create_assistant_threads.up.sql", size: 3773, mode: os.FileMode(420), modTime: time.Unix(1725085864, 0)}
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
	"000001_create_cache_table.down.sql":       _000001_create_cache_tableDownSql,
	"000001_create_cache_table.up.sql":         _000001_create_cache_tableUpSql,
	"000002_create_users_database.down.sql":    _000002_create_users_databaseDownSql,
	"000002_create_users_database.up.sql":      _000002_create_users_databaseUpSql,
	"000003_create_assistant_threads.down.sql": _000003_create_assistant_threadsDownSql,
	"000003_create_assistant_threads.up.sql":   _000003_create_assistant_threadsUpSql,
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
	"000001_create_cache_table.down.sql":       &bintree{_000001_create_cache_tableDownSql, map[string]*bintree{}},
	"000001_create_cache_table.up.sql":         &bintree{_000001_create_cache_tableUpSql, map[string]*bintree{}},
	"000002_create_users_database.down.sql":    &bintree{_000002_create_users_databaseDownSql, map[string]*bintree{}},
	"000002_create_users_database.up.sql":      &bintree{_000002_create_users_databaseUpSql, map[string]*bintree{}},
	"000003_create_assistant_threads.down.sql": &bintree{_000003_create_assistant_threadsDownSql, map[string]*bintree{}},
	"000003_create_assistant_threads.up.sql":   &bintree{_000003_create_assistant_threadsUpSql, map[string]*bintree{}},
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
