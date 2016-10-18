package cmd

import (
	"fmt"
	"io"
	"math"
	"os"
)

//readFileChunk will read a file from the given path, treat it as though it has been broken into chunks, and read the specified chunk into an io.Writer
//if it fails to read a chunk directly from the file, it will try to find a part file matching the chunk number and write that to the writer instead
//it returns the number of chunks it created and an error
func chunkFile(fPth string) (int, error) {
	f, err := os.Open(fPth)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return -1, err
	}
	fSize := fInfo.Size()
	chunkNum := 0
	for ; chunkNum < int(math.Ceil(float64(fSize)/float64(chunkSize))); chunkNum++ {
		out, err := os.Create(fmt.Sprintf("%s.part%d", fPth, chunkNum))
		if err != nil {
			return -1, err
		}
		//File.Seek sets the offset from which to read
		_, err = f.Seek(int64(chunkNum)*chunkSize, 0)
		if err != nil {
			return -1, err
		}
		//get file size to see if we need to read less than 1 chunk
		if int64(chunkNum+1) > fSize/chunkSize {
			_, err = io.CopyN(out, f, fSize%chunkSize)
		} else {
			_, err = io.CopyN(out, f, chunkSize)
		}
		out.Close()
	}
	return chunkNum, err
}

//readPartFile tries to read a file of the name "[path]/[name].[ext].partN" into an io.Writer
func readPartFile(fPth string, chunkNum int, conn io.Writer) error {
	partPath := fmt.Sprintf("%s.part%d", fPth, chunkNum)
	f, err := os.Open(partPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(conn, f)
	return err
}

func stitchOriginal(fPth string) error {
	out, err := os.Create(fPth)
	if err != nil {
		return err
	}
	for i := 0; ; i++ {
		f, err := os.Open(fmt.Sprintf("%s.part%d", fPth, i))
		if err != nil {
			if os.IsNotExist(err) && i != 0 {
				return nil
			}
			return err
		}
		io.Copy(out, f)
	}
}

func recvFileChunk(fPth string, chunkNum int, conn io.Reader) error {
	f, err := os.Create(fmt.Sprintf("%s.part%d", fPth, chunkNum))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, conn)
	return err
}
