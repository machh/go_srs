package app

import (
	"fmt"
	"errors"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
)

const SRS_PURE_AUDIO_GUESS_COUNT = 115

type SrsGopCache struct {
	enabled                  bool
	gopCache                    []*rtmp.SrsRtmpMessage
	cachedVideoCount         uint32
	audioAfterLastVideoCount uint32
}

func NewSrsGopCache() *SrsGopCache {
	return &SrsGopCache{
		enabled:true,
		gopCache:make([]*rtmp.SrsRtmpMessage, 0),
		cachedVideoCount:0,
		audioAfterLastVideoCount:0,
	}
}

func (this *SrsGopCache) set(enabled bool) {
	this.enabled = enabled
	if !this.enabled {
		this.clear()
	}
}

func (this *SrsGopCache) cache(msg *rtmp.SrsRtmpMessage) error {
	if !this.enabled {
		return nil
	}

	if msg.GetHeader().IsVideo() {
		if !flvcodec.VideoIsH264(msg.GetPayload()) {
			return errors.New("cache failed, video data is not h264")
		}
		this.cachedVideoCount++
		this.audioAfterLastVideoCount = 0
	}

	if this.pureAudio() {
		return errors.New("cache failed, video data must cached first")
	}

	if msg.GetHeader().IsAudio() {
		this.audioAfterLastVideoCount++
	}
	//only audio data in 3s?
	if this.audioAfterLastVideoCount > SRS_PURE_AUDIO_GUESS_COUNT {
		this.clear()
		return errors.New("cache failed, audio cache overflow detected")
	}
	//clear gop cache when got key frame
	if msg.GetHeader().IsVideo() && flvcodec.VideoIsKeyFrame(msg.GetPayload()) {
		this.clear()
		this.cachedVideoCount = 1
	}

	this.gopCache = append(this.gopCache, msg)
	return nil
}

func (this *SrsGopCache) clear() {
	this.gopCache = this.gopCache[0:0]
	this.cachedVideoCount = 0
	this.audioAfterLastVideoCount = 0
}

func (this *SrsGopCache) empty() bool {
	return len(this.gopCache) == 0
}

func (this *SrsGopCache) startTime() int64 {
	if this.empty() {
		return 0
	}

	return this.gopCache[0].GetTimestamp()
}

func (this *SrsGopCache) pureAudio() bool {
	return this.cachedVideoCount == 0
}

func (this *SrsGopCache) dump(consumer *SrsConsumer, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) error {
	for i:= 0; i < len(this.gopCache); i++ {
		consumer.Enqueue(this.gopCache[i], atc, jitterAlgorithm)
	}
	fmt.Println("****************dump count=", len(this.gopCache), "****************")
	return nil
}
