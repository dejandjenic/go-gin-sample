package entities

type Todo struct {
	Name string `firestore:"name,omitempty"`
	ID   string `firestore:"-"`
}
