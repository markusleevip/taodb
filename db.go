package taodb

type DB interface{
	Set(string,[]byte) error
	Get(string) ([]byte,error)
	Del(string) error
	State(string)  (string, error)
	Iterator(prefix string) (map[string] string,error)
}
