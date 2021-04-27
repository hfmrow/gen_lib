// signs.go

/*
	Copyright ©2019 H.F.M - Magic Numbers detection library v0.5 github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This library determine type of files by examining his magic number.
*/

package magicNumbers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	// gltsbh "github.com/hfmrow/gen_lib/tools/bench"
)

func mainTest() {
	// appcache text/cache-manifest
	var err error
	var signs *MagicNumbers

	var files []string
	root := "/home/syndicate/Documents/dev/Web-Design/hfmrow-html/public_html"
	// root = "/media/syndicate/storage/Documents/dev/go/src/github.com/hfmrow/gen_lib/files/magicNumbers/tst/"

	if signs, err = MagicNumbersNew(); err != nil {
		fmt.Println(err)
		return
	}

	// /home/syndicate/Documents/dev/Web-Design/hfmrow-html/public_html/wp-includes/certificates/ca-bundle.crt
	// /home/syndicate/Documents/dev/Web-Design/hfmrow-html/public_html/wp-content/uploads/wp-slimstat/maxmind.mmdb qt, mov video/quicktime

	// if err = signs.AddMagicNumber("class", "application/java-vm", true); err == nil {
	// 	if err = signs.AddMagicNumber("json", "application/json", true); err == nil {
	// 		err = signs.AddMagicNumber("js", "application/javascript", true)
	// 	}
	// }

	// mns := MagicNumbersSignatureNew(0, "774F464600010000")
	// mns1 := MagicNumbersSignatureNew(0, "774F46464F54544F")
	// err = signs.AddMagicNumber("woff", "application/font-woff", false, mns, mns1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// mns = MagicNumbersSignatureNew(0, "774F463200010000")
	// err = signs.AddMagicNumber("woff2", "application/font-woff2", false, mns)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// mns = MagicNumbersSignatureNew(22, "4C500000000000000000000000000000000001")
	// err = signs.AddMagicNumber("eot", "application/vnd.ms-fontobject", false, mns)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = signs.Rebuild()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// mns.Range = 512
	// mns.RangeEnd = 7
	// if err = signs.AddMagicNumber("svg", "image/svg+xml", true, mns); err != nil {
	// 	// if err = signs.AddMagicNumber("debug", "application/golang-library-debug-file",false,
	// 	// 	MagicNumbersNewSig(0, "213C617263683E0A5F5F2E504B47444546")); err != nil {
	// 	fmt.Println(err)
	// 	return
	// 	// }
	// }

	if err != nil {
		fmt.Println(err)
	}

	// if err = signs.Rebuild(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}
		return err
	})

	// bench := gltsbh.BenchNew()
	// bench.Lapse()

	for _, file := range files {
		// file := "/media/syndicate/storage/Documents/dev/go/src/github.com/hfmrow/gen_lib/files/magicNumbers/tst/content.inline.min.css"
		out, err := signs.CheckFile(file)
		if err != nil {
			fmt.Println(err)
		}
		out = out
		fmt.Println((*out)[0].Ext, (*out)[0].Mime, file)
	}

	// bench.Stop()

	if err != nil {
		fmt.Println(err)
	}
}

// MagicNumbersNew:
func MagicNumbersNew() (out *MagicNumbers, err error) {
	out = new(MagicNumbers)
	out.DispLen = 32
	out.structName = "signStruct.json"
	out.structAdded = "signAdded.json"
	// https://gist.github.com/qti3e/6341245314bf3513abb080677cd1c93b
	// the signature file come from the url above.
	out.importSigns = "signatures.json"

	// If the main structure has been removed, it will be automatically rebuild.
	if err = out.Read(out.structName); os.IsNotExist(err) {
		err = out.Rebuild()
	}
	return
}

