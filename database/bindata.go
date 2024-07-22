// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// database/migrations/000001_new_cache_table.down.sql
// database/migrations/000001_new_cache_table.up.sql
// database/migrations/000002_river.down.sql
// database/migrations/000002_river.up.sql
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

var __000001_new_cache_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4e\x4c\xce\x48\xb5\xe6\x02\x04\x00\x00\xff\xff\xf4\x50\x95\xa6\x1d\x00\x00\x00")

func _000001_new_cache_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_new_cache_tableDownSql,
		"000001_new_cache_table.down.sql",
	)
}

func _000001_new_cache_tableDownSql() (*asset, error) {
	bytes, err := _000001_new_cache_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_new_cache_table.down.sql", size: 29, mode: os.FileMode(420), modTime: time.Unix(1721380493, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000001_new_cache_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x90\xbd\x4e\xc3\x30\x14\x85\xf7\x3c\xc5\x19\x1b\xa9\x5d\x90\x3a\x75\x72\xd3\x5b\x61\x70\x9c\xe2\x1f\xd4\x4e\x96\x89\x2d\x11\x51\x9a\x2a\x75\x10\xbc\x3d\x6a\x42\x07\x32\x20\xb1\xde\xef\x9c\x2b\x9d\x6f\xb1\x40\xd1\x45\x9f\x22\x3c\x92\x7f\x39\x46\xa4\x16\x97\xd4\x76\x11\xb5\xaf\x5f\x23\x82\x4f\x3e\x2b\x14\x31\x43\x30\x6c\x2d\x08\x7c\x0b\x59\x19\xd0\x9e\x6b\xa3\x7f\x52\xb3\x0c\x00\x9a\x00\x4d\x8a\x33\x81\x9d\xe2\x25\x53\x07\x3c\xd2\x61\x3e\xa0\xd0\xbe\xfb\xe6\x84\x67\xa6\x8a\x7b\xa6\x66\x77\xcb\x65\x3e\x7c\x91\x56\x88\x31\xf1\x16\xbf\xfe\xc2\x1f\xfe\xd8\x47\x3c\xe8\x4a\xae\x27\x24\x7e\x9e\x9b\x2e\x5e\x9c\x4f\x30\xbc\x24\x6d\x58\xb9\x1b\x49\x3d\x2c\x0b\xbf\x08\x36\xb4\x65\x56\x18\x14\x56\x29\x92\xc6\x4d\x3a\xfd\x39\xfc\xa3\x93\xe5\xab\xec\x66\xc7\x4a\xfe\x64\x09\x5c\x6e\x68\x3f\x91\xd4\x9f\x1a\x37\x88\x72\xa3\x07\x77\x1d\x5b\xc9\x9b\xbc\xf1\x38\xbf\x2a\xc8\x57\xd9\x77\x00\x00\x00\xff\xff\xcc\xaf\x06\x66\x93\x01\x00\x00")

func _000001_new_cache_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_new_cache_tableUpSql,
		"000001_new_cache_table.up.sql",
	)
}

