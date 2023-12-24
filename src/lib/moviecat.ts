const ConstAudioBPSOptions = [
  { value: 32000, label: "32kbit/s" },
  { value: 44000, label: "44kbit/s" },
  { value: 48000, label: "48kbit/s" },
  { value: 64000, label: "64kbit/s" },
  { value: 96000, label: "96kbit/s" },
  { value: 128000, label: "128kbit/s" },
  { value: 196000, label: "196kbit/s" },
  { value: 256000, label: "256kbit/s" },
]

const ConstVideoBPSOptions = [
  { value: 1000000, label: "1Mbit/s" },
  { value: 2500000, label: "2,5Mbit/s" },
  { value: 6000000, label: "6Mbit/s" },
]

const ConstRecorderSlicesOptions = [
  { value: 1, label: "1 sec" },
  { value: 10, label: "10 sec" },
  { value: 30, label: "30 sec" },
  { value: 60, label: "1 min" },
  { value: 600, label: "10 min" },
  { value: 1800, label: "30 min" },
  { value: 3600, label: "1 h" },
  { value: 0, label: " undefine" },
]

const ConstPlaybackRateOptions = [
  { value: 0.25, label: "1/4" },
  { value: 0.5, label: "1/2" },
  { value: 1, label: "1" },
  { value: 2, label: "2" },
  { value: 4, label: "4" },
]

const ConstTrueMimeType = [
  { value: "video/webm", label: "WebM" },
  { value: "video/webm;codecs=vp8", label: "WebM Google VP8" },
  { value: "video/webm;codecs=vp9", label: "WebM Google VP9" },
  { value: "video/webm;codecs=daala", label: "WebM Daala" },
  { value: "video/webm;codecs=h264", label: "WebM H.264" },
  { value: "video/webm;codecs=av1", label: "WebM AV1" },
  { value: "video/mp4", label: "mp4" },
  { value: "video/mp4;codecs=vp8", label: "mp4 Google VP8" },
  { value: "video/mp4;codecs=vp9", label: "mp4 Google VP9" },
  { value: "video/mp4;codecs=h264", label: "mp4 H.264" },
  { value: "video/mpeg", label: "mpeg" },
];


function computeTime(iDuration: number): string {

  let
    seconds = Math.floor((iDuration) % 60),
    minutes = Math.floor((iDuration / (60)) % 60),
    hours = Math.floor((iDuration / (60 * 60)) % 24);

  return ("00" + hours).slice(-2) + ":" + ("00" + minutes).slice(-2) + ":" + ("00" + seconds).slice(-2);
}

export default {

  ConstAudioBPSOptions,

  ConstVideoBPSOptions,

  ConstRecorderSlicesOptions,

  ConstPlaybackRateOptions,

  ConstTrueMimeType,

  computeTime,
}