// TODO Thing to add multi-patterns matching in a specific range ...
// AddMagicNumber: Use MagicNumbersNewSig() to create a new signature before adding.
// i.e: notice, error not hendled here ... Rebuild() must be used to finalize operation !
//	1- Add 2 signatures for .woff type "application/font-woff", not a txt file
//		mns := MagicNumbersSignatureNew(0, "774F464600010000")
//		mns1 := MagicNumbersSignatureNew(0, "774F46464F54544F")
//		signs.AddMagicNumber("woff", "application/font-woff", false, mns, mns1)
//
//	2- Add 1 signature start at offset 22 for .eot type "application/vnd.ms-fontobject",
//	not a txt file
//		mns = MagicNumbersSignatureNew(22, "4C500000000000000000000000000000000001")
//		signs.AddMagicNumber("eot", "application/vnd.ms-fontobject", false, mns)
//
//	3- Add only extension .js and his mime type "application/javascript", is a txt file
//		signs.AddMagicNumber("js", "application/javascript", true)
//
//	4- Add 1 signature start at offset 0, to offset 512 (pattern if searched in this range),
//	Another search is done from the end of the file with a range of 7 to match "3C2F7376673E"
//	(</svg>). This signature is for the .svg type "image/svg+xml", is a txt file.
//		mns = MagicNumbersSignatureNew(0, "3C73766720786D6C6E73", "3C2F7376673E")
//		mns.Range=512  // Range to search from start offset
//		mns.RangeEnd=7 // Range from end, by default that contain length +1 of "ValueEnd"
//		signs.AddMagicNumber("svg", "image/svg+xml", true, mns)
//		signs.Rebuild()
func (m *MagicNumbers) AddMagicNumber(ext, mime string, isText bool, signaturesList ...MagicNumbersSignature) (err error) {
	tmpMn := new(MagicNumbers)
	var signs []MagicNumbersSignature
	if err = tmpMn.Read(m.structAdded); err != nil {
		return
	}
	if len(signaturesList) > 0 {
		signs = make([]MagicNumbersSignature, len(signaturesList))
		for idx, mns := range signaturesList {
			if b, err := hex.DecodeString(mns.Value); err == nil {
				signs[idx] = mns
				signs[idx].Bytes = b
			} else {
				return err
			}
		}
	} else {
		// Add to "MainStructure" without signatures only based on extension/mime type
		signs = []MagicNumbersSignature{}
	}
	// Add to "MainStructure"
	m.appendCheckExist(signatures{
		Ext:    ext,
		Signs:  signs,
		Mime:   mime,
		IsText: isText})
	m.compute()
	// Add to "AddedStructure"
	tmpMn.appendCheckExist(signatures{
		Ext:    ext,
		Signs:  signs,
		Mime:   mime,
		IsText: isText})
	// Write new record to "AddedStructure"
	return tmpMn.Write(m.structAdded)
}

// TODO Add (remove) and (modify) method ...
// appendCheckExist: Add or append whether already exist or not.
func (m *MagicNumbers) appendCheckExist(inSigns signatures) {
	var exist, extExist bool
	var inSig MagicNumbersSignature
	var inIdx int

	for idx, sigs := range m.List {
		if sigs.Ext == inSigns.Ext {
			if sigs.Mime == inSigns.Mime { // TODO the case where different mime/ext have the same signature ...
				extExist = true // We have same extension and mime type
				for _, sig := range sigs.Signs {
					exist = false
					for inIdx, inSig = range inSigns.Signs {
						if sig.Value == inSig.Value && sig.Offset == inSig.Offset {
							exist = true
							break
						}
					}
					if !exist { // Signature does not already exist, then add it
						m.List[idx].Signs = append(m.List[idx].Signs, inSigns.Signs[inIdx])
					}
				}
			}
		}
		if extExist {
			break
		}
	}
	if !extExist {
		m.List = append(m.List, inSigns)
	}
}

// Rebuild: using new added signature and the original json signatures files.
func (m *MagicNumbers) Rebuild() (err error) {
	if err = m.Read(m.structAdded); err == nil {
		if err = m.ImportFromFile(m.importSigns); err == nil {
			m.DispLen = 32
			err = m.Write(m.structName)
		}
	}
	return
}

// CheckFile:
// TODO find a way to seek only needed bytes to be analysed
func (m *MagicNumbers) CheckFile(filename string) (out *[]outInfos, err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		return m.CheckBytes(textFileBytes, filename)
	}
	return
}

