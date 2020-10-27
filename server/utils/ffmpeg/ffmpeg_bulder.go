package ffmpeg

import (
	"fmt"
	"strconv"
)

var resolutions = map[string]string{
	"720p": "",
}

type FFmpegCommand struct {
	Timelimit        string
	StartTime        string
	SourceFile       string
	SegementDuration string
	Resolution       string
	VideoCodec       string
	Preset           string
	AudioCodec       string
	PixFmt           string
	ForceKeyFrames   string
	ForceFormat      string
	SegementTime     string
	InitialOffset    string
	Pipe             string
}

func (ff *FFmpegCommand) SetTimelimit(limit int) *FFmpegCommand {
	ff.Timelimit = strconv.Itoa(limit)

	return ff
}

func (ff *FFmpegCommand) SetStartTime(startTime int) *FFmpegCommand {
	ff.StartTime = fmt.Sprintf("%v.00", startTime)

	return ff
}

func (ff *FFmpegCommand) SetSourceFile(path string) *FFmpegCommand {
	ff.SourceFile = path

	return ff
}

func (ff *FFmpegCommand) SetSegmentDuration(duration int) *FFmpegCommand {
	ff.SegementDuration = fmt.Sprintf("%v.00", duration)

	return ff
}

func (ff *FFmpegCommand) SetResolution(resolution string) *FFmpegCommand {
	res, ok := resolutions[resolution]
	if !ok {
		panic("Wrong resolution")
	}
	ff.Resolution = res

	return ff
}

func (ff *FFmpegCommand) SetVideoCodec(codec string) *FFmpegCommand {
	ff.VideoCodec = codec

	return ff
}

func (ff *FFmpegCommand) SetPreset(preset string) *FFmpegCommand {
	ff.Preset = preset

	return ff
}

func (ff *FFmpegCommand) SetAudioCodec(codec string) *FFmpegCommand {
	ff.AudioCodec = codec

	return ff
}

func (ff *FFmpegCommand) SetPixFmt(fmt string) *FFmpegCommand {
	ff.PixFmt = fmt

	return ff
}

func (ff *FFmpegCommand) EnableForceKeyFrames() *FFmpegCommand {
	ff.ForceKeyFrames = "expr:gte(t, n_forced*5.000)"

	return ff
}

func (ff *FFmpegCommand) SetForceFormat(format string) *FFmpegCommand {
	ff.ForceFormat = format

	return ff
}

func (ff *FFmpegCommand) SetSegmentTime(duration int) *FFmpegCommand {
	ff.SegementTime = fmt.Sprintf("%v.00", duration)

	return ff
}

func (ff *FFmpegCommand) SetInitialOffest(offset int) *FFmpegCommand {
	ff.InitialOffset = fmt.Sprintf("%v.00", offset)

	return ff
}

func (ff *FFmpegCommand) Build() string {
	stmt := "ffmpeg"
	if ff.Timelimit != "" {
		stmt += " " + fmt.Sprintf("-timelimit %s", ff.Timelimit)
	}
	if ff.StartTime != "" {
		stmt += " " + fmt.Sprintf("-ss %s", ff.StartTime)
	}
	if ff.SourceFile != "" {
		stmt += " " + fmt.Sprintf("-i %s", ff.SourceFile)
	}
	if ff.SegementDuration != "" {
		stmt += " " + fmt.Sprintf("-t %s", ff.SegementDuration)
	}
	if ff.Resolution != "" {
		stmt += " " + fmt.Sprintf("-vf %s", ff.Resolution)
	}
	if ff.VideoCodec != "" {
		stmt += " " + fmt.Sprintf("-vcodec %s", ff.VideoCodec)
	}
	if ff.Preset != "" {
		stmt += " " + fmt.Sprintf("-preset %s", ff.Preset)
	}
	if ff.AudioCodec != "" {
		stmt += " " + fmt.Sprintf("-acodec %s", ff.AudioCodec)
	}
	if ff.PixFmt != "" {
		stmt += " " + fmt.Sprintf("-pix_fmt %s", ff.PixFmt)
	}
	if ff.ForceKeyFrames != "" {
		stmt += " " + fmt.Sprintf("-force_key_frames %s", ff.ForceKeyFrames)
	}
	if ff.ForceFormat != "" {
		stmt += " " + fmt.Sprintf("-f %s", ff.ForceFormat)
	}
	if ff.SegementTime != "" {
		stmt += " " + fmt.Sprintf("-segment_time %s", ff.SegementTime)
	}
	if ff.InitialOffset != "" {
		stmt += " " + fmt.Sprintf("-initial_offset %s", ff.InitialOffset)
	}
	stmt += " " + ff.Pipe

	return stmt
}

func New() *FFmpegCommand {
	return &FFmpegCommand{
		Pipe: "pipe:out%03d.ts",
	}
}
