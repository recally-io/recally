// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// database/migrations/000001_create_cache_table.down.sql
// database/migrations/000001_create_cache_table.up.sql
// database/migrations/000002_create_users_database.down.sql
// database/migrations/000002_create_users_database.up.sql
// database/migrations/000003_create_assistant_threads.down.sql
// database/migrations/000003_create_assistant_threads.up.sql
// database/migrations/000004_assistant_embeddings_add_uuid.down.sql
// database/migrations/000004_assistant_embeddings_add_uuid.up.sql
// database/migrations/000005_create_new_bookmarks_tables.down.sql
// database/migrations/000005_create_new_bookmarks_tables.up.sql
// database/migrations/000006_set_embeddings_to_nullable.down.sql
// database/migrations/000006_set_embeddings_to_nullable.up.sql
// database/migrations/000007_new_auth_flows.down.sql
// database/migrations/000007_new_auth_flows.up.sql
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

var __000003_create_assistant_threadsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x97\x4b\x53\xdb\x30\x10\xc7\xef\xfe\x14\x7b\x23\x9e\x81\x03\x30\xe9\x81\x9c\x4c\x10\xe0\x36\x28\xa9\x2c\xd3\xd0\x8b\x47\x44\x3b\xc1\xd3\xf8\x31\x96\xd2\x09\xed\xf4\xbb\x77\xe2\xb7\x4d\x5e\xa5\x4d\x87\x4e\x7a\xf4\x6a\x1f\xd6\xee\xff\x27\xd9\x27\x27\xd0\x4f\x50\x68\x84\xd1\xcd\x3d\x4e\x74\x94\x00\x59\x68\x0c\x95\x1f\x85\x46\x9f\x11\x8b\x13\x20\x63\x4e\xa8\x63\x0f\x29\xd8\xd7\x40\x87\x1c\xc8\xd8\x76\xb8\x03\x5f\x53\xff\x9e\x61\x54\x49\xb8\x78\x9c\xa1\x02\x4b\x29\x5f\x69\x11\xea\x22\x05\xb7\x2e\x07\xa4\x15\x2e\x0a\x27\x05\x1d\x03\x00\xc0\x97\xe0\x10\x66\x5b\x03\x18\x31\xfb\xce\x62\x0f\xf0\x81\x3c\x1c\xa7\x4b\xf3\xb9\x2f\xc1\x75\xed\x2b\x70\xa9\xfd\xd1\x25\x69\x1e\xea\x0e\x06\x70\x45\xae\x2d\x77\xc0\x61\x8a\xa1\x97\x88\x50\x46\x81\xb7\x74\xee\x98\x79\xa0\xc2\xc4\x2b\x62\x19\xb9\x26\x8c\xd0\x3e\x71\x52\xbb\xea\x2c\x3d\x73\xc7\x50\x04\x08\xf7\x16\xeb\xdf\x5a\xac\x73\xd6\xed\x9a\x65\x89\x6c\x5d\xa2\x9a\x24\x7e\xac\xfd\x28\x04\x4e\xc6\x3c\xb3\xaa\x67\xa5\x31\xf0\xe2\x24\x0a\x62\x5d\xb3\x07\x91\xc4\x59\x99\xee\xfc\xac\x9d\x2d\x40\x2d\xa4\xd0\x02\xde\x3b\x43\x7a\x59\x6e\xe2\xe8\xfb\x8f\xa3\x8b\x8b\xd4\x96\xf9\x4d\xd2\xb6\x4a\x4f\x68\xe0\xf6\x1d\x71\xb8\x75\x37\x82\x4f\x36\xbf\x4d\x1f\xe1\xf3\x90\xae\x68\x45\xdf\x65\x8c\x50\xee\x95\x11\x79\x2b\x62\xf9\x07\x72\x19\x66\xaf\x18\xaa\x4d\xaf\xc8\xb8\x35\x54\x5f\x2e\xbc\xb4\xe7\xb5\x37\x1f\xd2\xc6\xac\xf3\x91\x1c\xd7\x36\x67\xf6\x0c\x63\x27\xa9\x78\xfa\x29\x41\x21\xdf\x84\x62\xaa\x97\x5a\xe1\x5d\x6d\xf8\xbf\xc8\xf6\x27\xb2\x6a\x04\x6b\xe4\x56\xe9\xa5\x54\x5d\x7d\x6c\x2d\x0d\x6e\x29\xf9\x0b\xd5\x36\x14\xd9\x51\xe8\x01\x2a\x25\xa6\xf8\xaf\x2a\x3d\xeb\xc4\x46\xff\xa2\x5b\xf5\xb0\x17\x92\xce\xcc\x49\x34\xdb\xc8\x8d\xc6\x45\x1d\x8c\x8c\x14\x4f\x47\x5f\x30\x04\x9b\x72\x72\x43\x58\x2e\xf5\x28\x88\x67\xb8\xe4\x6b\xd5\x2a\x06\x8f\x28\xa5\x1f\x4e\x55\x7e\xb7\x41\xe7\xb4\x7b\xfe\xce\x3c\x14\x9c\xd6\xc9\xba\xd2\xe2\xea\xb3\xfb\xb7\xb9\xa9\x0a\xbc\x9e\xce\x5c\x71\xdb\x4b\x94\xd2\x7c\x1d\x98\x42\x6b\x31\x79\x0a\xf0\x8d\x7c\xb7\xfc\x3d\x36\x5f\x5c\x5e\x79\xb6\xe7\x78\x95\x79\x9e\x54\x1c\x77\x4f\x0b\x90\x95\xff\x0d\x9b\xcc\x1d\x2e\x55\x0d\x21\xed\x09\xac\x46\x8d\xbd\xb1\xd5\xa8\xb2\x1e\xaf\xda\x6f\xc2\xf2\xbc\x26\xd5\x69\x9b\xfe\x36\xec\x86\x5f\x76\x46\x67\x61\xdb\xf1\xdb\x99\xa2\x72\x07\x85\x7b\xeb\x62\x69\x5d\x38\xeb\x6e\x8a\x43\xfb\x02\x6b\xea\xa0\x31\x9b\xbc\xf7\x3b\xa8\xb8\xec\xfd\x86\x6c\x8d\x01\x99\x3d\xe3\x67\x00\x00\x00\xff\xff\x7a\x3e\xcf\x42\xb8\x0e\x00\x00")

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

	info := bindataFileInfo{name: "000003_create_assistant_threads.up.sql", size: 3768, mode: os.FileMode(420), modTime: time.Unix(1731662456, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000004_assistant_embeddings_add_uuidDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\x2c\x2e\xce\x2c\x2e\x49\xcc\x2b\x89\x4f\xcd\x4d\x4a\x4d\x49\x49\xc9\xcc\x4b\x2f\x56\x70\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x53\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x2d\xcd\x4c\xb1\xe6\x02\x04\x00\x00\xff\xff\x4d\xb6\xd5\x15\x3e\x00\x00\x00")

func _000004_assistant_embeddings_add_uuidDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_assistant_embeddings_add_uuidDownSql,
		"000004_assistant_embeddings_add_uuid.down.sql",
	)
}