// CheckBytes:
func (m *MagicNumbers) CheckBytes(data []byte, basename ...string) (outInf *[]outInfos, err error) {
	var mNumbers []byte
	var ofstEnd, startEnd int
	var ext, baseExt string
	var tmpOut outInfos
	var out []outInfos

	if len(basename) > 0 {
		baseExt = strings.Trim(filepath.Ext(basename[0]), ".")
	}

	lenData := len(data)
	for _, signs := range m.sorted {
		tmpOut = outInfos{}
		// ext matching
		if len(baseExt) > 0 {
			for _, ext = range signs.Ext {
				if ext == baseExt {
					tmpOut.score = tmpOut.score + 0.1
					break
				}
			}
		}
		// From start
		if signs.Signature.Bytes != nil && signs.Signature.Offset < lenData {
			ofstEnd = signs.Signature.Offset + signs.Signature.Range
			if ofstEnd > lenData {
				ofstEnd = lenData
			}
			mNumbers = data[signs.Signature.Offset:ofstEnd]
			if bytes.Contains(mNumbers, signs.Signature.Bytes) {
				tmpOut.score = tmpOut.score + 0.1
			}
		}
		// From end
		if signs.Signature.BytesEnd != nil {
			if startEnd = (lenData - 1) - signs.Signature.RangeEnd; startEnd > 0 {
				if bytes.Contains(data[startEnd:lenData-1], signs.Signature.BytesEnd) {
					tmpOut.score = tmpOut.score + 0.1
				}
			}
		}
		// Have we found something ?
		if tmpOut.score > 0 {
			tmpOut.Ext = strings.Join(signs.Ext, ", ")
			tmpOut.Mime = signs.Mime
			tmpOut.IsText = signs.IsText
			out = append(out, tmpOut)
		}
	}
	if len(out) > 0 {
		sort.SliceStable(out, func(i, j int) bool { // sorting
			return out[i].score > out[j].score
		})
	} else {
		// Where there is no associated patterns with this type of file
		// we get only "Displen" first byte and put them as hexadecimal
		// representation and ascii in the returned information.
		tmpOut.Ext = fmt.Sprintf("Unknown: [%"+fmt.Sprintf("%d", m.DispLen)+"X]", data[0:m.DispLen])
		tmpOut.Mime = fmt.Sprintf("[%s]", strings.Join(strings.Split(strings.Join(strings.Split(string(data[0:m.DispLen]), "\r"), ""), "\n"), ""))
		out = append(out, tmpOut)
	}
	return &out, err
}

// ImportFromFile: Import base data containing signatures information.
func (m *MagicNumbers) ImportFromFile(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = m.importFromBytes(textFileBytes); err == nil {
			m.compute()
		}
	}
	return
}

// importFromBytes: see above
func (m *MagicNumbers) importFromBytes(textFileBytes []byte) (err error) {
	var lines []string
	var b []byte
	var s []MagicNumbersSignature
	var extName string

	extNameRe := regexp.MustCompile(`("\w*": \{)`)
	signSttRe := regexp.MustCompile(`("signs": \[)`)
	signDefRe := regexp.MustCompile(`("\d*,\w*")`)
	mimeDefRe := regexp.MustCompile(`("mime": ".*")`)

	lines = strings.Split(string(textFileBytes), m.getTextEOL(&textFileBytes))
	for idx := 0; idx < len(lines); idx++ {
		switch {
		case extNameRe.MatchString(lines[idx]):
			extName = strings.Split(lines[idx], `"`)[1]
		case signSttRe.MatchString(lines[idx]):
			idx++
			s = []MagicNumbersSignature{}
			for signDefRe.MatchString(lines[idx]) {
				splitted := strings.Split(strings.Split(lines[idx], `"`)[1], `,`)
				if ofst, err := strconv.Atoi(splitted[0]); err == nil {
					if b, err = hex.DecodeString(splitted[1]); err == nil {
						s = append(s, MagicNumbersSignature{
							Offset:   ofst,
							Range:    len(b),
							RangeEnd: -1,
							Value:    splitted[1],
							Bytes:    b,
						})
						idx++
						continue
					}
				}
				return fmt.Errorf("Unable to convert line: %d, got error: %s\n", idx, err.Error())
			}
		case mimeDefRe.MatchString(lines[idx]):
			m.List = append(m.List, signatures{
				Ext:   extName,
				Mime:  strings.Split(lines[idx], `"`)[3],
				Signs: s})
		}
	}
	return
}

// Write: main structure containing whole signatures
func (m *MagicNumbers) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(m); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			if err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm); err == nil {
				m.compute()
			}
		}
	}
	return
}

// Read: main structure containing whole signatures
func (m *MagicNumbers) Read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = m.ReadFromBytes(textFileBytes)
	}
	return
}

// ReadFromBytes:
func (m *MagicNumbers) ReadFromBytes(textFileBytes []byte) (err error) {
	if err = json.Unmarshal(textFileBytes, m); err == nil {
		m.compute()
	}
	return
}

