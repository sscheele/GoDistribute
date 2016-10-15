package cmd

import (
	"fmt"
	"io"
	"os"
)

//readFileChunk will read a file from the given path, treat it as though it has been broken into chunks, and read the specified chunk into an io.Writer
func readFileChunk(fPth string, chunkNum int, conn io.Writer) error {
	f, err := os.Open(fPth)
	if err != nil {
		return readPartFile(fPth, chunkNum, conn)
	}
	//File.Seek sets the offset from which to read
	_, err = f.Seek(int64(chunkNum)*chunkSize, 0)
	if err != nil {
		return readPartFile(fPth, chunkNum, conn)
	}
	//get file size to see if we need to read less than 1 chunk
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if chunkNum > fInfo.Size()/chunkSize {
		_, err = io.CopyN(conn, f, int64(fInfo.Size())%chunkSize)
	} else {
		_, err = io.CopyN(conn, f, chunkSize)
	}
	return err
}

//readPartFile tries to read a file of the name "[path]/[name].[ext].partN" into an io.Writer
func readPartFile(fPth string, chunkNum int, conn io.Writer) error {
	partPath := fmt.Sprintf("%s.part%d", fPth, chunkNum)
	f, err := os.Open(partPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(conn, f)
	return err
}

func recvFileChunk(fPth string, chunkNum int, conn io.Reader) error {
	f, err := os.OpenFile(fmt.Sprintf("%s.part%d", fPth, chunkNum), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, conn)
	return err
}
