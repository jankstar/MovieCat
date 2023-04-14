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

export default {

  ConstAudioBPSOptions,

  ConstVideoBPSOptions,

  ConstRecorderSlicesOptions,

  ConstPlaybackRateOptions,
}