func _000004_assistant_embeddings_add_uuidDownSql() (*asset, error) {
	bytes, err := _000004_assistant_embeddings_add_uuidDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_assistant_embeddings_add_uuid.down.sql", size: 62, mode: os.FileMode(420), modTime: time.Unix(1725438658, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000004_assistant_embeddings_add_uuidUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x1c\xca\xb1\x0a\xc2\x30\x10\x06\xe0\xdd\xa7\xf8\x47\x7d\x06\xa7\x68\xae\x10\x38\x53\xb4\x77\xe0\x16\x22\x17\x4a\x86\x46\x30\xcd\xfb\x8b\xce\xdf\xe7\x58\xe8\x01\x71\x17\x26\xe4\xde\x6b\xdf\x73\xdb\x53\xd9\x5e\xc5\xcc\x6a\x5b\x3b\x9c\xf7\xb8\xce\xac\xb7\x88\x30\x21\xce\x02\x7a\x86\x45\x16\x8c\x51\x0d\xaa\xc1\x43\x63\xb8\x2b\xfd\x2d\x2a\x33\x3c\x4d\x4e\x59\xb0\x96\x96\x3e\xb9\xd9\x7b\x4b\xbf\x7c\x3c\x9d\x0f\xdf\x00\x00\x00\xff\xff\xd5\xa2\x18\xce\x70\x00\x00\x00")

func _000004_assistant_embeddings_add_uuidUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_assistant_embeddings_add_uuidUpSql,
		"000004_assistant_embeddings_add_uuid.up.sql",
	)
}

func _000004_assistant_embeddings_add_uuidUpSql() (*asset, error) {
	bytes, err := _000004_assistant_embeddings_add_uuidUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_assistant_embeddings_add_uuid.up.sql", size: 112, mode: os.FileMode(420), modTime: time.Unix(1725438693, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000005_create_new_bookmarks_tablesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4c\xa9\x88\x4f\xca\xcf\xcf\xce\x4d\x2c\xca\x2e\x8e\x2f\x2d\x4e\x2d\x8a\xcf\x4c\xb1\xe6\x22\x46\x6d\x51\x0e\x51\xea\x92\x8b\x52\x13\x4b\x52\x53\xe2\x13\x4b\xac\xb9\x20\xea\x43\x1c\x9d\x7c\x5c\x91\xd4\xc3\xd5\x5a\x03\x02\x00\x00\xff\xff\x1c\xa0\xdd\x93\xa3\x00\x00\x00")

func _000005_create_new_bookmarks_tablesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000005_create_new_bookmarks_tablesDownSql,
		"000005_create_new_bookmarks_tables.down.sql",
	)
}

