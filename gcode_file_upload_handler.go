package main

// func handleGcodeFileUpload(w http.ResponseWriter, r *http.Request) {
// 	file, _, err := r.FormFile("FILE")
// 	if err != nil {
// 		respondJSON(w, jd{
// 			"error": err.Error(),
// 		})
// 		log.Printf("Gcode upload error: %v", err)
// 		return
// 	}
// 	defer file.Close()

// 	// copy example
// 	f, err := os.OpenFile("./downloaded", os.O_WRONLY|os.O_CREATE, 0666)
// 	defer f.Close()
// 	io.Copy(f, file)
// }