// compute: build temporary slice that contains whole information on signatures.
// this slice is sorted in a specifics ways that permit to handle correctly the
// searching process. In fact, bigger pattern are placed in first place just after
// signatures that contains multiples search process. The other thing of this step
// is to prepare in an efficient way tha data to reduce processing time.
// This process is only executed one time when a new structure is created.
func (m *MagicNumbers) compute() {
	var err error
	var tmpSorted []sortedSignatures

	for _, sig := range m.List {
		if sig.Mime == "application/octet-stream" {
			fmt.Println("application/octet-stream")
		}
		var exts []string
		// Get extensions associated with this mime type. Not really useful but ...
		if exts, err = mime.ExtensionsByType(sig.Mime); err != nil || len(exts) == 0 {
			exts = append(exts, sig.Ext)
		}
		for idx, e := range exts {
			exts[idx] = strings.Trim(e, ".")
		}

		if len(sig.Signs) == 0 {
			m.sorted = append(m.sorted, sortedSignatures{
				Ext:    exts,
				Mime:   sig.Mime,
				IsText: sig.IsText})
		} else {
			for _, s := range sig.Signs {
				m.sorted = append(m.sorted, sortedSignatures{
					Ext:       exts,
					Mime:      sig.Mime,
					Signature: s,
					IsText:    sig.IsText})
			}
		}
	}
	sort.SliceStable(m.sorted, func(i, j int) bool { // AlNum sorting, pre sorting to get a consistent second pass.
		return strings.ToLower(m.sorted[i].Signature.Value) > strings.ToLower(m.sorted[j].Signature.Value)
	})
	sort.SliceStable(m.sorted, func(i, j int) bool { // Len sorting, so, all patterns sorted from bigger to lower size
		return len(m.sorted[i].Signature.Bytes) > len(m.sorted[j].Signature.Bytes)
	})
	// Pull all signature with "End" data to the start of main slice.
	for idx, entry := range m.sorted {
		if entry.Signature.BytesEnd != nil && entry.Signature.RangeEnd > 0 {
			tmpSorted = append([]sortedSignatures{entry}, tmpSorted...) // Prepend to temp slice
			m.sorted = append(m.sorted[:idx], m.sorted[idx+1:]...)      // Remove from base slice
		}
	}
	m.sorted = append(tmpSorted, m.sorted...) // Prepend temp slice to base slice
}

// MagicNumbers:
type MagicNumbers struct {
	// Will be a saved file (read/write).
	List []signatures
	// whether there is no signature that match, output display "Displen"
	// byte from the start of the input data.
	DispLen int
	// internal usage, will be computed at runtime, this structure contain
	// both version, added one and base structure imported from json file
	// this structure is sorted in descending order by the size of the
	// signatures values.
	sorted      []sortedSignatures
	structName  string // This is the structure version of "importSigns" (json)
	structAdded string // This is the user added signatures
	importSigns string // This is original file formatted as json
}

// MagicNumbersSignature:
type MagicNumbersSignature struct {
	Offset,
	Range,
	RangeEnd int
	Value,
	ValueEnd string
	Bytes,
	BytesEnd []byte
}

// MagicNumbersSignatureNew:
func MagicNumbersSignatureNew(offset int, valueStart string, valueEnd ...string) MagicNumbersSignature {
	var bytesStart, bytesEnd []byte
	var err error
	var valueEndOut string

	if bytesStart, err = hex.DecodeString(valueStart); err != nil {
		return MagicNumbersSignature{}
	}
	if len(valueEnd) > 0 {
		valueEndOut = valueEnd[0]
		if bytesEnd, err = hex.DecodeString(valueEndOut); err != nil {
			return MagicNumbersSignature{}
		}
	}

	return MagicNumbersSignature{
		Offset:   offset,
		Range:    len(bytesStart),
		Value:    valueStart,
		Bytes:    bytesStart,
		RangeEnd: len(bytesEnd) + 1,
		ValueEnd: valueEndOut,
		BytesEnd: bytesEnd}
}

type outInfos struct {
	Ext, Mime string
	IsText    bool
	score     float64
}

type sortedSignatures struct {
	Ext       []string
	Mime      string
	Signature MagicNumbersSignature
	IsText    bool
}

type signatures struct {
	Ext    string
	Signs  []MagicNumbersSignature
	Mime   string
	IsText bool
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF)
func (m *MagicNumbers) getTextEOL(inTextBytes *[]byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	switch {
	case bytes.Contains((*inTextBytes), bCRLF):
		outString = string(bCRLF)
	case bytes.Contains((*inTextBytes), bCR):
		outString = string(bCR)
	default:
		outString = string(bLF)
	}
	return
}
