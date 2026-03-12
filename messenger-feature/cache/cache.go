package cache

type Cache interface {
	AddMessage(message string) error

	GetMessages() ([]string, error)
	GetMessageByID(id string) (string, error)
	
	DeleteMessage(message string) error
	DeleteMessageByID(id string) error
}