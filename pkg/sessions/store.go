package sessions

// Store ...
type Store interface {
	Save(s *Session) error
	Get(sessID string) (*Session, error)
	GetLen(nameAndClient string) (int, error)
}

// MongoStore ...
type MongoStore struct {

}

// NewMongoStore ...
func NewMongoStore() *MongoStore {
	// TODO:
	return nil
}

// Save ...
func (ms *MongoStore) Save(s *Session) error {
	return nil
}

// Get ...
func (ms *MongoStore) Get(sessID string) (*Session, error) {
	return nil, nil
}

// GetLen ...
func (ms *MongoStore) GetLen(nameAndClient string) (int, error) {
	return 0, nil
}
