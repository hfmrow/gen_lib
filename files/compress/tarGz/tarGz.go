// tarGz.go

// Source file auto-generated on Sun, 31 Mar 2019 19:42:54 using Gotk3ObjHandler ©2019-21 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

// CreateTarball: create tar.gz file from given filenames.
// UntarGzip: unpack tar.gz files list. Return len(filesList)=0 if all files has been restored.

package tarGz

import (
	"archive/tar"
	"compress/flate"
	"sort"

	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	glfsfo "github.com/hfmrow/gen_lib/files/filesOperations"
	gzip "github.com/klauspost/pgzip"
	"github.com/ulikunitz/xz"
)

var (
	countFiles int

	// Lib mapping
	ChangeFileOwner = glfsfs.ChangeFileOwner
)

// CreateTarball: create tar.gz file from given filenames.
// -2 = HuffmanOnly (linear compression, low gain, fast compression)
// -1 = DefaultCompression
//  0 = NoCompression
//  1 -> 9 = BestSpeed -> BestCompression
func CreateTarballGzip(tarballFilename string, filenames []string, compressLvl int) (countedWrittenFiles int, err error) {
	countFiles = 0
	file, err := os.Create(tarballFilename)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilename, err.Error()))
	}
	defer file.Close()

	zipWriter, err := gzip.NewWriterLevel(file, compressLvl)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Bad compression level '%d', got error '%s'", int(flate.BestCompression), err.Error()))
	}
	defer zipWriter.Close()

	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()
	// currentStoreFiles.inUse = true
	for _, filePath := range filenames {
		// select {
		// case <-quitGoRoutine:
		// quitGoRoutine = make(chan struct{})
		// err = changeFileOwner(tarballFilename)
		// currentStoreFiles.inUse = false
		// return countFiles, errors.New(sts["userCancelled"])
		// default:
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return countFiles, errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
		}
		// }
	}
	err = ChangeFileOwner(tarballFilename)
	// currentStoreFiles.inUse = false
	return countFiles, err
}

// addFileToTarWriter:
func addFileToTarWriter(filePath string, tarWriter *tar.Writer) (err error) {
	var stat os.FileInfo
	var linkname string

	stat, err = os.Lstat(filePath)
	modeType := stat.Mode() & os.ModeType
	switch {
	case modeType&os.ModeNamedPipe > 0:
		return nil
	case modeType&os.ModeSocket > 0:
		return nil
	case modeType&os.ModeDevice > 0:
		return nil
	case err != nil:
		return nil
	}

	file, err := os.Open(filePath)
	switch {
	case os.IsPermission(err):
		return nil
	default:
		defer file.Close()
	}

	if link, err := os.Readlink(filePath); err == nil {
		linkname = link
	}

	header, err := tar.FileInfoHeader(stat, filepath.ToSlash(linkname))
	if err != nil {
		return errors.New(fmt.Sprintf("Could not build header for '%s', got error '%s'", filePath, err.Error()))
	}
	header.Name = filePath

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}
	countFiles++
	if header.Typeflag == tar.TypeReg {
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
		}
	}
	return nil
}

