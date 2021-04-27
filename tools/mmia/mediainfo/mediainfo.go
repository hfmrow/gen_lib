// mediainfo.go

/*
	Use: go-mediainfo
	github.com/zelenin/go-mediainfo
	Author: Aleksandr Zelenin, e-mail: aleksandr@zelenin.me
	Unspecified license type.

	Copyright ©2021 hfmrow - mediainfo wrapper library github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package mediainfo

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	glts "github.com/hfmrow/gen_lib/tools"
)

var (
	// Lib mapping
	ExecCommand         = glts.ExecCommand
	ExecCommandProgress = glts.ExecCommandProgress
)

// MediaInfoStructNew create a new MediaInfo structure
func MediaInfoStructNew(filename string) (*MediaInfoStruct, error) {
	var (
		out string
		err error
	)
	mis := new(MediaInfoStruct)
	cmd := []string{
		"mediainfo",
		"--Output=JSON",
		"--full",
		filename,
	}

	var m mediaInfoStruct
	if out, err = ExecCommand(cmd); err == nil {
		err = json.Unmarshal([]byte(out), &m)
	}
	if err != nil {
		return nil, err
	}
	if len(m.Media.Streams) > 0 {
		m.Media.AudioCount = countTypeStream(&m, "Audio")
		m.Media.VideoCount = countTypeStream(&m, "Video")
		m.Media.TextCount = countTypeStream(&m, "Text")
		m.Media.MenuCount = countTypeStream(&m, "Menu", true)
		if m.Media.MenuCount > 0 {

			// Retrieve Menu with mediainfo
			cmd = []string{
				"mediainfo",
				"mediainfo --full",
				filename,
			}
			if out, err = ExecCommand(cmd); err != nil {
				return nil, err
			}
			lines := strings.Split(out, "\n")
			for idx := len(lines) - 1; idx >= 0; idx-- {
				if strings.TrimSpace(lines[idx]) == "Menu" {
					lines = lines[idx:]
					break
				}
			}
			m.Media.Chapters = wrapMenu(lines)
		}
		mis.Media = append(mis.Media, m.Media)
	}
	return mis, nil

}

type menu struct {
	Entries []menuEntry
}

type menuEntry struct {
	Time,
	Text,
	Lang string
}

func wrapMenu(lines []string) menu {
	regSpace := regexp.MustCompile(`\s+`)
	regEntry := regexp.MustCompile(`(?m)\d{2}:\d{2}:\d{2}\.\d{3}`)
	var m menu
	for _, line := range lines {
		if regEntry.MatchString(line) {
			items := regSpace.Split(line, -1)
			time := strings.TrimSpace(items[0])
			items = strings.Split(strings.Join(items[1:], " "), ":")
			m.Entries = append(m.Entries, menuEntry{
				Time: time,
				Text: strings.TrimSpace(strings.Join(items[2:], ":")),
				Lang: strings.TrimSpace(items[1]),
			})
		}
	}
	return m
}

/*
 * Helper functions
 */
func countTypeStream(m *mediaInfoStruct, typeStream string, removeStream ...bool) int {
	count := 0
	rs := false
	if len(removeStream) > 0 {
		rs = true
	}
	for idx, stream := range m.Media.Streams {
		if stream.Type == typeStream {
			count++
			if rs {
				// Remove stream
				m.Media.Streams = append(m.Media.Streams[:idx], m.Media.Streams[idx+1:]...)
			} else {
				m.Media.Streams[idx].InStreamID = count
			}
		}
	}
	return count
}

func toInt(in string) int {
	if len(in) > 0 {
		val, _ := strconv.ParseFloat(in, 64)
		return int(val)
	}
	return -1
}

// MediaInfoStruct that hold MediaInfo data provided by 'mediainfo'
// shell command.
type MediaInfoStruct struct {
	Media []media `json:"media"`
}

