package main

import (
	"fmt"
	"imageServer/primitiveUtil"
	"io"
	"net/http"
	"path/filepath"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
		<body>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<input type="file" name="image" />
			<button type="submit">
				Upload
			</button>
		</form>
		</body>
		</html>

		`
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)[1:]
		out, err := primitiveUtil.Tranform(file, ext, 5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch ext {
		case "jpeg":
			fallthrough
		case "jpg":
			w.Header().Set("Content-Type", "image/jpeg")
		case "png":
			w.Header().Set("Content-Type", "image/png")
		default:
			http.Error(w, "invalid image type", http.StatusBadRequest)
		}

		io.Copy(w, out)
	})
	// inFile, err := os.Open(".\\test2.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer inFile.Close()
	// out, err := primitiveUtil.Tranform(inFile, 10)
	// if err != nil {
	// 	fmt.Println(fmt.Sprint(err))
	// 	panic(err)
	// }
	// // fmt.Println(string(b))
	// os.Remove("out.jpg")
	// outFile, err := os.Create("out.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer outFile.Close()

	// io.Copy(outFile, out)
	http.ListenAndServe(":3000", mux)
}
