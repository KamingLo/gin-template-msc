// DONT CHANGE THIS CODE, ITS FOR AUTOMIGRATION
// WITHOUT CHANGING CODE IN config/database.go
package models

var ModelsRegistry []interface{}

func RegisterModel(model interface{}) {
	ModelsRegistry = append(ModelsRegistry, model)
}
