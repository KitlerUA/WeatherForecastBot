package indexbuilder

import (
	"context"

	"github.com/KitlerUA/WeatherForecastBot/db"
)

const mapping = `{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	}
}`

func BuildIndices(indices ...string) error {
	ctx := context.Background()
	for i := range indices {
		if exists, _ := db.Get().IndexExists(indices[i]).Do(ctx); !exists {
			_, err := db.Get().CreateIndex(indices[i]).BodyString(mapping).Do(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
