package arango

import "github.com/solher/arangolite"

var (
	imageStore ImageStore
)

func init() {
	dbName := "ims"
	endpoint := "http://localhost:8529"

	db := arangolite.Connect(endpoint, "_system", "", "")

	_, _ := db.Run(&arangolite.DropDatabase{Name: dbName})
	_, err := db.Run(&arangolite.CreateDatabase{Name: dbName})
	if err != nil {
		panic(err)
	}

	db.SwitchDatabase(dbName)

	for _, collection := range []string{imageCollection, tagCollection} {
		_, _ := db.Run(&arangolite.DropCollection{Name: collection})
		_, err := db.Run(&arangolite.CreateCollection{Name: collection})
		if err != nil {
			panic(err)
		}
	}

	_, _ := db.Run(&arangolite.DropGraph{Name: dbName, DropCollection: true})

	from := make([]string, 1)
	from[0] = imageCollection
	to := make([]string, 1)
	to[0] = tagCollection
	edgeDefinition := arangolite.EdgeDefinition{Collection: "image_tag", From: from, To: to}
	edgeDefinitions := make([]arangolite.EdgeDefinition, 1)
	edgeDefinitions[0] = edgeDefinition
	db.Run(&arangolite.CreateGraph{Name: dbName, EdgeDefinition: edgeDefinitions})

	imageStore = NewImageStore(db)
}
