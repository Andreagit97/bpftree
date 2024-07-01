package capture

// var (
// 	// All the capture files will be written in little endian.
// 	captureEndianness binary.ByteOrder = binary.LittleEndian
// )

// // GetCaptureEndianness returns the endianness used to read and write.
// func GetCaptureEndianness() binary.ByteOrder {
// 	return captureEndianness
// }

// func getTasksFromFile() ([]*task.Task, error) {

// 	// todo!: implement this function

// 	// /* 1. Parse Magic number */
// 	// if err := parseMagicNumber(reader); err != nil {
// 	// 	displayError("cannot parse the magic number:", err)
// 	// 	return err
// 	// }

// 	// /* 2. Parse Header len */
// 	// if err := parseHeaderLen(reader); err != nil {
// 	// 	displayError("cannot parse the header len:", err)
// 	// 	return err
// 	// }

// 	// /* 3. Parse Header */
// 	// if err := parseHeader(reader); err != nil {
// 	// 	displayError("cannot parse the header:", err)
// 	// 	return err
// 	// }

// 	// /* 4. Parse task Info. */
// 	// if err := parseTaskInfos(reader); err != nil {
// 	// 	displayError("unable to parse task info:", err)
// 	// 	return err
// 	// }

// 	// /* Order the list of tasks by tid */
// 	// sort.SliceStable(tasksList, func(i, j int) bool {
// 	// 	return tasksList[i].Info.Tid < tasksList[j].Info.Tid
// 	// })
// 	return nil, nil
// }

// func GetAllTasks() ([]*task.Task, error) {
// 	// todo!: we need to implement this
// 	return nil, fmt.Errorf("not implemented")
// }

// // // SetCaptureFilePath is used to set the capture file.
// // func SetCaptureFilePath(filepath string) {
// // 	captureFilePath = filepath
// // }

// // func GetCaptureFilePath() string {
// // 	return captureFilePath
// // }

// /* todo!: today we dump all the files and all the tasks but maybe in the future we could dump only some of them.
//  *
//  * Format we dump the file:
//  * - Major version (uint64)
//  * - Minor version (uint64)
//  * - Patch version (uint64)
//  * -------------------------------- The following schema could change for each version of the tool.
//  * - Size of the task struct (packed) (uint64)
//  * - Number of tasks (uint64)
//  * - Tasks
//  * --------------------------------
//  * - Size of the file struct (packed) (uint64)
//  * - Task id (int32)
//  * - Number of files (uint32)
//  * - Files
//  * - ... (repeat for each task)
//  */

// func dumpTasksIntoFile(path string) error {

// 	// todo!: implement this function

// 	// absPath, err := filepath.Abs(path)
// 	// if err != nil {
// 	// 	render.DisplayError(err)
// 	// 	return err
// 	// }

// 	// f, err := os.Create(filepath.Clean(absPath))
// 	// if err != nil {
// 	// 	render.DisplayError(err)
// 	// 	return err
// 	// }
// 	// defer f.Close()

// 	// // Check if the kernel supports necessary BPF features.
// 	// if err := utils.CheckeBPFSupport(); err != nil {
// 	// 	render.DisplayError(err)
// 	// 	return err
// 	// }

// 	// // We write the version of the tool that dumped the file.
// 	// // Only a tool with the same major and minor and with a greater or equal patch can read the file.
// 	// if err = binary.Write(f, utils.GetCaptureEndianness(), utils.GetVersionMajor()); err != nil {
// 	// 	render.DisplayError("cannot write major version: %v", err)
// 	// 	return err
// 	// }

// 	// if err = binary.Write(f, utils.GetCaptureEndianness(), utils.GetVersionMinor()); err != nil {
// 	// 	render.DisplayError("cannot write minor version: %v", err)
// 	// 	return err
// 	// }

// 	// if err = binary.Write(f, utils.GetCaptureEndianness(), utils.GetVersionPatch()); err != nil {
// 	// 	render.DisplayError("cannot write patch version: %v", err)
// 	// 	return err
// 	// }

// 	// // open BPF iterator.
// 	// reader, err := task.OpenBPFIterator()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer reader.Close()

// 	// /* Read chuncks of the file. */
// 	// chunckBytes := make([]byte, fileChunck)
// 	// for {
// 	// 	readBytes, err := reader.Read(chunckBytes)
// 	// 	if errors.Is(err, io.EOF) {
// 	// 		/* We read all the file. */
// 	// 		break
// 	// 	}
// 	// 	if err != nil {
// 	// 		displayError("unable to read file chunck: %v", err)
// 	// 		return err
// 	// 	}

// 	// 	if _, err := f.Write(chunckBytes[:readBytes]); err != nil {
// 	// 		displayError("unable to write file chunck: %v", err)
// 	// 		return err
// 	// 	}
// 	// }
// 	// displayGraph(imageNewspaper, "Capture correctly dumped:", absPath)
// 	return nil
// }