func _000005_create_new_bookmarks_tablesDownSql() (*asset, error) {
	bytes, err := _000005_create_new_bookmarks_tablesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000005_create_new_bookmarks_tables.down.sql", size: 163, mode: os.FileMode(420), modTime: time.Unix(1731990367, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000005_create_new_bookmarks_tablesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\x41\x6f\xe2\x30\x10\x85\xef\xf9\x15\x73\x4c\x24\x2e\xbb\x2b\xf6\xc2\xc9\x84\x41\x78\x37\x38\xd4\x71\x5a\xe8\x25\x32\xd8\x82\x88\x38\xa9\x1c\x53\xb5\xff\xbe\x4a\x6a\x50\x40\x55\xb9\xf4\x98\x97\xf7\x3d\x8f\x66\x5e\xcc\x91\x08\x04\x41\xa6\x09\x02\x9d\x03\x4b\x05\xe0\x9a\x66\x22\x83\x6d\xd3\x1c\x8d\xb4\xc7\x16\xc2\x00\x00\xa0\x54\x90\x21\xa7\x24\x81\x15\xa7\x4b\xc2\x37\xf0\x1f\x37\xa3\xfe\xd7\xe9\x54\x2a\xc8\x73\x3a\x83\x9c\xd1\x87\x1c\xfb\x18\x96\x27\x09\xcc\x70\x4e\xf2\x44\xc0\x5e\xd7\x85\x95\xb5\x6a\x4c\xd1\x99\xc3\xc8\x83\xad\xb6\xc5\x99\xe5\x38\x47\x8e\x2c\xc6\xac\xd7\xdb\xb0\x73\x9e\x8d\xb6\x02\x81\x6b\x71\x49\xfe\x94\x5d\xe9\x2a\x0d\x8f\x84\xc7\x0b\xc2\xc3\xdf\xe3\xb1\xb7\xb7\x27\x63\xa4\x7d\xef\x91\x2b\xa5\xd0\x66\xab\x95\x2a\xeb\x7d\x0b\xaf\x7a\xe7\x1a\x1b\xfe\x1a\xff\xf9\xeb\xb1\x5d\x53\x3b\x5d\xbb\x01\xe6\x95\xef\xb1\x83\x33\xd5\x80\x31\xda\x49\x25\x9d\x84\x7f\x59\xca\xa6\xfe\xf9\x9d\xd5\xba\x6e\x0f\xcd\x55\xb8\xd5\xd2\x69\x55\x48\x07\x82\x2e\x31\x13\x64\xb9\x82\x27\x2a\x16\xfd\x27\x3c\xa7\xec\x8b\x4d\xc6\x39\xe7\xc8\x44\x71\x21\xfc\x82\x5e\xd4\x0f\x64\x05\xd1\x24\x08\x7c\x27\x28\x9b\xe1\xfa\xa6\x13\xa5\x7a\x2b\xfa\x9b\x0d\x46\x4f\xd9\xb0\x2a\xfe\xa2\xa3\xc1\x40\xd1\xe4\x6e\x64\xd7\x9f\x9b\x9c\xee\xf6\xf7\x41\x5b\xdd\x72\xb6\x8a\x26\xc1\x47\x00\x00\x00\xff\xff\x1b\x3c\x35\xe8\xd7\x02\x00\x00")

func _000005_create_new_bookmarks_tablesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000005_create_new_bookmarks_tablesUpSql,
		"000005_create_new_bookmarks_tables.up.sql",
	)
}

func _000005_create_new_bookmarks_tablesUpSql() (*asset, error) {
	bytes, err := _000005_create_new_bookmarks_tablesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000005_create_new_bookmarks_tables.up.sql", size: 727, mode: os.FileMode(420), modTime: time.Unix(1731990367, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000006_set_embeddings_to_nullableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\x2c\x2e\xce\x2c\x2e\x49\xcc\x2b\x89\x4f\xcd\x4d\x4a\x4d\x49\x49\xc9\xcc\x4b\x2f\x56\x80\xa8\x70\xf6\xf7\x09\xf5\xf5\x53\x80\x48\x80\xc5\x43\x22\x03\x5c\x15\xca\x52\x93\x4b\xf2\x8b\x34\x0c\x4d\x8d\xcd\x34\x75\x70\x2a\x0d\x76\x0d\x51\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\xb1\xe6\x02\x04\x00\x00\xff\xff\x54\xfe\xd3\xbb\x73\x00\x00\x00")

func _000006_set_embeddings_to_nullableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000006_set_embeddings_to_nullableDownSql,
		"000006_set_embeddings_to_nullable.down.sql",
	)
}

func _000006_set_embeddings_to_nullableDownSql() (*asset, error) {
	bytes, err := _000006_set_embeddings_to_nullableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000006_set_embeddings_to_nullable.down.sql", size: 115, mode: os.FileMode(420), modTime: time.Unix(1731990367, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000006_set_embeddings_to_nullableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\x2c\x2e\xce\x2c\x2e\x49\xcc\x2b\x89\x4f\xcd\x4d\x4a\x4d\x49\x49\xc9\xcc\x4b\x2f\x56\x80\xa8\x70\xf6\xf7\x09\xf5\xf5\x53\x80\x48\x80\xc5\x43\x22\x03\x5c\x15\xca\x52\x93\x4b\xf2\x8b\x34\x0c\x4d\x8d\xcd\x34\x75\x70\x2a\x75\x09\xf2\x0f\x50\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\xb1\xe6\x02\x04\x00\x00\xff\xff\x46\x02\x96\x56\x74\x00\x00\x00")

