package mongo

func GetMessages() {
	collection := Mongo.Database("chat").Collection("message")
	_ = collection
}
