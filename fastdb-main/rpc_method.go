// rpc_methods.go

package fastdb

// RPCArgs represents the arguments for RPC methods.
type RPCArgs struct {
	Method  string
	Bucket  string
	Key     int
	Value   []byte
	SyncIime int
}

// RPCResult represents the result of RPC methods.
type RPCResult struct {
	Data []byte
	Ok   bool
	Info string
}

// FastDBRPC is a type that will hold the DB instance for RPC methods.
type FastDBRPC struct {
	db *DB
}

// OpenDB opens a database at the provided path.
func (f *FastDBRPC) OpenDB(args *RPCArgs, result *RPCResult) error {
	db, err := Open(args.Bucket, args.SyncIime)
	f.db = db
	result.Info = db.Info()
	return err
}

// DefragDB optimizes the file to reflect the latest state.
func (f *FastDBRPC) DefragDB(args *RPCArgs, result *RPCResult) error {
	err := f.db.Defrag()
	result.Info = f.db.Info()
	return err
}

// DelDB deletes one map value in a bucket.
func (f *FastDBRPC) DelDB(args *RPCArgs, result *RPCResult) error {
	ok, err := f.db.Del(args.Bucket, args.Key)
	result.Ok = ok
	result.Info = f.db.Info()
	return err
}

// GetDB returns one map value from a bucket.
func (f *FastDBRPC) GetDB(args *RPCArgs, result *RPCResult) error {
	data, ok := f.db.Get(args.Bucket, args.Key)
	result.Data = data
	result.Ok = ok
	result.Info = f.db.Info()
	return nil
}

// GetAllDB returns all map values from a bucket.
func (f *FastDBRPC) GetAllDB(args *RPCArgs, result *RPCResult) error {
	data, err := f.db.GetAll(args.Bucket)
	if err == nil {
		result.Data = make([]byte, 0, 0)
		for _, v := range data {
			result.Data = append(result.Data, v...)
		}
	}
	result.Info = f.db.Info()
	return err
}

// SetDB stores one map value in a bucket.
func (f *FastDBRPC) SetDB(args *RPCArgs, result *RPCResult) error {
	err := f.db.Set(args.Bucket, args.Key, args.Value)
	result.Info = f.db.Info()
	return err
}

// CloseDB closes the database.
func (f *FastDBRPC) CloseDB(args *RPCArgs, result *RPCResult) error {
	err := f.db.Close()
	result.Info = f.db.Info()
	return err
}
