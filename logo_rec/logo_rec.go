package logorec

import (
	"context"
	"fmt"
	"io"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

func RecognizeLogo(file io.Reader) ([]string, error) {
	sb := strings.Builder{}
	w := io.Writer(&sb)
	err := detectLogos(w, file)
	if err != nil {
		return nil, err
	}
	names := strings.Split(sb.String(), "\n")
	return names, nil

}

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogos(w io.Writer, file io.Reader) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

func detectLogosURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}
