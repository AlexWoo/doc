# RTMP Chunk Header

RTMP 是 FLV 的流式传播格式，类似于 H.264 中 TS 流与 PS 流的关系，一种便于网络传输，一种文件存储。RTMP Chunk Header 实际就是 FLV Tag 的变形，RTMP 中除了 Handshake 消息外，其它所有消息都使用 RTMP Chunk 的方式进行传输，每个 RTMP Chunk 都包含一个 RTMP Chunk Header，本部分主要对 RTMP Chunk Header 进行总结。RTMP 部分参考规范《Adobe’s Real Time Messaging Protocol》，FLV 部分参考《Video File Format Specification Version 10》

## Chunk Header 结构

一个 Chunk Header 包含：

- 1-3 字节的 Basic Header
- 0，3，7 或 11 字节的 Message Header
- 0 或 4 字节的 Extended Timestamp

### Basic Header

Basic Header 包含 fmt 和 csid 两个字段

**fmt**

fmt 包含 2 个bit：

- 0: Message Header 为 11 字节编码
- 1: Message Header 为 7 字节编码
- 2: Message Header 为 3 字节编码
- 3: Message Header 为 0 字节编码

**csid**

第一个字节中，剩余 6 bit 为 2 ~ 63 时，Basic Header 只有一个字节
第一个字节中，剩余 6 bit 为 0 时，Basic Header 为 2 个字节，csid 范围为 64 + byte2，即 64 ~ 319
第一个字节中，剩余 6 bit 为 1 时，Basic Header 为 3 个字节，csid 范围为 (byte3) * 256 + byte2 + 64，即 64-65599

### Message Header

**fmt = 0**

- timestamp: 3 字节，绝对时间戳(单位 ms)，如果大于等于 0xFFFFFF，该位设置为 0xFFFFFF，时间戳设置于 Extended Timestamp 中
- message length：3 字节，与 flv 中 flvtag 的 DataSize 一致，帧长度(不是本 Chunk 的长度)
- message type id：1 字节，帧类型，如 8 为音频帧，9 为视频帧，18 为 metadata 帧，与 flv 中 flvtag 的 TagType 一致
- msg stream id：4 字节，little-endian，用来标示同一个流，一般一个连接里保持不变，与 flv 中 flvtag 的 StreamID 一致，区别是 flv 中该字段为 3 字节，并且始终设置为 0。该字段在实际应用中基本没有意义

**fmt = 1**

- timestamp delta：3 字节，相对时间戳，当前帧与前一个 csid 相同帧时间戳差值
- message length：同上
- message type id：同上

**fmt = 2**

- timestamp delta：同上

**fmt = 3**

无 Messsage Header 部分

### Extended Timestamp

在 RTMP 中，当时间大于 0xffffff 时，Message Header 中的 timestamp 填为 0xffffff，真实时间戳填到 Extended Timestamp 中。否则，真实时间填在 timestamp 中。

而在 FLV 中，当时间大于 0xffffff 时，高 1 字节填在 FLV 的 TimestampExtended 中，低三字节填在 Timestamp 中。
