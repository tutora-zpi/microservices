package recorder

import (
	"context"
	"fmt"
	"path"
	"recorder-service/internal/infrastructure/webrtc/writer"
	"time"

	"github.com/pion/webrtc/v3"
)

type Recorder interface {
	StartRecording(ctx context.Context, roomID, userID string, track *webrtc.TrackRemote, newWriter writer.WriterFactory)
	StopRecording(userID string) *RecordingInfo
	StopAll()
}

type RecordingInfo struct {
	Timestamps     Timestamp
	BasePath       string
	RecordedUserID string
	Ext            string
}

func (r *RecordingInfo) FilePath() string {
	return path.Join(r.BasePath, fmt.Sprintf("%s.%s", r.RecordedUserID, r.Ext))
}

func (r *RecordingInfo) TmpFilePath() string {
	return path.Join(r.BasePath, fmt.Sprintf("tmp_%s.%s", r.RecordedUserID, r.Ext))
}

func NewRecordingInfo(detail *Detail, basePath, userID, ext string) RecordingInfo {
	return RecordingInfo{
		Timestamps:     detail.Timestamp,
		BasePath:       basePath,
		RecordedUserID: userID,
		Ext:            ext,
	}
}

type Timestamp struct {
	Join int64
	Left int64
}

type Detail struct {
	Writer    writer.Writer
	Cancel    context.CancelFunc
	Timestamp Timestamp
}

func NewDetail(w writer.Writer, c context.CancelFunc) *Detail {
	return &Detail{Writer: w, Cancel: c}
}

func (d *Detail) SetJoinTime() {
	d.Timestamp.Join = time.Now().UTC().UnixMilli()
}

func (d *Detail) SetLeftTime() {
	d.Timestamp.Left = time.Now().UTC().UnixMilli()
}