func _000001_new_cache_tableUpSql() (*asset, error) {
	bytes, err := _000001_new_cache_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_new_cache_table.up.sql", size: 403, mode: os.FileMode(420), modTime: time.Unix(1721443399, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_riverDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x54\x5d\x6f\xdb\x36\x14\x7d\xe7\xaf\x38\x0f\x01\x64\x03\x76\x91\x75\x7b\x18\xa2\x75\x80\x2a\xd1\x8e\x30\x55\x32\x68\x09\xdd\x30\x0c\x06\x2d\xdd\x38\xcc\x68\x52\x95\xe8\x64\xd9\xaf\x2f\x24\xd9\x8d\x63\xc7\x40\xab\x27\x91\xf7\x83\xe7\x9e\x73\xef\x9d\x4e\x21\xd4\x23\x35\xd8\xaa\x4d\x23\x9d\xb2\x06\xd7\xd7\xbf\xe0\xef\xca\x3e\x99\x7f\x58\x90\xe4\x5c\x20\x0f\x3e\x26\x1c\x4d\xe7\xb6\x7a\xb0\x6b\x0c\xb7\x61\x96\x14\x9f\x52\xc8\x66\xd3\x22\x12\xd9\x02\x69\x96\x23\x2d\x92\xc4\x67\xdf\x13\xb7\x25\x27\x2b\xe9\xe4\x69\xec\x0f\x87\x46\x7c\x16\x14\x49\xee\x33\x36\x9d\x22\x76\x50\x2d\x8c\x75\xa8\x6d\xdb\xaa\xb5\x26\x38\x8b\x56\xde\x91\x7e\x46\x43\x5b\xfb\x48\xf0\x6a\x32\x95\x32\x1b\x0f\x77\x8d\xdd\xc2\xdd\xd3\xcb\x43\xab\xd6\x49\x47\x20\xb3\xdb\x4e\xba\x7c\xad\x85\x26\xf9\x48\x50\x0e\xca\xa0\xd6\xb2\xa4\x77\x97\xea\xeb\xe1\x84\x59\xba\xcc\x45\x10\xa7\x39\xee\x94\x91\x5a\xfd\x4f\xd5\xca\x36\xab\x97\x83\x74\x2b\xb3\xd3\xfa\x62\xa9\x51\xf4\xfd\x49\x10\xde\xf2\xf0\x0f\x8c\x18\x30\x1a\xa0\xc7\x29\x46\x5e\x29\x4d\x49\x5a\x53\xe5\x4d\xe0\x95\x76\x5b\x6b\x72\xc3\xa1\x52\x6d\x29\x9b\x8a\x2a\x6f\x8c\x20\x8d\x70\x9c\x11\xf1\xf2\x9b\x10\x63\x64\xe2\xdc\x58\x24\x09\x1b\xfb\x8c\x85\x82\x07\x39\xef\x5c\x04\x5f\x24\x41\xc8\x31\x2b\xd2\x30\x8f\xb3\xf4\x88\x4a\x63\x9d\xba\x7b\x1e\x8d\x19\x20\x78\x5e\x88\x74\x89\x5c\xc4\xf3\x39\x17\x0c\x08\x96\xb8\xba\x62\x11\x0f\x93\x40\x70\x06\xd4\xf2\x59\x5b\x59\xe1\xa1\xb5\xc6\x67\x1f\xf9\x3c\x4e\x19\x10\xcf\x90\xf2\xcf\xef\x86\xca\x3e\xc0\x93\x8f\x52\x69\xb9\xd6\xe4\x21\xbf\xe5\x9d\x07\x30\x9d\x22\xed\x5f\xc2\x93\xd2\x1a\xa5\x95\x9a\xda\x92\x50\xed\x6a\xad\xca\x2e\xb0\x07\xa2\x4a\xd5\xb7\x77\x8b\x27\xe5\xee\x95\x81\x84\x6b\xa4\x69\x65\xd9\xdd\x4e\xd0\xda\x43\xb6\x7f\x89\xea\xae\x2b\x5a\x3a\xa0\x6a\xb1\x21\x43\xcd\xc0\xc5\x4d\xef\x77\xc0\xfb\xa1\x47\xbc\x5a\xef\x94\xae\x56\x76\xfd\x40\xa5\x1b\x79\x5f\x76\xb4\x23\x6f\xd2\x63\xef\xff\xc7\x7e\x1f\xb4\xe0\x62\x96\x89\x4f\xfd\x3f\x50\x6f\x0e\x14\x79\x03\x69\xca\xb4\xd4\x38\x6f\x72\x48\x7e\x73\xe3\xe8\x3f\xd7\xc7\xf2\x34\x42\x3c\xf3\xbf\x51\xb9\x1f\x16\x9e\x46\x3e\xbb\xba\x62\x49\x90\xce\x8b\x60\xce\x51\xeb\x7a\xd3\x7e\xd1\x2f\x12\xed\x19\xdf\xcb\x32\xbc\xd7\xd1\x3f\xeb\x5a\x2f\x4e\x97\x5c\xe4\x38\x56\x8d\x01\xb3\x4c\x80\x07\xe1\x2d\x44\xf6\xb9\x7b\xfa\x4f\x1e\x16\x39\xc7\x42\x64\x21\x8f\x0a\xc1\xdf\x90\xd8\x67\xac\x6f\xfe\xe3\x5e\xee\x0b\x7f\x73\x17\x68\x92\x15\x35\x3d\x0b\xaf\xe6\xda\xc8\x2d\xbd\x9a\xe9\x49\xef\x73\x3a\x56\x9d\xdb\x4a\x93\xd9\xb8\xfb\xc1\xe1\x64\x62\x8e\xec\x87\xf1\x28\xef\x65\xb3\xbf\x1a\x75\xe6\x31\x7e\xc7\x75\xdf\xff\xe7\x96\xdf\xf0\xd3\xfb\x5f\xc7\xc3\x3a\x39\x5f\x8d\x3f\xff\xc0\x6a\x74\xf2\x74\x35\x4e\xf6\xd2\x9f\x7e\x17\xc2\x5e\xad\xb6\x73\x2c\xef\x0f\x58\xce\xb8\x7f\xb0\x6b\x7f\xb8\xbd\x38\x97\x7b\x7b\xfe\xd7\x82\x9f\xae\xbf\xb7\xd4\x1c\x34\xf3\xd9\xd7\x00\x00\x00\xff\xff\xc6\x52\xa5\x3d\x30\x06\x00\x00")

