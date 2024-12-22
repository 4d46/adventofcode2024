package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type File struct {
	blank bool
	id    int
	size  int
}

func main() {
	fmt.Println("Hello, day 09!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	{
		// Parse the input data
		diskmap := make([]File, 0, len(inputData[0]))
		isBlank := false
		fileNum := 0
		// Loop over each character of the first line
		for _, char := range inputData[0] {
			// Convert the character to an integer
			num, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatalf("failed to convert character to integer: %s", err)
			}
			// Don't add 0 length files or gaps
			if num > 0 {
				diskmap = append(diskmap, File{blank: isBlank, id: fileNum, size: num})
			}
			if !isBlank {
				fileNum++
			}
			isBlank = !isBlank
		}

		// Calculate the sum of the sizes of the files and blanks
		totalDiskSize := 0
		for _, file := range diskmap {
			totalDiskSize += file.size
		}
		fmt.Println("Total disk size:", totalDiskSize)

		printDiskmap(diskmap)
		// fmt.Println()

		// Looping over the diskmap, find the first file that is blank
		fileLoc := 0
	DONE:
		for {
			// If we reach the end of the diskmap, break out of the loop
			if fileLoc >= len(diskmap) {
				break
			}
			// If the file is blank, we want to replace it with data from files at the end of the diskmap
			if diskmap[fileLoc].blank {
				// Find the last file, if we find the last file is blank then remove it from the diskmap and try again
				lastFileLoc := len(diskmap) - 1
				for {
					// If there are no files left, break out of the loop
					if lastFileLoc <= fileLoc {
						break DONE
					}
					// If the last file is blank, remove it from the diskmap
					if diskmap[lastFileLoc].blank {
						// fmt.Printf("Length before last blank file trim %d\n", len(diskmap))
						diskmap = diskmap[:lastFileLoc]
						lastFileLoc--
						// fmt.Printf("Length after last blank file trim %d\n", len(diskmap))
					} else {
						// We've found a file that is not blank, stop searching
						break
					}
				}
				// spew.Dump(diskmap[fileLoc])
				// spew.Dump(diskmap[lastFileLoc])
				// fmt.Printf("diskmap length: %d\n", len(diskmap))
				if diskmap[fileLoc].size == diskmap[lastFileLoc].size {
					// File is the same size as the gap, replace the gap with the file
					diskmap[fileLoc] = diskmap[lastFileLoc]
					// Remove the last file
					diskmap = diskmap[:lastFileLoc]
				} else if diskmap[fileLoc].size < diskmap[lastFileLoc].size {
					// File is bigger than the gap, fill the gap and leave remainder of file at the end
					diskmap[fileLoc].id = diskmap[lastFileLoc].id
					diskmap[fileLoc].blank = false
					diskmap[lastFileLoc].size -= diskmap[fileLoc].size
				} else {
					// fmt.Printf("Length before gap shrink %d\n", len(diskmap))
					// File is smaller than the gap, insert file and shrink gap
					// Shrink the gap
					diskmap[fileLoc].size -= diskmap[lastFileLoc].size
					// Insert extra file in diskmap before the gap
					diskmap = slices.Insert(diskmap, fileLoc, diskmap[lastFileLoc])
					// Remove the last file
					diskmap = diskmap[:len(diskmap)-1]
					// fmt.Printf("Length after gap shrink %d\n", len(diskmap))
				}

			}
			// Move onto next file
			fileLoc++
			// printDiskmap(diskmap)
		}
		// Should all now be packed, calculate the sum of the sizes of the files
		// Then add blanks to the end of the diskmap to make it the same size as the original diskmap
		packedDiskSize := 0
		for _, file := range diskmap {
			packedDiskSize += file.size
		}
		fmt.Println("Packed disk size:", packedDiskSize)
		diskmap = append(diskmap, File{blank: true, size: totalDiskSize - packedDiskSize})

		printDiskmap(diskmap)
		// fmt.Println("0099811188827773336446555566..............")
		checksum := caclulateChecksum(diskmap)
		fmt.Printf(" Part 1 checksum: %d\n", checksum)
	}
	fmt.Println("-----------------------------")

	// Part 2
	{
		// Parse the input data
		diskmap := make([]File, 0, len(inputData[0]))
		isBlank := false
		fileNum := 0
		// Loop over each character of the first line
		for _, char := range inputData[0] {
			// Convert the character to an integer
			length, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatalf("failed to convert character to integer: %s", err)
			}
			// Don't add 0 length files or gaps
			if length > 0 {
				if isBlank {
					diskmap = append(diskmap, File{blank: isBlank, id: -1, size: length})
				} else {
					diskmap = append(diskmap, File{blank: isBlank, id: fileNum, size: length})
					fileNum++
				}
			}
			isBlank = !isBlank
		}

		// spew.Dump(diskmap)

		// Calculate the sum of the sizes of the files and blanks
		totalDiskSize := 0
		for _, file := range diskmap {
			totalDiskSize += file.size
		}
		fmt.Println("Total disk size:", totalDiskSize)

		printDiskmap(diskmap)

		// Starting from the last file and work forwards
		// lastFileLoc := len(diskmap) - 1

		// Start with the last file id and work backwards
		for fileId := diskmap[len(diskmap)-1].id; fileId >= 0; fileId-- {
			// fmt.Printf("Before moving file %d ", fileId)
			// printDiskmap(diskmap)
			// Find the file with the given id
			fileLoc := len(diskmap) - 1
			for {
				// If we reach the start of the diskmap without finding the file, raise an error
				if fileLoc < 0 {
					log.Fatalf("failed to find file with id: %d", fileId)
				}
				// If we find the file, break out of the loop
				if diskmap[fileLoc].id == fileId {
					break
				}
				fileLoc--
			}
			// spew.Dump(diskmap[fileLoc])
			fileSize := diskmap[fileLoc].size
			// fmt.Printf("Moving file %d of size %d\n", fileId, fileSize)

			// Start from the beginning of the diskmap and work forwards, trying yo find a gap that the file can fit into
			for gapLoc := 0; gapLoc < fileLoc; gapLoc++ {
				// If not a blank space then skip
				if !diskmap[gapLoc].blank {
					continue
				}
				// If the gap is too small then skip
				if diskmap[gapLoc].size < fileSize {
					continue
				}
				// If the gap is the right size then move the file
				if diskmap[gapLoc].size == fileSize {
					diskmap[gapLoc] = diskmap[fileLoc]
					// Replace the file with blank space
					// But we don't want to create multiple adjacent blank spaces, so check either side of the file for blanks
					// First convert the file to a blank
					diskmap[fileLoc].blank = true
					diskmap[fileLoc].id = -1
					// Check the file after the fileLoc, if it isn't at the end of the diskmap
					if fileLoc < len(diskmap)-1 {
						if diskmap[fileLoc+1].blank {
							// If the file after the fileLoc is blank, expand the blank space and remove that file
							diskmap[fileLoc].size += diskmap[fileLoc+1].size
							diskmap = append(diskmap[:fileLoc+1], diskmap[fileLoc+2:]...)
						}
					}
					// Check the file before the fileLoc, if it isn't at the start of the diskmap
					if fileLoc > 0 {
						if diskmap[fileLoc-1].blank {
							// If the file before the fileLoc is blank, expand the blank space and remove this file
							diskmap[fileLoc-1].size += diskmap[fileLoc].size
							diskmap = append(diskmap[:fileLoc], diskmap[fileLoc+1:]...)
						}
					}
					break
				}
				// If the gap is bigger than the file then move the file and shrink the gap
				if diskmap[gapLoc].size > fileSize {
					diskmap[gapLoc].size -= fileSize
					// Take a copy of the file, then convert the file in the diskmap to a blank
					tempfile := diskmap[fileLoc]
					diskmap[fileLoc].blank = true
					diskmap[fileLoc].id = -1
					// Check the file after the fileLoc, if it isn't at the end of the diskmap
					if fileLoc < len(diskmap)-1 {
						if diskmap[fileLoc+1].blank {
							// If the file after the fileLoc is blank, expand the blank space and remove that file
							diskmap[fileLoc].size += diskmap[fileLoc+1].size
							diskmap = append(diskmap[:fileLoc+1], diskmap[fileLoc+2:]...)
						}
					}
					// Check the file before the fileLoc, if it isn't at the start of the diskmap
					if fileLoc > 0 {
						if diskmap[fileLoc-1].blank {
							// If the file before the fileLoc is blank, expand the blank space and remove this file
							diskmap[fileLoc-1].size += diskmap[fileLoc].size
							diskmap = append(diskmap[:fileLoc], diskmap[fileLoc+1:]...)
						}
					}
					// Now insert the file into the diskmap at the gap location
					diskmap = slices.Insert(diskmap, gapLoc, tempfile)
					break
				}
			}
			// lastFileLoc--
		}
		printDiskmap(diskmap)
		// fmt.Println("00992111777.44.333....5555.6666.....8888..")
		checksum := caclulateChecksum(diskmap)
		fmt.Printf(" Part 2 checksum: %d\n", checksum)
	}

}

func caclulateChecksum(diskmap []File) int {
	// Loop over total length of diskmap
	// If the file is blank, skip it
	pos := 0
	checksum := 0
	for _, file := range diskmap {
		if !file.blank {
			for i := 0; i < file.size; i++ {
				checksum += (pos * file.id)
				pos++
			}
		} else {
			pos += file.size
		}
	}
	return checksum
}

func printDiskmap(diskmap []File) {
	for _, file := range diskmap {
		if file.blank {
			fmt.Print(strings.Repeat(".", file.size))
		} else {
			fmt.Print(strings.Repeat(fmt.Sprintf("%d", file.id), file.size))
		}
	}
	fmt.Println()
}

// LoadInputData reads the input file and returns a slice of strings
func LoadInputData(filename string) []string {
	// Read the file as strings a line at a time
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new reader
	reader := bufio.NewReader(file)

	var inputData []string

	for {
		// Read until we encounter a newline character
		line, err := reader.ReadString('\n')
		if err != nil {
			// If we encounter EOF, break out of the loop
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("error reading file: %s", err)
		}
		// Remove the newline character from the end of the line
		line = strings.TrimSuffix(line, "\n")
		// Append the line to the inputData slice
		inputData = append(inputData, line)
	}
	return inputData
}