func _000006_set_embeddings_to_nullableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000006_set_embeddings_to_nullableUpSql,
		"000006_set_embeddings_to_nullable.up.sql",
	)
}

func _000006_set_embeddings_to_nullableUpSql() (*asset, error) {
	bytes, err := _000006_set_embeddings_to_nullableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000006_set_embeddings_to_nullable.up.sql", size: 116, mode: os.FileMode(420), modTime: time.Unix(1731990367, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000007_new_auth_flowsDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\x4f\x6f\x9b\x40\x10\xc5\xef\xf9\x14\x73\xb4\x0f\xb9\x54\xca\xc9\x27\x6a\x48\x8a\xe4\x42\x85\x49\x95\x1b\xda\xec\x4e\x60\x65\x98\x41\xbb\x43\x9a\x7e\xfb\x8a\x3f\x76\x71\xea\x9a\x36\x17\x04\xab\x37\xbf\x79\xf3\x66\xb9\xbd\x85\xd0\x71\x0b\xa2\x9e\x6b\xf4\xa0\xc8\x80\x54\x68\x1d\x28\xef\x59\x5b\x25\x68\x40\x9c\x2d\x4b\x74\xfe\x26\xcc\xd2\x6f\x90\x07\x9f\x77\x11\xc4\xf7\x10\x3d\xc5\xfb\x7c\x0f\xaa\x93\xaa\x70\xf8\xca\x07\x34\x85\xf0\x01\xc9\x6f\xae\x28\x55\x6b\x8b\x03\xfe\xbc\xaa\xe9\x3c\xba\x82\x87\x57\xcd\x44\xa8\xc5\x72\x4f\xbd\x39\xb9\x9d\x1c\x0d\x7e\x5f\x3a\x1a\x05\x13\x31\x8b\x1f\x1e\xa2\x6c\xc6\xec\x5a\xa3\x04\x07\xaa\x2f\xc6\x0f\x53\x28\x81\x34\x81\xe1\x6c\xb3\x50\xf8\x87\x93\x77\x90\xab\x9e\xaf\xa3\x8f\x69\x5c\x22\xbe\x4b\xea\xfe\x31\xd9\xe6\x71\x9a\x5c\x18\xec\x54\x5b\x68\xae\xbb\x86\x56\xeb\x59\x56\x96\x0c\xbe\xa1\x07\xa6\x71\xda\x71\xd5\x23\x33\x4e\xc2\xe8\x69\x06\xb4\xe6\x6d\x8a\x09\x1b\x65\xeb\xcd\x92\xaa\xad\x98\x70\x51\xd5\x3f\x49\x35\xcb\x42\x2f\x4a\xba\x69\xcf\x19\x36\xfc\x8a\xa0\x99\xbc\x38\x65\x49\x3c\xbc\x38\x6e\xce\x66\x08\x76\x79\x94\x4d\x57\x68\x3c\x1f\x1a\x6c\xd3\x64\x9f\x67\x41\x9c\xe4\xf3\xac\x7e\x8f\x55\xe8\x0a\xf5\x61\xf3\x81\x72\xcd\x24\x4a\xcb\x11\x70\xe6\xb3\x0f\xfe\x7f\x3c\xee\x1e\xbf\xce\x77\x39\x45\xf9\xaf\x72\x8f\x22\x96\xca\x53\x58\x5e\xd8\x21\xb0\xb3\xa5\x25\x55\x9f\xec\xac\xa4\x42\x8f\xf0\x03\x1d\x82\x71\xdc\xb6\x68\xc0\x52\xff\x8b\x43\xd7\x42\x63\x4b\xa7\xfa\x7b\xba\xbe\xd0\x37\x08\xc3\x59\xdb\x24\xcd\x8f\xad\x4b\x2b\x55\xf7\x0c\xdf\x83\x6c\xfb\x25\xc8\x56\x9f\xee\xee\xd6\x97\x7c\xff\xbd\x9e\xb9\xac\xf1\xe3\xf5\x82\x35\x96\x4e\x35\xe7\x84\x5f\x01\x00\x00\xff\xff\xcd\x18\xc0\x75\xc9\x04\x00\x00")

func _000007_new_auth_flowsDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000007_new_auth_flowsDownSql,
		"000007_new_auth_flows.down.sql",
	)
}

