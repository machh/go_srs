package rtmp

type SrsRtmpMessage struct {
	// 4.1. Message Header
	header SrsMessageHeader
	// 4.2. Message Payload
	/**
	 * current message parsed size,
	 *       size <= header.payload_length
	 * for the payload maybe sent in multiple chunks.
	 */
	size int32
	/**
	 * the payload of message, the SrsCommonMessage never know about the detail of payload,
	 * user must use SrsProtocol.decode_message to get concrete packet.
	 * @remark, not all message payload can be decoded to packet. for example,
	 *       video/audio packet use raw bytes, no video/audio packet.
	 */
	payload []byte
	/**
     * Four-byte field that contains a timestamp of the message.
     * The 4 bytes are packed in the big-endian order.
     * @remark, used as calc timestamp when decode and encode time.
     * @remark, we use 64bits for large time for jitter detect and hls.
     */
	timestamp int64
}

func NewSrsRtmpMessage() *SrsRtmpMessage {
	return &SrsRtmpMessage{}
}

func (this *SrsRtmpMessage) DeepCopy() *SrsRtmpMessage {
	msg := &SrsRtmpMessage{
		header:this.header,
		size:this.size,
	}

	msg.payload = make([]byte, len(this.payload))
	copy(msg.payload, this.payload)
	return msg
}

func (this *SrsRtmpMessage) GetHeader() *SrsMessageHeader {
	return &(this.header)
}

func (this *SrsRtmpMessage) GetPayload() []byte {
	return this.payload
}

func (this *SrsRtmpMessage) SetTimestamp(t int64) {
	this.timestamp = t
}

func (this *SrsRtmpMessage) GetTimestamp() int64 {
	return this.timestamp
}