// Used for temporary storage
type mediaInfoStruct struct {
	Media media `json:"media"`
}
type media struct {
	AudioCount,
	MenuCount,
	TextCount,
	VideoCount int
	Filename string  `json:"@ref"`
	Streams  []Track `json:"track"`
	Chapters menu
}
type Track struct {
	Type                           string `json:"@type"`
	UniqueID                       string `json:"UniqueID"`
	videoCount                     string `json:"VideoCount,omitempty"`
	audioCount                     string `json:"AudioCount,omitempty"`
	textCount                      string `json:"TextCount,omitempty"`
	FileExtension                  string `json:"FileExtension,omitempty"`
	Format                         string `json:"Format"`
	FormatVersion                  string `json:"Format_Version,omitempty"`
	FileSize                       string `json:"FileSize,omitempty"`
	Duration                       string `json:"Duration,omitempty"`
	OverallBitRate                 string `json:"OverallBitRate,omitempty"`
	FrameRate                      string `json:"FrameRate,omitempty"`
	FrameCount                     string `json:"FrameCount,omitempty"`
	StreamSize                     string `json:"StreamSize,omitempty"`
	IsStreamable                   string `json:"IsStreamable,omitempty"`
	EncodedDate                    string `json:"Encoded_Date,omitempty"`
	FileModifiedDate               string `json:"File_Modified_Date,omitempty"`
	FileModifiedDateLocal          string `json:"File_Modified_Date_Local,omitempty"`
	EncodedApplication             string `json:"Encoded_Application,omitempty"`
	EncodedLibrary                 string `json:"Encoded_Library,omitempty"`
	StreamOrder                    string `json:"StreamOrder,omitempty"`
	ID                             string `json:"ID,omitempty"`
	InStreamID                     int
	FormatProfile                  string `json:"Format_Profile,omitempty"`
	FormatLevel                    string `json:"Format_Level,omitempty"`
	FormatSettingsCABAC            string `json:"Format_Settings_CABAC,omitempty"`
	FormatSettingsRefFrames        string `json:"Format_Settings_RefFrames,omitempty"`
	CodecID                        string `json:"CodecID,omitempty"`
	BitRate                        string `json:"BitRate,omitempty"`
	Width                          string `json:"Width,omitempty"`
	Height                         string `json:"Height,omitempty"`
	StoredHeight                   string `json:"Stored_Height,omitempty"`
	SampledWidth                   string `json:"Sampled_Width,omitempty"`
	SampledHeight                  string `json:"Sampled_Height,omitempty"`
	PixelAspectRatio               string `json:"PixelAspectRatio,omitempty"`
	DisplayAspectRatio             string `json:"DisplayAspectRatio,omitempty"`
	FrameRateMode                  string `json:"FrameRate_Mode,omitempty"`
	FrameRateModeOriginal          string `json:"FrameRate_Mode_Original,omitempty"`
	ColorSpace                     string `json:"ColorSpace,omitempty"`
	ChromaSubsampling              string `json:"ChromaSubsampling,omitempty"`
	BitDepth                       string `json:"BitDepth,omitempty"`
	ScanType                       string `json:"ScanType,omitempty"`
	Delay                          string `json:"Delay,omitempty"`
	Default                        string `json:"Default,omitempty"`
	Forced                         string `json:"Forced,omitempty"`
	Typeorder                      string `json:"@typeorder,omitempty"`
	FormatCommercialIfAny          string `json:"Format_Commercial_IfAny,omitempty"`
	FormatSettingsEndianness       string `json:"Format_Settings_Endianness,omitempty"`
	FormatAdditionalFeatures       string `json:"Format_AdditionalFeatures,omitempty"`
	BitRateMode                    string `json:"BitRate_Mode,omitempty"`
	Channels                       string `json:"Channels,omitempty"`
	ChannelPositions               string `json:"ChannelPositions,omitempty"`
	ChannelLayout                  string `json:"ChannelLayout,omitempty"`
	SamplesPerFrame                string `json:"SamplesPerFrame,omitempty"`
	SamplingRate                   string `json:"SamplingRate,omitempty"`
	SamplingCount                  string `json:"SamplingCount,omitempty"`
	CompressionMode                string `json:"Compression_Mode,omitempty"`
	DelaySource                    string `json:"Delay_Source,omitempty"`
	StreamSizeProportion           string `json:"StreamSize_Proportion,omitempty"`
	Language                       string `json:"Language,omitempty"`
	ServiceKind                    string `json:"ServiceKind,omitempty"`
	ElementCount                   string `json:"ElementCount,omitempty"`
	Title                          string `json:"Title,omitempty"`
	Movie                          string `json:"Movie,omitempty"`
	FormatTier                     string `json:"Format_Tier,omitempty"`
	EncodedLibraryName             string `json:"Encoded_Library_Name,omitempty"`
	EncodedLibraryVersion          string `json:"Encoded_Library_Version,omitempty"`
	EncodedLibrarySettings         string `json:"Encoded_Library_Settings,omitempty"`
	Extra                          extra  `json:"extra,omitempty"`
	ColourDescriptionPresent       string `json:"colour_description_present,omitempty"`
	ColourDescriptionPresentSource string `json:"colour_description_present_Source,omitempty"`
	ColourRange                    string `json:"colour_range,omitempty"`
	ColourRangeSource              string `json:"colour_range_Source,omitempty"`
	ColourPrimaries                string `json:"colour_primaries,omitempty"`
	ColourPrimariesSource          string `json:"colour_primaries_Source,omitempty"`
	TransferCharacteristics        string `json:"transfer_characteristics,omitempty"`
	TransferCharacteristicsSource  string `json:"transfer_characteristics_Source,omitempty"`
	MatrixCoefficients             string `json:"matrix_coefficients,omitempty"`
	MatrixCoefficientsSource       string `json:"matrix_coefficients_Source,omitempty"`
	FormatSettingsSBR              string `json:"Format_Settings_SBR,omitempty"`
	BitRateNominal                 string `json:"BitRate_Nominal,omitempty"`
	FormatSettingsGOP              string `json:"Format_Settings_GOP,omitempty"`
	DisplayAspectRatioOriginal     string `json:"DisplayAspectRatio_Original,omitempty"`
	menuCount                      string `json:"MenuCount,omitempty"`
	PixelAspectRatioOriginal       string `json:"PixelAspectRatio_Original,omitempty"`
}
type extra struct {
	ComplexityIndex         string `json:"ComplexityIndex"`
	NumberOfDynamicObjects  string `json:"NumberOfDynamicObjects"`
	BedChannelCount         string `json:"BedChannelCount"`
	BedChannelConfiguration string `json:"BedChannelConfiguration"`
	Bsid                    string `json:"bsid"`
	Dialnorm                string `json:"dialnorm"`
	Compr                   string `json:"compr"`
	Acmod                   string `json:"acmod"`
	Lfeon                   string `json:"lfeon"`
	DialnormAverage         string `json:"dialnorm_Average"`
	DialnormMinimum         string `json:"dialnorm_Minimum"`
	ComprAverage            string `json:"compr_Average"`
	ComprMinimum            string `json:"compr_Minimum"`
	ComprMaximum            string `json:"compr_Maximum"`
	ComprCount              string `json:"compr_Count"`
	DynrngAverage           string `json:"dynrng_Average"`
	DynrngMinimum           string `json:"dynrng_Minimum"`
	DynrngMaximum           string `json:"dynrng_Maximum"`
	DynrngCount             string `json:"dynrng_Count"`
	DATE                    string `json:"DATE"`
	ErrorDetectionType      string `json:"ErrorDetectionType"`
}
