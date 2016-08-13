package copydir

import (
	"fmt"
	"io"
	"log"
	"os"
)

func copyFile(source, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func copyDir(source, dest string) (err error) {

	// Get properties of source dir.
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create dest dir.

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively.
			err = copyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Perform copy.
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func Copy(source_dir, dest_dir string, overwrite bool) (err error) {

	log.Println("Source: " + source_dir)

	// Check if the source dir exist.
	src, err := os.Stat(source_dir)
	if err != nil {
		return err
	}
	if !src.IsDir() {
		return fmt.Errorf("Source is not a directory.")
	}

	log.Println("Destination: " + dest_dir)

	// We will continue to copy if we meet either condition:
	// 1. The destination does not exist.
	// 2. The destination exists and it is a dir and overwrite is true.
	dest, err := os.Stat(dest_dir)
	if err == nil {
		if !dest.IsDir() {
			return fmt.Errorf("Destination is not a directory.")
		}
		if !overwrite {
			return fmt.Errorf("We will not overwrite the destination.")
		}
	}

	err = copyDir(source_dir, dest_dir)
	if err != nil {
		return err
	}
	log.Println("Directory copied.")
	return nil
}