func _000007_new_auth_flowsDownSql() (*asset, error) {
	bytes, err := _000007_new_auth_flowsDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000007_new_auth_flows.down.sql", size: 1225, mode: os.FileMode(420), modTime: time.Unix(1733879939, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000007_new_auth_flowsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\x6d\x8f\xdb\xb8\x11\xfe\xee\x5f\x31\x1f\x5c\xd8\x6e\xac\x45\x7a\x40\x80\x76\x17\x77\xa8\xd6\xe6\x66\x75\xf1\xca\x5b\x59\xee\xe6\x72\x4d\x05\x46\xa2\x6d\xc6\x32\xa9\x90\xd4\x66\x7d\x87\xeb\x6f\x2f\x48\xea\xc5\xb2\xfc\xb2\x1b\x04\x68\xfd\xc9\xa6\x38\xcf\x70\x86\x33\xcf\x33\xb2\xe3\x40\xb8\xa2\x12\x36\x74\x29\xb0\xa2\x9c\x81\x24\x4a\x42\x9e\x81\x5a\x11\xc0\xb9\x5a\x11\xa6\x68\x6c\x1f\x61\x96\x98\x25\x2e\xe8\x6f\xc5\xe6\x78\x45\x36\xb8\xe3\x38\xe0\x29\xa0\x2c\x4e\xf3\x84\x48\xc8\x25\x11\xb0\xc1\x0c\x2f\xc9\x86\x30\x35\x84\xa9\x9b\xab\x15\x50\xa6\x48\xe1\x65\x08\x8f\x44\xd0\x45\x09\x2c\xb7\x52\x91\x8d\x1c\x1a\x0f\xc1\xb5\x3b\xea\x74\xdc\x49\x88\x02\x08\xdd\xeb\x09\x32\x78\x12\xdc\xf1\x18\x46\xd3\xc9\xfc\xce\x07\xef\x06\xfc\x69\x08\xe8\xbd\x37\x0b\x67\x90\xad\x38\x23\xf0\x4f\x37\x18\xdd\xba\x41\xff\xcd\xeb\xc1\xd5\x01\xeb\x71\x30\xbd\xdf\x31\x2f\x4c\x97\x54\xad\xf2\x4f\x2f\xd8\xcf\xf9\x32\x25\xcf\xdf\xaf\x48\xaa\x43\xde\x1c\xb2\x38\x1a\x8f\x24\x4a\x51\xb6\x94\xf0\xf3\x6c\xea\x5f\xc3\x18\xdd\xb8\xf3\x49\x08\xbd\xdf\xff\xe8\x5d\x5e\x9a\xb5\xab\xce\x78\x0a\xdd\x2e\x5c\xa3\xb7\x9e\xdf\x01\x00\x38\x06\xef\xcf\xc2\xc0\xf5\xfc\xd0\x2e\x46\x64\x83\x69\x1a\xc5\x2b\x12\xaf\xc1\xd8\xe9\xcf\xe8\x16\x8d\xde\x41\xdf\x3c\x83\xff\xfc\x19\x7a\xff\xfe\xd5\x75\x3e\x60\xe7\xb7\xd7\xce\xdf\x2e\xa2\x3f\xbd\x72\x3e\xbe\xfa\xfb\xce\x8a\xf3\xf1\xd5\xbf\x2e\x8a\xdf\x1f\x7f\xff\x61\xf8\x47\xb7\x37\xb8\x32\x60\xe8\xfd\x08\xdd\x87\xde\xd4\x87\x87\x5b\xe4\x43\x92\x67\xa9\xbe\x61\x12\xf1\x4f\x9f\x49\xac\x20\xd4\xab\xfe\x7c\x32\xb9\xea\x20\x7f\x0c\xdd\xee\x37\xc7\x11\x73\xa6\x70\xac\x4e\x46\xe2\xcd\x4c\x52\xb5\x3f\x98\x06\x45\x95\xec\x2d\x6a\x34\x86\x37\x8d\xf5\x6f\x0f\x46\xf7\xc1\x28\x40\x6e\x88\x8a\x08\x9a\xf7\x6a\xe3\xe9\xeb\x5d\xfa\x43\x13\x98\xa1\xc0\x73\x27\x70\x1f\x78\x77\x6e\xf0\x0b\xbc\x43\xbf\x0c\xcb\xa7\x79\x4e\x13\x98\xcf\xbd\x31\xcc\x7d\xef\x1f\x73\x54\x1f\xbb\xac\x87\x25\x61\x91\xc0\x2c\xe1\x9b\x48\x6f\xee\x0f\x6a\xdb\x32\xaa\xb2\x27\x7e\x78\xf3\xa6\x7e\x6a\xb3\xb3\xfb\xa8\x70\x51\xed\x68\x35\x54\x6b\x03\x96\xf2\x2b\x17\x49\xb4\xc2\x72\x05\x21\x7a\x1f\x56\x8f\x70\xac\xe8\xa3\x4e\x13\x96\x92\x4a\x85\x99\x8a\x8a\x40\xda\x5b\xd4\x4a\x10\x9c\xb4\x9e\x4b\x85\x55\x2e\x9b\x27\x6c\x45\xdf\xcb\x08\x4b\x28\x5b\xf6\x6a\xb3\xb3\x6d\x53\x6d\x8d\x05\xc1\x8a\x24\x11\x56\x10\x7a\x77\x68\x16\xba\x77\xf7\xf0\xe0\x85\xb7\xe6\x27\x7c\x98\xfa\x07\xf2\x3d\x9a\x07\x01\xf2\xc3\xa8\xb2\xa8\xf3\x9d\x25\xdf\x01\x4e\xa3\x0d\x6c\x11\xb9\x49\x02\x19\x16\x8a\xe2\x14\x28\x4b\xc8\x13\x91\xb0\xe0\x02\x32\x22\x16\x5c\x6c\x30\x8b\x49\xa7\x28\xb4\xa2\x3a\x3c\x7f\x8c\xde\xef\xd5\x1b\x4d\x9e\xa2\x9d\xb6\x87\xa9\x5f\x96\xe0\x64\xfa\x80\x02\xdb\x26\x83\x81\xae\xef\x00\x41\xab\x69\xae\x5e\xe4\xc2\xd6\x4c\xed\xc2\xfc\x2e\xb1\x5b\xbd\xf7\x32\xec\xaa\x9e\x6b\xf8\x72\xa9\xf4\x70\xa8\x91\x2b\x27\xa7\xd1\x8b\x72\xab\xb1\xa9\x8c\x4c\x89\x56\xe0\xd5\x02\xfc\x08\x4a\xe4\xa4\xbe\x24\x25\xe8\x72\x49\x84\xb9\x9c\xba\x0a\x4a\xbf\xd3\x00\x02\x74\x3f\x71\x47\x08\x6e\xe6\xfe\xc8\x90\x89\xdd\x15\xd5\x9b\xa3\x98\xa7\xf9\x86\xf5\x07\x9d\x00\x85\xf3\xc0\x9f\x41\x18\x78\x6f\xdf\xa2\x00\xdc\x19\x74\xbb\x9d\x9a\x1b\x7d\xf4\x70\xb1\x53\x6a\x3f\xb6\x4b\xc8\x32\x97\xc5\xd1\xdb\x0d\x35\x5d\x75\xba\x5d\x48\x31\x5b\xe6\x78\x49\xa0\x97\xa5\xd9\x52\x7e\x49\x7b\x57\x9d\x8e\xd1\xac\xd2\x5b\x2d\x5a\xe5\x11\x6d\xee\x6b\x87\x65\x86\xaa\xbc\x96\xa6\x47\x0c\xcc\x61\xae\xd1\xcd\x34\x40\x30\xbf\x1f\x9b\x8c\x14\x10\xe6\xd1\xcd\x34\x00\xe4\x8e\x6e\x21\x98\x3e\x14\x94\x8b\x46\xf3\xf0\x59\xc9\xb2\x57\x60\xe7\x8a\x11\x67\x8c\xc4\x7a\x90\x90\x97\x70\x67\xc6\x0e\x09\x6a\x45\x45\xe2\xe8\x26\xda\xee\x8f\x31\x99\xe0\x8f\x34\xd1\xa7\x70\x1c\x40\x0c\x7f\x4a\x8b\x99\x45\x82\xe2\x90\x52\xb6\x86\x4d\x9e\x2a\x9a\xa5\xa4\xf0\x50\x59\xe8\x0d\x6a\x45\xa8\x00\x1c\xc7\x3c\x67\x4a\x43\xcc\x14\x17\x44\x02\x23\x31\x91\x12\x8b\x2d\x28\xbe\x26\x4c\x9a\x61\xa6\xb4\x74\x64\x46\x62\x3d\xf0\x40\x82\x15\xee\x9c\x50\x09\x7d\x58\x93\xca\x88\x9b\xaf\x71\x1d\x1d\xf4\x3b\x85\x70\x18\x61\xd8\x91\x8d\x53\xa2\x50\x2a\x42\xc9\xb3\x35\x15\x05\xe8\x06\x05\xc8\x1f\xa1\x42\x9b\xfa\x34\x19\xe8\x2b\x1a\xa3\x09\x0a\x11\x8c\xdc\xd9\xc8\x1d\x23\x8b\x50\x06\xd2\x90\x85\x12\x69\x08\xf5\xa7\xba\x96\xca\xc2\xf4\x65\x9f\x5c\x2c\x2f\x86\xd0\xb3\x13\x54\x4f\x7f\x33\xb3\x57\x6f\xd0\x80\x8f\xca\x93\x1e\x64\xff\xa1\x81\x9f\x4b\x22\x7a\x12\xbc\x31\x50\x66\xc6\xd4\xd2\xb8\x27\x8b\x49\xb2\x09\xd9\x16\x3c\x1b\x12\x8e\xf5\x85\x45\xe6\xb6\xac\x88\xc1\xe1\x4f\x15\x92\xb5\xb0\xf7\x6b\x7a\xde\xbd\xf7\x20\xc6\x69\x6a\x0b\x5a\x90\x85\x20\x72\x75\x1e\x51\x4f\xdd\x15\x46\x61\x45\xd9\xb2\x71\x22\x83\x68\xbe\x45\xe4\x29\xa3\x82\xc8\x86\xbe\x84\x1f\x9a\xf7\x12\xe9\xba\xb2\xd2\x77\xd8\xa9\xe3\xc0\x4d\x4a\x9e\xe8\xa7\x94\x80\x54\x5c\x68\x32\x30\x9a\xd2\xaa\x50\xca\x8c\xc8\xe8\x92\xeb\x1c\x53\xcb\xf0\xc3\xb3\x04\xf2\x98\x3a\xbe\xc0\x7c\x77\xf6\xfb\xd2\xea\x89\x42\x40\xfa\x65\x14\xc3\x56\x25\x0d\x3a\x9a\x2c\xce\x28\x01\xaf\xbb\x8e\x26\xba\x05\x4e\x74\x61\xbf\x04\x3e\xab\x2f\xd6\xb4\x3a\x50\xca\xf9\x3a\xcf\xce\xa1\x9f\x88\xe4\x99\x0e\x77\x8a\x66\x7b\xce\xdb\x7e\x81\x95\x92\xd7\x2a\xbc\x86\xae\x9e\x51\x8e\x96\x9b\x3d\x15\x39\x71\xa0\x63\xda\x72\x0a\xf2\xb0\xce\x9c\x70\xf2\x9d\xd4\x47\x37\xff\x3b\xb2\x2d\x24\x47\xbf\xe9\x5e\x56\x72\x92\x09\xae\xdf\xfb\x74\x17\xc5\x35\x6d\x18\xb6\x72\xef\x3d\xa3\x1a\x79\x96\x71\xa1\x64\xad\x35\xc5\x8c\xb1\x26\x5b\xa9\x27\x3d\xfb\x1a\xfd\x95\xaa\x15\x24\x74\xb1\x20\x82\x30\x05\x32\xe6\x19\x31\xc2\xe5\x6d\xb2\xd4\x38\x95\xda\x02\x04\x57\x3b\x2f\xe9\x1a\x89\xaa\x2d\x28\x81\xe3\x35\x65\xcb\xb3\x82\x83\x33\x1a\x19\xc7\xff\x43\x89\x69\xbd\xb2\x1c\xd3\x17\x2d\x00\x4e\x42\x16\x94\x91\xc4\x5a\x69\x22\xd3\xb9\x5d\x93\xad\x81\x5a\x93\x6d\x94\x09\xb2\xa0\x4f\x15\xe0\x5f\xdb\x70\x9a\x12\xb9\x30\xe9\xa3\x89\x1e\x10\xaa\x3f\x24\x74\xd6\x79\xae\x80\x3c\x65\x5c\x6a\x66\x5e\xe4\x69\xda\x40\x37\x6f\x3e\x27\x0f\xab\xef\x98\xc4\xb9\x20\xe9\xd6\x10\x2e\x49\x20\xe6\xfa\xd2\x94\x39\x27\x68\x04\x03\x67\xef\xd4\xe8\xc5\xaf\x1f\x0f\xc6\x5c\xc0\xb9\x42\xe0\x2d\xf0\x85\xae\x8e\x0d\x95\xd2\xcc\x04\x4b\x81\x99\x22\x89\xad\x2e\x2a\xab\x43\x1e\x11\x8d\x03\xd2\x96\xe9\x90\xb1\x09\xcf\x5a\xd5\xd4\x9f\x62\xa9\x74\x13\xed\xb3\xf7\xb0\x25\x67\xba\xd0\x6c\x25\xe6\x12\x2f\xc9\xff\x9d\x70\x98\xf2\xd4\x17\x67\xea\xa5\x50\x8d\xa2\x66\x87\xa6\x88\x06\x2d\xb3\x78\x6d\x2c\x0a\x26\x2d\xff\x5c\x68\x72\x62\xf1\x47\xc2\xce\xea\x4f\x3b\x71\x3f\x4b\x7c\x1a\x0d\x58\x96\x6d\xc9\x60\xe5\x7a\xbf\x2e\xe9\xf3\x3a\xd0\x44\x34\x3c\xd2\xc2\x7b\xb6\xa8\x34\xc1\xf6\x64\xa5\x82\x6b\x4b\xc8\x37\x8a\xc7\xde\xd9\xdb\xc2\x51\x3e\x3b\x26\x15\xc7\x00\x4e\xc8\x44\xb9\xfb\x3b\x09\x83\x9d\xed\x02\xf2\xc8\x2d\x9d\x5c\xee\x92\xb5\x34\x9c\x00\x29\x5f\x6a\x82\xd1\x6c\x6d\xc7\x45\xca\x1e\x71\x4a\x13\xdb\x7b\x75\x47\x09\xf2\xc8\xd7\x24\x81\x9f\x1f\xc2\xf2\xbd\x22\x67\x8a\xa6\xc5\x6b\x08\xc3\x2a\x17\x38\xdd\xed\xdb\x5d\x6d\x79\xc4\x82\xf2\xdc\xa2\x14\xd4\x26\x08\x96\x9a\x37\x34\x65\xe2\x3c\xa1\x0a\xb2\x5c\x64\x5c\x12\x79\x56\x24\x8a\xb3\x44\xc5\x39\xac\x54\x7c\x56\xb4\xa5\x15\xed\x09\x54\xb3\x36\xa3\x5f\x72\x62\x02\x29\xd9\x96\x88\x17\x4a\xc7\xf0\x04\xb5\xed\x33\xa7\xe3\xc0\xc3\x8a\xb0\x22\xbd\x5f\x79\x9e\x26\x65\xba\xd2\x82\xe8\x48\x31\xb6\xdb\xb0\xbe\x91\x67\x6c\x42\x2b\x2d\xf8\xcb\xeb\xd7\x83\x03\x13\xb8\xe3\x40\x60\x37\xda\x99\xbf\xba\x8f\xbe\xad\x84\xa1\x2d\x0c\x2d\xda\x31\x67\x31\x11\x6c\x08\x44\xc5\x17\x83\x03\x5c\xd6\xb8\x87\xf2\x2f\x94\x9a\xcc\x3e\x2b\x7a\x90\xcb\xec\x4c\xb7\xe3\xba\xcd\x68\x3f\xed\x24\xe3\xf9\xdc\xd5\xac\x8b\x7d\x8a\x68\x3e\xdd\x25\x8a\xe7\x51\xcf\x1e\x7a\x83\xcd\xf6\xb0\x6b\x4e\xfb\x6f\x00\x00\x00\xff\xff\xef\x9c\xd0\xe6\xd8\x18\x00\x00")