// UntarGzip: unpack tar.gz files list. Return len(filesList)=0 if all files has been restored.
// removeDir means that, before restore folder content, the existing dir will be removed.
func UntarGzip(sourcefile, dst string, filesList *[]string, removeDir bool) (countedWrittenFiles int, err error) {
	var storePath, target string
	var header *tar.Header
	var tmpFilesList = *filesList
	var skipReadHeader bool
	// Initialise readers
	file, err := os.Open(sourcefile)
	if err != nil {
		return countedWrittenFiles, err
	}
	defer file.Close()

	var r io.ReadCloser = file

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return countedWrittenFiles, err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	sort.SliceStable(tmpFilesList, func(i, j int) bool {
		return tmpFilesList[i] < tmpFilesList[j]
	})

	// ownThis: set file/dir owner owner stored in tar archive.
	var ownThis = func(target string) (err error) {
		err = os.Chown(target, header.Uid, header.Gid)
		if !os.IsPermission(err) {
			if err != nil {
				return err
			}
		}
		return nil
	}
	// buildDirStruct: and set owner respective for each dir created.
	var buildDirStruct = func(target string) (err error) {
		if _, err := os.Stat(target); os.IsNotExist(err) {
			if err := os.MkdirAll(target, os.ModePerm); err == nil {
				for {
					if err = ownThis(target); err != nil {
						return err
					}
					target = filepath.Dir(target)
					if target == filepath.Dir(dst) {
						break
					}
				}
			}
		}
		return err
	}
	// restoreFile: Handling Directory, regular file, symlink and own them. Others
	// kind of files are ignored cause my packing function don't handle them.
	var restoreFile = func(target string) (err error) {
		countedWrittenFiles++
		switch header.Typeflag {
		case tar.TypeDir:
			if removeDir {
				if _, err := os.Stat(target); !os.IsNotExist(err) {
					err = os.RemoveAll(target)
					if err != nil {
						return err
					}
				}
			}
			err = buildDirStruct(target)
			if err != nil {
				return err
			}
			// fmt.Println(header.Name)
		case tar.TypeSymlink:
			if _, err := os.Lstat(target); err == nil {
				os.Remove(target)
			}
			err = os.Symlink(header.Linkname, target)
			if err != nil {
				return err
			}
			// fmt.Println(header.Name)
		case tar.TypeReg:
			err = buildDirStruct(filepath.Dir(target)) // in case of regular file where there is no directory to receive it.
			if err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			if err = ownThis(target); err == nil {
				f.Close()
				return err
			}
			f.Close()
			// fmt.Println(header.Name)
		}
		return err
	}
	// readHeader:
	var readHeader = func() (target string, err error) {
		header, err = tr.Next()
		switch {
		case err == io.EOF:
			return "", nil
		case err != nil:
			return "", nil
		}
		target = filepath.Join(dst, header.Name)
		return target, err
	}
	// writeFile: and dir, set perms and owner stored in tar archive.
	var writeFile = func(target *string) (err error) {
		if header != nil {
			if header.Typeflag == tar.TypeDir {
				storePath = header.Name
				for {
					// select {
					// case <-quitGoRoutine:
					// 	quitGoRoutine = make(chan struct{})
					// 	// currentStoreFiles.inUse = false
					// 	header = nil
					// 	return errors.New(sts["userCancelled"])
					// default:
					err = restoreFile(*target)
					if err != nil {
						return err
					}
					*target, err = readHeader()
					if err != nil {
						return err
					}
					if header == nil || !strings.Contains(filepath.Dir(header.Name), storePath) {
						return nil
					}
					// }
				}
			} else {
				err = restoreFile(*target)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Parse desired files and restore them.
	// currentStoreFiles.inUse = true
	for {
		if !skipReadHeader {
			target, err = readHeader()
		}
		if err != nil || header == nil || len(target) == 0 {
			return countedWrittenFiles, err
		}
		if len(tmpFilesList) != 0 {
			for idx := 0; idx < len(tmpFilesList); idx++ {
				source := filepath.Join(dst, tmpFilesList[idx])
				if source == target && len(tmpFilesList) != 0 {
					tmpFilesList = append(tmpFilesList[:idx], tmpFilesList[idx+1:]...)
					idx--
					if err = writeFile(&target); err != nil {
						return countedWrittenFiles, err
					}
					if len(tmpFilesList) == 0 {
						*filesList = tmpFilesList
						return countedWrittenFiles, err
					}
				}
				skipReadHeader = false
				if idx == -1 || filepath.Join(dst, tmpFilesList[idx]) == target {
					skipReadHeader = true
				}
			}
		} else {
			if err = writeFile(&target); err != nil {
				return countedWrittenFiles, err
			}
		}
	}
	// End of game ...
	*filesList = tmpFilesList
	// currentStoreFiles.inUse = false
	return countedWrittenFiles, err
}

// UnGzip: Unpack
func UnGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

// Gzip: Pack file
func Gzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	//	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

// Tar: Make standalone tarball
func Tar(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Lstat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if baseDir != "" {
			//			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			header.Name = path
		}

		if link, err := os.Readlink(path); err == nil {
			header.Linkname = link
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if header.Typeflag == tar.TypeReg {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
		}
		return err
	})
}

// Untar: extract file from tarball
func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open tar.gz file '%s', got error '%s'", tarball, err.Error()))
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New(fmt.Sprintf("Could not read file in tar.gz archive, got error '%s'", err.Error()))
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTarballLXz:  create tar.xz file from given filenames. Just for fun test ...
// but it's really slow, no multi-threading and low size gain.
func CreateTarballXz(tarballFilename string, filenames []string, compressLvl int) (countedWrittenFiles int, err error) {
	countFiles = 0
	file, err := os.Create(tarballFilename)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilename, err.Error()))
	}
	defer file.Close()

	xzWriter, err := xz.NewWriter(file)
	if err != nil {
		return countFiles, errors.New(fmt.Sprintf("Could not create xzWriter, got error '%s'", err.Error()))
	}
	defer xzWriter.Close()
	tarWriter := tar.NewWriter(xzWriter)
	defer tarWriter.Close()
	// currentStoreFiles.inUse = true
	for _, filePath := range filenames {
		// select {
		// case <-quitGoRoutine:
		// 	quitGoRoutine = make(chan struct{})
		// 	err = changeFileOwner(tarballFilename)
		// 	currentStoreFiles.inUse = false
		// 	return countFiles, errors.New(sts["userCancelled"])
		// default:
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return countFiles, errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
			// }
		}
	}
	err = ChangeFileOwner(tarballFilename)
	// currentStoreFiles.inUse = false
	return countFiles, err
}
