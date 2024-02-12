package useCase

import (
	"io"
	"net/http"
	errorsCFG "priceComp/pkg/errors"
)

func downloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check that the response status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return nil, errorsCFG.ErrBadUrl
	}

	// Create a file to save the downloaded image
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Copy the body of the response to the file
	return imageData, nil
}