func _000007_new_auth_flowsUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000007_new_auth_flowsUpSql,
		"000007_new_auth_flows.up.sql",
	)
}

func _000007_new_auth_flowsUpSql() (*asset, error) {
	bytes, err := _000007_new_auth_flowsUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000007_new_auth_flows.up.sql", size: 6360, mode: os.FileMode(420), modTime: time.Unix(1733881703, 0)}
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
	"000001_create_cache_table.down.sql":            _000001_create_cache_tableDownSql,
	"000001_create_cache_table.up.sql":              _000001_create_cache_tableUpSql,
	"000002_create_users_database.down.sql":         _000002_create_users_databaseDownSql,
	"000002_create_users_database.up.sql":           _000002_create_users_databaseUpSql,
	"000003_create_assistant_threads.down.sql":      _000003_create_assistant_threadsDownSql,
	"000003_create_assistant_threads.up.sql":        _000003_create_assistant_threadsUpSql,
	"000004_assistant_embeddings_add_uuid.down.sql": _000004_assistant_embeddings_add_uuidDownSql,
	"000004_assistant_embeddings_add_uuid.up.sql":   _000004_assistant_embeddings_add_uuidUpSql,
	"000005_create_new_bookmarks_tables.down.sql":   _000005_create_new_bookmarks_tablesDownSql,
	"000005_create_new_bookmarks_tables.up.sql":     _000005_create_new_bookmarks_tablesUpSql,
	"000006_set_embeddings_to_nullable.down.sql":    _000006_set_embeddings_to_nullableDownSql,
	"000006_set_embeddings_to_nullable.up.sql":      _000006_set_embeddings_to_nullableUpSql,
	"000007_new_auth_flows.down.sql":                _000007_new_auth_flowsDownSql,
	"000007_new_auth_flows.up.sql":                  _000007_new_auth_flowsUpSql,
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
	"000001_create_cache_table.down.sql":            &bintree{_000001_create_cache_tableDownSql, map[string]*bintree{}},
	"000001_create_cache_table.up.sql":              &bintree{_000001_create_cache_tableUpSql, map[string]*bintree{}},
	"000002_create_users_database.down.sql":         &bintree{_000002_create_users_databaseDownSql, map[string]*bintree{}},
	"000002_create_users_database.up.sql":           &bintree{_000002_create_users_databaseUpSql, map[string]*bintree{}},
	"000003_create_assistant_threads.down.sql":      &bintree{_000003_create_assistant_threadsDownSql, map[string]*bintree{}},
	"000003_create_assistant_threads.up.sql":        &bintree{_000003_create_assistant_threadsUpSql, map[string]*bintree{}},
	"000004_assistant_embeddings_add_uuid.down.sql": &bintree{_000004_assistant_embeddings_add_uuidDownSql, map[string]*bintree{}},
	"000004_assistant_embeddings_add_uuid.up.sql":   &bintree{_000004_assistant_embeddings_add_uuidUpSql, map[string]*bintree{}},
	"000005_create_new_bookmarks_tables.down.sql":   &bintree{_000005_create_new_bookmarks_tablesDownSql, map[string]*bintree{}},
	"000005_create_new_bookmarks_tables.up.sql":     &bintree{_000005_create_new_bookmarks_tablesUpSql, map[string]*bintree{}},
	"000006_set_embeddings_to_nullable.down.sql":    &bintree{_000006_set_embeddings_to_nullableDownSql, map[string]*bintree{}},
	"000006_set_embeddings_to_nullable.up.sql":      &bintree{_000006_set_embeddings_to_nullableUpSql, map[string]*bintree{}},
	"000007_new_auth_flows.down.sql":                &bintree{_000007_new_auth_flowsDownSql, map[string]*bintree{}},
	"000007_new_auth_flows.up.sql":                  &bintree{_000007_new_auth_flowsUpSql, map[string]*bintree{}},
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
