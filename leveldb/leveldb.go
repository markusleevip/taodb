package leveldb

import (
	_leveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	opt2 "github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)

type LevelDB struct {
	DB *_leveldb.DB
}

func NewDB(dbPath string) *LevelDB {

	opt := new(opt2.Options)
	opt.CompactionTableSize = 4 * MiB
	opt.IteratorSamplingRate = 2 * MiB
	opt.WriteBuffer = 32 * MiB

	dbInstance, err := _leveldb.OpenFile(dbPath, opt)
	if err != nil {
		panic(err)
	}
	return &LevelDB{dbInstance}
}

func (db *LevelDB) Set(key string, value []byte) error {
	return db.DB.Put([]byte(key), []byte(value), nil)
}

func (db *LevelDB) Get(key string) ([]byte, error) {
	data, err := db.DB.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *LevelDB) Del(key string) error {
	return db.DB.Delete([]byte(key),nil)
}

func (db *LevelDB) State(value string )(string, error){
	if value==""{
		value = "leveldb.stats"
	}
	if value =="type"{
		return "leveldb",nil
	}

	return db.DB.GetProperty(value)
}

func (db *LevelDB) Iterator(prefix string) (map[string] string,error){
	data :=make(map[string] string)
	var iter iterator.Iterator
	if prefix==""{
		iter = db.DB.NewIterator(nil,nil)
		for ok:= iter.Seek([]byte(""));ok;ok=iter.Next(){
			data[string(iter.Key())] = string(iter.Value()[:])
		}
	}else {
		iter = db.DB.NewIterator(util.BytesPrefix([]byte(prefix)),nil)
		for iter.Next(){
			data[string(iter.Key())] = string(iter.Value()[:])
		}
	}
	iter.Release()
	return data, iter.Error()
}