func _000002_riverDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_riverDownSql,
		"000002_river.down.sql",
	)
}

func _000002_riverDownSql() (*asset, error) {
	bytes, err := _000002_riverDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_river.down.sql", size: 1584, mode: os.FileMode(420), modTime: time.Unix(1721491008, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_riverUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x58\x51\x8f\xe2\x38\x12\x7e\xcf\xaf\xa8\x87\x91\x80\x3b\x18\xcd\xcc\xcd\x4a\xa3\x61\x7b\x25\x16\xdc\x0c\x5a\x26\xb4\x42\xb8\x9e\xd1\x6a\x45\x3b\x49\x01\xee\x36\x76\xd6\x76\xba\x87\x3d\xdd\x7f\x3f\xd9\x26\x21\xa1\x81\xeb\x7e\xdc\xb7\xc4\x2e\x7f\x55\xae\xaa\x7c\x55\x95\x5e\x0f\x22\xf6\x88\x0a\xb6\x6c\xad\xa8\x61\x52\xc0\xbb\x77\x1f\xe0\xf7\x22\xff\x23\x18\x46\x64\x10\x13\x88\xbf\xdf\x10\x50\x56\x68\x79\x2f\x93\xa5\x36\xd4\x20\x0c\xe6\x40\xc2\xc5\xd7\x76\x00\xd0\xa2\x8f\x94\x71\x9a\x70\x6c\x75\xed\x6b\x4a\x45\x8a\x9c\x63\xb6\x7f\x95\xdb\x9c\xa3\x29\x5f\x33\xa6\x53\xaa\xb2\xf2\x55\xa1\x51\xbb\xc3\x59\x55\x08\xc1\xc4\xda\xbf\xe8\x74\x83\x59\x61\x81\x82\x4e\x3f\xa8\xcc\x19\xfc\x3a\xad\xd9\x63\x2d\xe8\xf5\xe0\x13\x24\x3b\x83\x3a\x00\x60\x19\x24\x6c\xad\x51\x31\xca\xe1\x26\x9a\x7c\x1d\x44\xdf\xe1\x37\xf2\xbd\x1b\x34\x24\xa1\xfd\x71\xff\xf0\x4f\xf8\x70\xfc\xd4\x71\xa2\x5e\xfe\xce\x5d\xf8\x0e\x98\x86\x07\xcc\x0d\x08\xa4\x0a\xcc\x06\xc1\xc8\x1c\xe4\xca\x3f\xda\x1b\xc0\x4a\x2a\x90\x39\x2a\x6a\xa4\x82\x54\x8a\x47\x14\x0c\x45\x8a\x16\xe5\x69\x83\xc2\xe3\x71\x29\x1f\x98\x58\x03\x35\x70\x2f\x13\x0d\x4f\xcc\x6c\xe0\x6e\x4e\xa6\x64\x18\xc3\x3f\xee\x80\x99\x16\xe7\x40\xf3\xdc\xea\x59\x31\xa5\x0d\xd0\x95\x41\x05\x93\xd1\x5b\x88\x37\x08\xd2\x6c\x50\x81\x79\x92\x1e\x6f\xc5\x90\x67\x1a\xa8\x42\xd1\x32\x40\x35\xb0\x6d\x2e\x95\xa1\xc2\x40\x52\x18\xbb\xee\xed\xa6\xd9\x3d\x4d\x51\x18\x30\xb2\xba\x93\xb5\x98\x72\xb6\x16\x5b\x14\xc6\xc3\x19\x09\x6b\x34\x40\x05\x7c\xea\x59\x5f\x40\xc2\x65\xfa\xf0\x36\x00\xf0\x81\x3f\x4e\x84\x70\x16\x43\xb8\x98\x4e\x61\x44\xae\x07\x8b\x69\x5c\x4f\x07\xf8\xfc\xf9\x48\xdc\xc6\x95\x1a\x83\xdb\xdc\x80\xde\x52\xce\x99\x30\xcf\x21\xde\x59\xb1\x2d\xfd\xb1\xdc\x8b\xea\xe7\xb2\xc7\xd1\x44\x9a\x6e\xa0\x2d\xe4\xe1\x3a\x20\x10\x33\xcc\x3a\x07\x8d\x98\x2d\xa9\x01\xc3\xb6\xa8\x0d\xdd\xe6\xe6\x2f\xab\x26\x55\x48\x9f\xef\x3c\xb7\x29\x9c\xdd\xb6\x3b\xf6\xc0\x8a\x09\xca\xd9\x5f\x27\xc1\xaa\x8c\x7d\x29\x9c\xbf\x44\x99\x7f\x6d\x2d\xb7\x08\x4f\x54\x1b\xcc\x20\xa7\x59\x66\xf3\x24\x57\x32\xa1\x09\xdf\xd9\x8b\xe4\x8a\x49\xc5\xcc\xee\x82\xef\xde\x97\xa0\x66\x97\xa3\x06\x6d\xa4\xc2\x0c\x64\x61\x7a\x72\xd5\x4b\xa8\xc8\xac\x3f\xd4\x5a\xc3\xbd\x96\x22\xe9\x36\xbc\x93\xec\xc0\xe0\x0f\xf3\xfb\x1f\x76\x19\x95\x92\x6a\x2f\xe6\x57\x1e\x98\xc8\x9c\x40\x2d\x0a\x00\x5b\x34\x34\xa3\x86\x7a\xc9\x13\xf9\xf0\x9f\xff\xda\x44\xa8\xd4\xfd\x59\x60\x81\x4d\x98\x83\x6c\x86\x2b\x5a\x70\x63\x0f\x58\x09\x2b\x6f\xe8\x5a\xc3\x23\x55\xe9\x86\xaa\xf6\x87\x9f\x7e\xea\x58\x63\x02\x80\xe1\x2c\x9c\xc7\xd1\x60\x12\xc6\xb5\x90\x48\xb5\xac\xc7\x67\x29\x0a\xce\x61\xf8\x85\x0c\x7f\x83\x76\xdb\x27\xec\x24\x84\x76\x9d\xa3\x1a\x0c\x55\xe7\xa7\x0e\x0c\xc2\x51\x33\xdc\x93\x79\x65\x73\x07\x66\xd1\xf3\x4d\xbb\xd1\x6d\x1a\x57\xcf\xe3\x25\xd3\xcb\x5c\x6a\x66\xd8\x23\x96\x66\x35\xf2\xfc\x17\x78\x77\x7c\xbe\x8c\xf9\x92\x89\xa5\xa2\x62\x5d\x1d\xac\x92\xe1\x97\x2b\x78\xef\x6c\xad\x56\x7e\xbe\x82\x8f\xc7\x38\xce\xef\x4b\x8e\x62\x6d\x36\x25\x84\xf5\xe9\x7e\xa9\xed\xf6\x3b\xd6\x02\x87\x75\x62\xeb\x67\x78\xff\xe1\xd3\x31\xac\x4d\x8a\x0b\xa8\x76\xfb\x34\xa8\xdf\xf1\x98\x8e\xdf\x7b\x3d\xb8\x45\xd8\xd2\x1d\x3c\x51\x4f\x53\xa9\x14\x9a\x65\xa8\x60\xff\x29\x50\xe1\xb9\x2f\x57\x96\x66\xcd\x0e\x36\xa8\x70\xcf\x8e\x77\x16\xef\x0e\xd8\x0a\x98\x01\x8d\xb8\xd5\x16\x90\xb3\x07\xb4\x8c\x9a\x41\x82\x50\x68\x5c\x15\xdc\x91\x9e\xfd\xd0\xcc\x86\x89\xf5\xdb\xb2\xac\x4c\xc2\x11\xf9\x56\x63\x37\x97\xec\xb3\xf0\xb0\x02\x8b\xf9\x24\x1c\x43\x62\x14\xa2\x37\xfe\x50\x93\x8e\x0f\xbb\x4c\x5b\x52\x91\x35\xb3\x91\x89\x0c\x7f\x9c\x07\xf5\x0c\xd9\x48\xaa\x0e\xdc\x7e\x21\x11\x39\x9b\x85\xe7\x4d\xd8\xe7\x82\x3b\xb3\x42\x93\xda\xcb\xbe\xd0\x00\x17\xee\x6e\x95\x4d\xdd\x06\xad\x75\x81\x5d\xba\xb9\xa5\x96\xf3\x6a\xc6\x93\xb0\x6d\x25\x2e\x00\x94\x74\x72\x19\xa4\x94\xaa\x01\xcd\x22\x88\xc8\xcd\x74\x30\x24\x70\xbd\x08\x87\xf1\xa4\x7e\x74\x29\xa4\x61\xab\x5d\xdb\x92\x68\x44\xe2\x45\x14\xce\x21\x8e\x26\xe3\x31\x89\x02\xb0\xcd\xcc\x9b\x37\xc1\x88\x0c\xa7\x83\x88\x58\x9a\xa5\x3b\x2e\x69\xe6\x18\xad\x1f\xfc\x4a\xc6\x93\x30\x00\x98\x5c\x43\x48\x6e\xdf\x7a\x16\xb9\x6a\xd4\xb9\xf8\x0b\xb1\x12\x8e\x77\x43\xa7\x09\x9e\x18\xe7\x90\x4a\xca\x51\xa7\x08\x59\x91\x73\x96\xda\x83\xce\x10\x96\x32\xd7\x6b\xf9\xea\xcf\x04\x50\x30\x8a\x0a\x4d\x53\xbb\xda\x05\x2d\x4b\xb4\x07\xc4\xdc\xb6\x18\x1a\x4b\xab\x34\xac\x51\xa0\xf2\xe9\xf0\xd9\xc9\x95\xf6\x5e\x39\x8b\x97\x49\xc1\x78\xb6\x94\xc9\x3d\xa6\xa6\xdd\x72\xd1\x6c\x75\x9d\xed\xfe\x43\xee\xbb\x43\x37\x24\xba\x9e\x45\x5f\xdd\x33\x40\xbe\x2e\x5d\xd4\xf2\x4e\x63\x42\xa3\x32\xad\x6e\x09\xee\x09\xd9\x9d\x25\xe1\x08\x26\xd7\xfd\xca\x95\xfb\x4c\x24\xe1\xa8\x1f\xbc\x79\x13\x4c\x07\xe1\x78\x31\x18\x13\xc8\x79\xbe\xd6\x7f\xf2\x5a\xeb\xe6\x3d\xbe\x0f\x8b\xd7\x67\xdd\x7f\x1d\x93\x08\x26\xe1\x9c\x44\x71\x23\xe0\x01\xc0\xf5\x2c\x02\x32\x18\x7e\x81\x68\x76\x6b\x55\x7f\x23\xc3\x45\x4c\xe0\x26\x9a\x0d\xc9\x68\x11\x91\x13\x21\x3e\xe8\x5b\x84\xd3\xd9\x78\x4c\x46\x8d\x9e\x91\x23\xcd\x50\xb5\x5f\xd1\x3e\x20\xc7\xf4\x42\x8b\xe0\x8a\xe5\x8f\x9c\x29\xd4\xe7\x45\xfe\x6f\x4d\xf6\x66\x2d\xd9\x89\x0a\x2b\xe8\x76\x5f\x30\x8f\xdb\xd9\x1a\x0f\x5b\xa1\x0b\x3c\x6c\xb7\x4f\xf3\xb0\xdf\x39\xc9\xed\x95\x4d\x17\x80\x2b\x99\xd3\xe8\xb5\xed\x26\xd5\x3f\x9f\x3a\xfe\xe5\xa7\x8e\xc1\xd4\x66\xc3\x51\x97\x0f\x7e\x75\x38\x9b\x2e\xbe\x86\xbe\x25\x98\x93\xb8\xd1\x63\xf4\x83\xc5\xcd\xc8\x46\xfd\x70\xc8\x8a\x38\xd9\x2b\xdf\x85\x78\x2e\x75\x2b\xfb\x62\xdd\x7f\x95\xba\x1a\xeb\x9e\xbc\xc1\x47\x7f\x83\x5e\xcf\xf5\xe9\xae\xcb\x4a\x25\x2f\xb6\x02\x04\x5a\xd9\x0d\xcd\x80\x1e\xba\x1e\x5b\xde\x8c\xa2\xb6\x8d\x93\x0a\xf6\xad\x0f\x3c\x52\x5e\xa0\x9d\x0e\xcc\x06\x2d\x94\xa5\xb9\x84\x6a\x04\x8e\x8f\xc8\xbb\x60\x36\xb2\x58\x6f\xe0\x09\xc1\x28\x86\x99\xad\x93\x28\x74\xa1\x10\xa4\x28\xcf\xd9\xf1\xc1\xf1\x8d\xb5\xcb\x9d\x7b\xfb\x92\x8b\x3a\x8b\x5f\xe8\x57\x27\xdb\xf0\xab\x5b\x79\x8d\x5f\x2b\x75\x07\xbf\xbe\xf4\xd4\x28\x9a\xdd\x94\x56\xfa\x68\x58\x8f\x57\xad\xe8\xdf\xca\xeb\x95\xd5\x2f\xf4\x7c\x25\xdf\xf0\x7e\xb5\xfa\x9a\x08\x34\x54\x37\xb3\xdb\xfa\xb3\x95\xa3\xb0\x8d\x57\xcb\x4e\xaa\xfb\xf9\xcf\xd5\x35\xdf\x4c\x65\xae\x95\x2a\xf2\x54\x6e\x6d\x77\xb6\x2a\x84\x2b\x5f\x94\x33\xb3\xfb\x5c\xaa\x3f\xf9\xf3\x60\x34\x82\x7f\x0f\xa6\x0b\xe2\x6a\xea\x2c\x06\xf2\x6d\x32\x8f\xe7\x35\x7d\xbe\x24\xd4\x9a\xf1\x7e\x70\xe6\x3e\x2e\x13\x5e\x3c\x0b\x9c\xf5\xca\x68\xf4\xfa\x81\xc2\x95\xce\xf6\xa9\x19\xc0\xd1\xe0\x61\x3e\x7e\xf9\xc8\x61\xa7\x8a\x33\xb0\x65\xfe\x1e\xa0\x5f\x01\xeb\x58\xd7\xb9\xea\x54\x0d\x6e\x14\xdd\xbe\x97\x3b\xdb\x43\x9d\xfe\x13\xe3\x3a\x8b\xf6\xb9\x62\xd5\x28\x67\xaf\x9e\xb8\x5f\x3b\x61\xe6\xd4\x26\xe7\x89\xf1\xbc\xc8\xb3\x4b\x9a\x9d\x93\x9e\xe7\x87\x2f\x61\x2e\x2a\x8d\x8f\xc7\xdd\xb4\xf1\xcd\x96\xd3\x6b\xd7\x09\x1f\x67\x66\xad\x42\x7b\x81\xa3\xa4\x3b\x51\xc1\x9d\x8a\xab\x03\x70\xa7\x1f\xfc\x2f\x00\x00\xff\xff\xd8\x85\xad\x9c\xb1\x13\x00\x00")

func _000002_riverUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_riverUpSql,
		"000002_river.up.sql",
	)
}

func _000002_riverUpSql() (*asset, error) {
	bytes, err := _000002_riverUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_river.up.sql", size: 5041, mode: os.FileMode(420), modTime: time.Unix(1721491002, 0)}
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
	"000001_new_cache_table.down.sql": _000001_new_cache_tableDownSql,
	"000001_new_cache_table.up.sql":   _000001_new_cache_tableUpSql,
	"000002_river.down.sql":           _000002_riverDownSql,
	"000002_river.up.sql":             _000002_riverUpSql,
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
	"000001_new_cache_table.down.sql": &bintree{_000001_new_cache_tableDownSql, map[string]*bintree{}},
	"000001_new_cache_table.up.sql":   &bintree{_000001_new_cache_tableUpSql, map[string]*bintree{}},
	"000002_river.down.sql":           &bintree{_000002_riverDownSql, map[string]*bintree{}},
	"000002_river.up.sql":             &bintree{_000002_riverUpSql, map[string]*bintree{}},
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
