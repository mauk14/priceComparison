package main

import (
	"fmt"
	"priceComp/pkg/config"
	logger2 "priceComp/pkg/logger"
	"priceComp/pkg/postgres"
	"priceComp/services/DataCollection/internal/delivery/httpDataCollection"
	"priceComp/services/DataCollection/internal/repository"
	"priceComp/services/DataCollection/internal/useCase"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.MongoDb_uri)
	logger := logger2.SetUpLogger()
	db, err := postgres.OpenDb(cfg.PostgresDsn)
	if err != nil {
		logger.LogError(nil, err)
		return
	}

	//products, err := useCase.DnsParseManyProducts("https://www.dns-shop.kz/catalog/17a8a01d16404e77/smartfony/?order=6&stock=now-today-tomorrow-later")
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//
	//fmt.Println(products)
	app := httpDataCollection.NewApp(useCase.New(repository.New(db)), logger)

	err = app.Route().Run(fmt.Sprintf(":%d", cfg.DataCollection))
	if err != nil {
		logger.LogError(nil, err)
		return
	}

	//imageFilePath := "./services/DataCollection/cmd/app/airpods.jpg"
	//imageFile, err := os.Open(imageFilePath)
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//defer imageFile.Close()
	//
	//imageBytes, err := io.ReadAll(imageFile)
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//
	//sqlStatement := `
	//	INSERT INTO images (image_data, main_image, product_id)
	//	VALUES ($1, false, 1)
	//	RETURNING id
	//`
	//
	//var id int
	//err = db.QueryRow(context.Background(), sqlStatement, imageBytes).Scan(&id)
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//
	//fmt.Println("Image saved with ID:", id)
	//
	//sqlStatement = `
	//	SELECT image_data FROM images WHERE id = $1
	//`
	//
	//var imageData []byte
	//err = db.QueryRow(context.Background(), sqlStatement, id).Scan(&imageData)
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//imageFilePath = "retrieved_image.jpg"
	//file, err := os.Create(imageFilePath)
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = io.Copy(file, bytes.NewReader(imageData))
	//if err != nil {
	//	logger.LogError(nil, err)
	//	return
	//}
	//
	//fmt.Println("Image retrieved and saved to:", imageFilePath)
}
