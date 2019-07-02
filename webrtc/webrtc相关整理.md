# webrtc 相关整理

杂记，后续会持续增加

## 规范及资料

### webrtc

- [基础概念](https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24)：WebRTC 基础概念，模型，交互逻辑等
- [W3C Webrtc](https://www.w3.org/TR/webrtc)：W3C webrtc 官方文档
- [RFC 8445](https://tools.ietf.org/html/rfc8445)：ICE

### SDP 相关

- [RFC 4566](https://tools.ietf.org/html/rfc4566)：SDP 基础协议
- [RFC 3264](https://tools.ietf.org/html/rfc3264)：SDP offer answer 协商基础协议
- [RFC 5576](https://tools.ietf.org/html/rfc5576)：ssrc 和 ssrc-group
- [RFC 5761](https://tools.ietf.org/html/rfc5761)：RTP 和 RTCP 复用传输通道
- [Bundle](https://tools.ietf.org/html/draft-ietf-mmusic-sdp-bundle-negotiation-38)：Bundle 属性
- [Unified Plan](https://tools.ietf.org/html/draft-roach-mmusic-unified-plan-00)：Unified Plan
- [WMS](https://tools.ietf.org/html/draft-alvestrand-rtcweb-msid-02)：WebRTC Media Stream ID

### 媒体相关

- [RFC 3550](https://tools.ietf.org/html/rfc3550)：RTP 协议
- [RFC 3551](https://tools.ietf.org/html/rfc3551)：RTCP 协议
- [RFC 4588](https://tools.ietf.org/html/rfc4588)：RTP 重传，rtx
- [RFC 5109](https://tools.ietf.org/html/rfc5109)：RTP 前向纠错，ULPFEC

## SDP

### Stream

- msid-semantic

	- [a=msid-semantic: WMS \<streamid\>](https://tools.ietf.org/html/draft-alvestrand-rtcweb-msid-02#section-4)

	WMS 后跟的即为 streamid，该属性在 Session 层，在 PlanA 中一般使用该参数识别 StreamId。对于 PlanB 和 Unified-Plan 有多个 stream 的情况，一般不适用。

- ssrc

	- a=ssrc:\<ssrcid\> cname:\<cname\>
	- a=ssrc:\<ssrcid\> msid:\<streamid\> \<trackid\>
	- a=ssrc:\<ssrcid\> mslabel:\<streamid\>
	- a=ssrc:\<ssrcid\> label:\<trackid\>

	通过 ssrc 中的 msid，可以直接拿到媒体对应的 streamid 和 trackid，该属性一般放在 Media 层

- msid

	- [a=msid:\<streamid\> \<trackid\>](https://tools.ietf.org/html/draft-ietf-mmusic-msid-16)

	一般 PlanB 和 Unified Plan 会使用 msid 属性来标识 streamid 和 trackid

### SSRC

参考规范 [https://tools.ietf.org/html/rfc5576](https://tools.ietf.org/html/rfc5576)

SDP 中 Media 层中的 ssrc 属性是与 RTP 中的 SSRC 对应的，不同的 ssrc 标识着

- ssrc

	- a=ssrc:\<ssrcid\> cname:\<cname\>
	- a=ssrc:\<ssrcid\> msid:\<streamid\> \<trackid\>
	- a=ssrc:\<ssrcid\> mslabel:\<streamid\>
	- a=ssrc:\<ssrcid\> label:\<trackid\>

	ssrcid 与 RTP 包头中的 SSRC 对应，一般常见的四个属性即为 cname，msid，mslabel 和 label。在 PlanB 中，一个 media 下，可能存在多个 ssrc + msid 的属性，用于标识多个 stream

		a=ssrc:3764413287 cname:nW8yslSe7TqdnK4L
		a=ssrc:3764413287 msid:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG 3ab09921-fd9e-46ab-9cf1-b8b252821b20
		a=ssrc:3764413287 mslabel:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG
		a=ssrc:3764413287 label:3ab09921-fd9e-46ab-9cf1-b8b252821b20

- ssrc-group

	ssrc-group 用于在同一 media 层下，同一个 stream 中有两个 ssrc 的情况(目前还未知两个 SSRC 是出于什么考虑)，和 ssrc 属性配合情况如下

		a=ssrc-group:FID 1838686046 3030128730
		a=ssrc:1838686046 cname:nW8yslSe7TqdnK4L
		a=ssrc:1838686046 msid:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG d394365a-b865-4c3c-b1d3-7db2e007a5fb
		a=ssrc:1838686046 mslabel:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG
		a=ssrc:1838686046 label:d394365a-b865-4c3c-b1d3-7db2e007a5fb
		a=ssrc:3030128730 cname:nW8yslSe7TqdnK4L
		a=ssrc:3030128730 msid:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG d394365a-b865-4c3c-b1d3-7db2e007a5fb
		a=ssrc:3030128730 mslabel:nNI3GmkRwoAm74BVarqgukAt4BaHqDwT68MG
		a=ssrc:3030128730 label:d394365a-b865-4c3c-b1d3-7db2e007a5fb

### BUNDLE

- mid

	- a=mid:\<string\>

	mid 即 media ID，每个 Media 层中一个 mid，形如：

		a=mid:0

- group BUNDLE

	- a=group:BUNDLE <mid-value> <mid-value>

	group BUNDLE 在 session 层定义，后面的 mid-value 是每个 media 层的 mid 值，形如：

		a=group:BUNDLE 0 1
	
	该属性用于标识多个媒体行使用相同的 UDP 端口进行传输，在目前应用场景中，客户端需要通过打洞才能实现互通，所以一般都会带有 group BUNDLE 属性