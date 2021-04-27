// gitCmdTerm.go

// Source file auto-generated on Fri, 27 Sep 2019 19:41:41 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	gcr "github.com/hfmrow/gen_lib/crypto"
	gst "github.com/hfmrow/gen_lib/strings"
)

/***********************************************/
/* Github acces via "git" command (terminal). */
/*********************************************/
type gitStruct struct {
	RemoteList []string
	BranchList []string

	baseReposDir            string
	repositoryformatversion string
	filemode                string // int
	bare                    string // bool
	logallrefupdates        string // bool
	remote                  []remoteS
	branch                  []branchS
	crypted                 unCrypt
}

type remoteS struct {
	name  string
	url   string
	fetch string
}

type branchS struct {
	name   string
	remote string
	merge  string
}

type unCrypt struct {
	Passphrase string
	PwKey      string
	Token      string
	// local ssh parameters   // Default values
	sshIdRsaPath       string // `~/.ssh/id_rsa`
	serveStorageFile   string // `./assets/cmd/passBridge`
	servePassphraseCmd string // `absoluteRealPath + assets/cmd/cmd`

}

// TODO find a way to indicate how to find pw for pre-decryption
// Actually based on md5 of the caller executable ...
var pathMain = "gitReposMgr"

// Write crypted access to file
func (u *unCrypt) Write(inDev ...bool) (err error) {
	var devMode bool
	if len(inDev) > 0 {
		devMode = inDev[0]
	}
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(&u); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			cp := new(gcr.AES256CipherStruct)
			cp.Data = out.Bytes()
			if err = cp.Encrypt(gcr.Md5File(os.Args[0], devMode)); err == nil {
				err = ioutil.WriteFile(u.serveStorageFile, cp.Data, os.ModePerm)
			}
		}
	}
	return err
}

