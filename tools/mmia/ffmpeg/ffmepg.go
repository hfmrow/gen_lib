// ffmepg.go

/*
	Copyright ©2021 hfmrow - ffmpeg go wrapper v1.0 github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	glts "github.com/hfmrow/gen_lib/tools"
)

var (
	// Lib mapping
	ExecCommand         = glts.ExecCommand
	ExecCommandProgress = glts.ExecCommandProgress
)

// FFmpegMetadataNew create a new structure to hold media information
func FFmpegMetadataNew(filename string) (*FFmpegMetadata, error) {
	var err error
	var out string
	ffmpg := new(FFmpegMetadata)
	ffmpg.Filename = filename
	cmd := []string{
		"ffprobe",
		"-i",
		ffmpg.Filename,
		"-print_format",
		"json",
		"-show_format",
		"-show_streams",
		"-show_error"}
	if out, err = ExecCommand(cmd); err == nil {
		err = json.Unmarshal([]byte(out), ffmpg)
	}
	if err != nil {
		return nil, err
	}
	return ffmpg, nil
}

// FFmpegCutMedia if 'cutAtEnd' = true, 'timeSec' will be deleted
// at the start, otherwise the time will be deleted at the end.
// A 'progress' function is optionally available to retrieve realtime
// information (shell stdout) while current operation.
func FFmpegCutMedia(in, out string, timeSec int, cutAtEnd, overwrite bool,
	progress ...func(items map[string]string)) error {
	var (
		err       error
		direction = "-t"
		ovrWrt    = "-n"
		// Define default progress function
		callbackProgress = func(items map[string]string) {
			for n, v := range items {
				fmt.Println(n, v)
			}
			// fmt.Printf("\n---start---\n%s\n---end---\n\n", strings.Join(lines, "\n"))
			time.Sleep(time.Second)
		}
	)
	if len(progress) > 0 {
		callbackProgress = progress[0]
	}
	if cutAtEnd {
		direction = "-ss"
	}
	if overwrite {
		ovrWrt = "-y"
	}
	// build command
	cmd := []string{
		"ffmpeg",
		"-hide_banner",
		"-loglevel",
		"error",
		"-progress",
		"-",
		"-nostats",
		ovrWrt,
		"-i",
		in,
		direction,
		// Format 'timeSec' to 'hh:mm:ss'
		time.Date(2006, 01, 02, 0, 0, timeSec, 0, time.Local).Format("15:04:05"),
		"-vcodec",
		"copy",
		"-acodec",
		"copy",
		out}

	err = ExecCommandProgress(cmd, 12, func(lines []string) {
		items := make(map[string]string)
		for _, line := range lines {
			splitted := strings.Split(line, "=")
			if len(splitted) > 1 {
				items[splitted[0]] = splitted[1]
			}
		}
		callbackProgress(items)
	})
	if err != nil {
		if strings.Contains(err.Error(), "No such file or directory") {
			err = os.ErrNotExist
		}
		if strings.Contains(err.Error(), "already exists. Exiting.") {
			err = os.ErrExist
		}
	}
	return err
}

// FFmpegMetadata structure to hold media information provided by ffprobe.
type FFmpegMetadata struct {
	Filename string
	Streams  []Stream `json:"streams"`
	Format   format   `json:"format"`
}
type Stream struct {
	ID                 string      `json:"id"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile"`
	CodecType          string      `json:"codec_type"`
	CodecTimeBase      string      `json:"codec_time_base"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio"`
	DisplayAspectRatio string      `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ChromaLocation     string      `json:"chroma_location"`
	Refs               int         `json:"refs"`
	QuarterSample      string      `json:"quarter_sample"`
	DivxPacked         string      `json:"divx_packed"`
	RFrameRrate        string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	TimeBase           string      `json:"time_base"`
	DurationTs         int         `json:"duration_ts"`
	Duration           string      `json:"duration"`
	Disposition        disposition `json:"disposition"`
	BitRate            string      `json:"bit_rate"`
	Tags               streamTags  `json:"tags"`
}
type streamTags struct {
	Title          string `json:"title"`
	Language       string `json:"language"`
	Bps            string `json:"BPS-eng"`
	Duration       string `json:"DURATION-eng"`
	FramesCount    string `json:"NUMBER_OF_FRAMES-eng"`
	BytesCount     string `json:"NUMBER_OF_BYTES-eng"`
	WritingApp     string `json:"_STATISTICS_WRITING_APP-eng"`
	WritingDateUTC string `json:"_STATISTICS_WRITING_DATE_UTC-eng"`
}
type disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
}
type format struct {
	NbStreams      int        `json:"nb_streams"`
	NbPrograms     int        `json:"nb_programs"`
	FormatName     string     `json:"format_name"`
	FormatLongName string     `json:"format_long_name"`
	Duration       string     `json:"duration"`
	Size           string     `json:"size"`
	BitRate        string     `json:"bit_rate"`
	ProbeScore     int        `json:"probe_score"`
	Tags           formatTags `json:"tags"`
}
type formatTags struct {
	Encoder      string `json:"encoder"`
	CreationTime string `json:"creation_time"`
}
