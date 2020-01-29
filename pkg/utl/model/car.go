package gorsk

// Car represents car model
type Car struct {
	Base
	Name      string     `json:"name"`
	Owner     User       `json:"owner"`
}