// Read access from file. After data was readed, the file
// is written 5 times with differents values and destroyed.
// for security reasons.
func (u *unCrypt) Read(filename string) (err error) {
	var file *os.File
	var fi os.FileInfo
	filename = filepath.Join(getAbsRealPath(), filename)
	rewriteValues := []byte{0, 64, 96, 128, 255}
	if file, err = os.OpenFile(filename, os.O_RDWR, os.ModePerm); err == nil {
		defer file.Close()
		if fi, err = file.Stat(); err == nil {
			data := make([]byte, fi.Size())
			if _, err = file.Read(data); err == nil {
				cp := new(gcr.AES256CipherStruct)
				cp.Data = data
				if err = cp.Uncrypt(gcr.Md5File(pathMain)); err == nil {
					if err = json.Unmarshal(cp.Data, &u); err == nil {
						b := new(gcr.Base64)
						u.PwKey = string(b.Decode(u.PwKey))
						cp.Data = b.Decode(u.Passphrase)
						cp.Uncrypt(u.PwKey)
						u.Passphrase = string(cp.Data)
						// Secure delete: Re-Write 5 times differents values to file
						for _, val := range rewriteValues {
							for idx, _ := range data {
								data[idx] = val // Fill with current value {0, 64, 96, 128, 255}
							}
							file.Seek(0, 0)
							if _, err = file.Write(data); err == nil { // Write data to file
								if err = file.Sync(); err != nil { // To be sure the file is fully written
									break
								}
							}
						}
						file.Close()
						os.Remove(filename) // Delete file
					}
				} else {
					fmt.Println(err.Error())
				}
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

// Retrieve current realpath
func getAbsRealPath() (absoluteRealPath string) {
	if absoluteBaseName, err := os.Executable(); err == nil {
		absoluteRealPath = filepath.Dir(absoluteBaseName)
	} else {
		log.Fatal(err)
	}
	return
}

// gitNoSshCommand:
func (gs *gitStruct) CmdNoSSH(gitCmd string, skipErr ...bool) (stdOut string, err error) {
	var noErr string
	if len(skipErr) > 0 {
		if skipErr[0] {
			noErr = " 2> /dev/null"
		}
	}
	var out []byte
	var execCmd *exec.Cmd
	execCmd = exec.Command(`/bin/bash`, `-c`, gitCmd+noErr)
	out, err = execCmd.CombinedOutput()
	stdOut = string(out)

	return
}

func (gs *gitStruct) Cmd(gitCmd string) (stdOut string, err error) {
	// var out bytes.Buffer
	var out []byte
	var execCmd *exec.Cmd
	// write encrypted data to bridge file
	if err = gs.crypted.Write(); err == nil {
		// Build command to get access without passphrase prompt.
		execCmd = exec.Command(`/bin/bash`, `-c`,
			`eval "$(ssh-agent)" && cat `+gs.crypted.sshIdRsaPath+
				` | SSH_ASKPASS="`+gs.crypted.servePassphraseCmd+`" ssh-add - ; `+gitCmd)
		// execCmd.Stdout = &out
		if out, err = execCmd.CombinedOutput(); err == nil {
			stdOut = string(out)
		}
	}
	return
}

func (gs *gitStruct) ReadConfig() (err error) {
	var data []byte
	var inStrSl []string
	var lStart, lEnd int
	var ok bool
	// Terget file: "./.git/config"
	if data, err = ioutil.ReadFile(filepath.Join(gs.baseReposDir, ".git", "config")); err == nil {
		inStrSl = strings.Split(string(data), gst.GetTextEOL(data))
		// [core] part
		ok, lStart, lEnd = getGitCfgSection("core", 0, inStrSl)
		if ok {
			for lIdx := lStart + 1; lIdx <= lEnd; lIdx++ {
				spLine := strings.Split(inStrSl[lIdx], "=")
				switch {
				case strings.TrimSpace(spLine[0]) == "repositoryformatversion":
					gs.repositoryformatversion = spLine[1]
				case strings.TrimSpace(spLine[0]) == "filemode":
					gs.filemode = spLine[1]
				case strings.TrimSpace(spLine[0]) == "bare":
					gs.bare = spLine[1]
				case strings.TrimSpace(spLine[0]) == "logallrefupdates":
					gs.logallrefupdates = spLine[1]
				}
			}
		}
		// [remote ...] part
		ok, lStart, lEnd = getGitCfgSection("remote", lEnd+1, inStrSl)
		for ok {
			r := remoteS{}
			for lIdx := lStart; lIdx <= lEnd; lIdx++ {
				spRemote := strings.Split(inStrSl[lIdx], `"`)
				if len(spRemote) > 1 {
					r.name = spRemote[1]
					continue
				}
				spLine := strings.Split(inStrSl[lIdx], "=")
				switch {
				case strings.TrimSpace(spLine[0]) == "url":
					r.url = spLine[1]
				case strings.TrimSpace(spLine[0]) == "fetch":
					r.fetch = spLine[1]
				}
			}
			gs.RemoteList = append(gs.RemoteList, r.name)
			gs.remote = append(gs.remote, r)
			ok, lStart, lEnd = getGitCfgSection("remote", lEnd+1, inStrSl)
		}
		// [branch ...] part
		ok, lStart, lEnd = getGitCfgSection("branch", lEnd+1, inStrSl)
		for ok {
			b := branchS{}
			for lIdx := lStart; lIdx <= lEnd; lIdx++ {
				spBranch := strings.Split(inStrSl[lIdx], `"`)
				if len(spBranch) > 1 {
					b.name = spBranch[1]
					continue
				}
				spLine := strings.Split(inStrSl[lIdx], "=")
				switch {
				case strings.TrimSpace(spLine[0]) == "remote":
					b.remote = spLine[1]
				case strings.TrimSpace(spLine[0]) == "merge":
					b.merge = spLine[1]
				}
			}
			gs.BranchList = append(gs.BranchList, b.name)
			gs.branch = append(gs.branch, b)
			ok, lStart, lEnd = getGitCfgSection("branch", lEnd+1, inStrSl)
		}
	}
	return
}

func getGitCfgSection(section string, lBegin int, inStrSl []string) (ok bool, lStart, lEnd int) {
	matchSectName := regexp.MustCompile(`^(\[` + section + `.*\])`)
	matchSection := regexp.MustCompile(`^(\[.*\])`)

	for lIdx := lBegin; lIdx < len(inStrSl); lIdx++ {
		line := inStrSl[lIdx]
		if matchSectName.Match([]byte(line)) {
			lStart = lIdx
			for lIdx = lStart + 1; lIdx < len(inStrSl); lIdx++ {
				line := inStrSl[lIdx]
				if matchSection.Match([]byte(line)) || lIdx == len(inStrSl)-1 {
					lEnd = lIdx - 1
					ok = true
					break
				}
			}
		}
	}
	return
}
