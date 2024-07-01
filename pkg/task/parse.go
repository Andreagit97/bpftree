package task

// /* If we want to read taskInfo from a file. */
// var (
// 	captureFilePath                  = ""
// 	byteOrder       binary.ByteOrder = binary.LittleEndian
// )

// func decodeByteEndianness(reader io.ReadCloser, numBytes uint32, data any) error {
// 	rawBytes := make([]byte, numBytes)
// 	readBytes, err := reader.Read(rawBytes)
// 	if errors.Is(err, io.EOF) {
// 		/* The caller should check for end of file */
// 		return io.EOF
// 	}
// 	if err != nil {
// 		displayError(fmt.Sprintf("cannot read '%d' bytes: %v", numBytes, err))
// 		return err
// 	}
// 	if uint32(readBytes) != numBytes {
// 		err = fmt.Errorf("read '%d' bytes instead of '%d'", readBytes, numBytes)
// 		return err
// 	}

// 	byteBuffer := bytes.NewBuffer(rawBytes)
// 	err = binary.Read(byteBuffer, byteOrder, data)
// 	if err != nil {
// 		displayError(fmt.Sprintf("cannot decode '%d' bytes: %v", numBytes, err))
// 		return err
// 	}
// 	return nil
// }

// func parseMagicNumber(reader io.ReadCloser) error {
// 	var magic uint16
// 	if err := decodeByteEndianness(reader, 2, &magic); err != nil {
// 		return err
// 	}

// 	// Swap endianness if necessary.
// 	if magic != magicLittleEndian {
// 		if magic == magicBigEndian {
// 			byteOrder = binary.BigEndian
// 		} else {
// 			return fmt.Errorf("the provided file doesn't have the right format")
// 		}
// 	}
// 	return nil
// }

// func OpenTaskIterator() (io.ReadCloser, error) {

// 	// Load pre-compiled programs and maps into the kernel.
// 	objs := iterObjects{}
// 	if err := loadIterObjects(&objs, nil); err != nil {
// 		displayError("unable to load BPF objects into the kernel:", err)
// 		return nil, err
// 	}
// 	defer objs.Close()

// 	opts := link.IterOptions{
// 		Program: objs.DumpTask,
// 	}

// 	iter, err := link.AttachIter(opts)
// 	if err != nil {
// 		displayError("unable to attach Iter prog:", err)
// 		return nil, err
// 	}

// 	reader, err := iter.Open()
// 	if err != nil {
// 		displayError("unable to open a reader for the iter:", err)
// 		return nil, err
// 	}
// 	return reader, nil
// }

// func openCaptureFile() (io.ReadCloser, error) {
// 	reader, err := os.Open(captureFilePath)
// 	if err != nil {
// 		displayError(fmt.Sprintf("unable to open file '%s': %v", captureFilePath, err))
// 		return nil, err
// 	}

// 	// We need to check that everything is correct here.
// 	// The len

// 	return reader, nil
// }

// // PopulateTaskInfo populate the thread table reading from a capture file
// // or using the BPF iterator.
// func PopulateTaskInfo() error {
// 	var err error
// 	var reader io.ReadCloser
// 	if captureFilePath == "" {
// 		// We need to read from the system with BPF.
// 		reader, err = openBPFIterator()
// 	} else {
// 		// We need to use the provided file.
// 		reader, err = openCaptureFile()
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	defer reader.Close()

// 	/* 1. Parse Magic number */
// 	if err := parseMagicNumber(reader); err != nil {
// 		displayError("cannot parse the magic number:", err)
// 		return err
// 	}

// 	/* 2. Parse Header len */
// 	if err := parseHeaderLen(reader); err != nil {
// 		displayError("cannot parse the header len:", err)
// 		return err
// 	}

// 	/* 3. Parse Header */
// 	if err := parseHeader(reader); err != nil {
// 		displayError("cannot parse the header:", err)
// 		return err
// 	}

// 	/* 4. Parse task Info. */
// 	if err := parseTaskInfos(reader); err != nil {
// 		displayError("unable to parse task info:", err)
// 		return err
// 	}

// 	/* Order the list of tasks by tid */
// 	sort.SliceStable(tasksList, func(i, j int) bool {
// 		return tasksList[i].Info.Tid < tasksList[j].Info.Tid
// 	})

// 	/* Compute children and order them by tid */
// 	computeChildren()
// 	return nil
// }

// func PrintTaskFiles(task string) {
// 	thread_id, err := strconv.Atoi(task)
// 	if err != nil {
// 		displayError("invalid thread id:", err)
// 		os.Exit(1)
// 	}

// 	/* Detect necessary BPF features */
// 	if err := bpfDetectionFeatures(); err != nil {
// 		displayError("cannot detect right features:", err)
// 		os.Exit(1)
// 	}

// 	/* Load pre-compiled programs and maps into the kernel. */
// 	spec, err := loadFileIter()
// 	if err != nil {
// 		displayError("cannot load the file iter:", err)
// 		os.Exit(1)
// 	}

// 	// we need to rewrite the map
// 	err = spec.RewriteConstants(map[string]interface{}{
// 		"target_thread_id": int32(thread_id),
// 	})
// 	if err != nil {
// 		displayError("cannot rewrite the map:", err)
// 		os.Exit(1)
// 	}

// 	objs := fileIterObjects{}
// 	if err := spec.LoadAndAssign(&objs, nil); err != nil {
// 		displayError("unable to load BPF iter file objects into the kernel:", err)
// 		os.Exit(1)
// 	}
// 	defer objs.Close()

// 	opts := link.IterOptions{
// 		Program: objs.DumpTaskFile,
// 	}

// 	iter, err := link.AttachIter(opts)
// 	if err != nil {
// 		displayError("unable to attach Iter file prog:", err)
// 		os.Exit(1)
// 	}

// 	reader, err := iter.Open()
// 	if err != nil {
// 		displayError("unable to open a reader for the iter:", err)
// 		os.Exit(1)
// 	}
// 	defer reader.Close()

// 	var fdInfoList []FileInfo
// 	for {
// 		fdInfo, err := parseFileInfo(reader)
// 		if err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}
// 			displayError("failed while reading:", err)
// 			os.Exit(1)
// 		}
// 		fdInfoList = append(fdInfoList, fdInfo)
// 	}

// 	displayGraph("Files for thread:", thread_id)
// 	for _, f := range fdInfoList {
// 		displayGraph(f.print())
// 	}
// }
