package ffmpeg

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"recorder-service/internal/domain/recorder"
	"sync"
	"time"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func MixAudio(infos []recorder.RecordingInfo) (outputPath string, err error) {
	if len(infos) < 1 {
		return "", fmt.Errorf("no audio to mix")
	}

	streams := []*ffmpeg_go.Stream{}
	for _, info := range infos {
		streams = append(streams, ffmpeg_go.Input(info.FilePath()))
	}

	amixNode := ffmpeg_go.FilterMultiOutput(streams, "amix", nil, ffmpeg_go.KwArgs{
		"inputs": len(streams),
	})

	outputPath = mixPath(infos[0].BasePath, infos[0].Ext)
	err = amixNode.Stream("", "").Output(outputPath, ffmpeg_go.KwArgs{"y": ""}).Run()
	if err != nil {
		log.Printf("Failed to mix audio: %v", err)
		return "", fmt.Errorf("failed to mix audio")
	}

	return outputPath, nil
}

func AddSilence(infos []recorder.RecordingInfo) error {
	min, max := findRange(infos)

	var wg sync.WaitGroup
	errors := make(chan error, len(infos))

	for _, info := range infos {
		wg.Go(func() {
			lhsSec, rhsSec := computeOffsets(info, min, max)

			err := connectAudio(lhsSec, info.FilePath(), rhsSec, info.TmpFilePath())
			if err != nil {
				errors <- fmt.Errorf("failed to connect audio for %s: %w", info.RecordedUserID, err)
			}

			if err := os.Remove(info.FilePath()); err != nil {
				errors <- err
			}
			if err := os.Rename(info.TmpFilePath(), info.FilePath()); err != nil {
				errors <- err
			}
		})
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	log.Println("Successfully added silence")
	return nil
}

func computeOffsets(info recorder.RecordingInfo, min, max int64) (lhs, rhs float64) {
	tmpLhs := info.Timestamps.Join - min
	tmpRhs := max - info.Timestamps.Left

	if tmpLhs < 0 {
		tmpLhs = 0
	}
	if tmpRhs < 0 {
		tmpRhs = 0
	}

	lhsSec := float64(tmpLhs) / 1000.0
	rhsSec := float64(tmpRhs) / 1000.0

	return lhsSec, rhsSec
}

func craftSilence(durationSec float64) *ffmpeg_go.Stream {
	return ffmpeg_go.Input("anullsrc=r=48000:cl=mono",
		ffmpeg_go.KwArgs{
			"f": "lavfi",
			"t": fmt.Sprintf("%.3f", durationSec),
		},
	)
}

func connectAudio(lhsSec float64, audioPath string, rhsSec float64, outputPath string) error {
	var streams []*ffmpeg_go.Stream
	n := 0

	if lhsSec > 0 {
		streams = append(streams, craftSilence(lhsSec))
		n++
	}

	streams = append(streams, ffmpeg_go.Input(audioPath))
	n++

	if rhsSec > 0 {
		streams = append(streams, craftSilence(rhsSec))
		n++
	}

	return ffmpeg_go.Concat(streams, ffmpeg_go.KwArgs{
		"v": 0,
		"a": 1,
		"n": n,
	}).Output(outputPath, ffmpeg_go.KwArgs{"y": ""}).Run()
}

func findRange(infos []recorder.RecordingInfo) (min int64, max int64) {
	min = math.MaxInt64
	max = math.MinInt64

	for _, info := range infos {
		if info.Timestamps.Join < min {
			min = info.Timestamps.Join
		}
		if info.Timestamps.Left > max {
			max = info.Timestamps.Left
		}
	}

	return min, max
}

func mixPath(basePath, ext string) string {
	t := time.Now().UTC().UnixNano()
	return path.Join(basePath, fmt.Sprintf("merged_%d.%s", t, ext))
}
