package task

// const (
// 	fileChunck int = 40
// )

// todo!: we need to move this

// func DumpTasksIntoFile(path string) error {
// 	absPath, err := filepath.Abs(path)
// 	if err != nil {
// 		displayError(err)
// 		return err
// 	}

// 	f, err := os.Create(filepath.Clean(absPath))
// 	if err != nil {
// 		displayError(err)
// 		return err
// 	}
// 	defer f.Close()

// 	/* open BPF iterator. */
// 	reader, err := openBPFIterator()
// 	if err != nil {
// 		return err
// 	}
// 	defer reader.Close()

// 	// todo!: we should write some info at the beginning of the file

// 	/* Read chuncks of the file. */
// 	chunckBytes := make([]byte, fileChunck)
// 	for {
// 		readBytes, err := reader.Read(chunckBytes)
// 		if errors.Is(err, io.EOF) {
// 			/* We read all the file. */
// 			break
// 		}
// 		if err != nil {
// 			displayError("unable to read file chunck: %v", err)
// 			return err
// 		}

// 		if _, err := f.Write(chunckBytes[:readBytes]); err != nil {
// 			displayError("unable to write file chunck: %v", err)
// 			return err
// 		}
// 	}
// 	displayGraph(imageNewspaper, "Capture correctly dumped:", absPath)
// 	return nil
// }
