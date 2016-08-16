// Package copydir is a tool for copying contents of a directory to a
// destination directory, including all files and sub-directories.
package copydir

import (
	"fmt"
	"io"
	"log"
	"os"
)

func copyFile(source, dest string, ch chan error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		ch <- err
		return
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		ch <- err
		return
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	ch <- err
}

func copyDir(source, dest string, ch chan error) {

	// Get properties of source dir.
	sourceinfo, err := os.Stat(source)
	if err != nil {
		ch <- err
		return
	}

	// Create dest dir.

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		ch <- err
		return
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively.
			ch := make(chan error)
			go copyDir(sourcefilepointer, destinationfilepointer, ch)
			if err = <-ch; err != nil {
				log.Println(err)
			}
		} else {
			// Perform copy.
			ch := make(chan error)
			go copyFile(sourcefilepointer, destinationfilepointer, ch)
			if err = <-ch; err != nil {
				log.Println(err)
			}
		}

	}
	ch <- err
}

// Copy executes copying contents from dir to dir. Overwriting is optional.
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

	ch := make(chan error)
	go copyDir(source_dir, dest_dir, ch)
	if err = <-ch; err != nil {
		return err
	}
	log.Println("Directory copied.")
	return nil
}
