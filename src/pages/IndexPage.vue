<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "IndexPage",
  props: ["langu"],

  components: {},

  data: () => {
    return {
      $t: undefined,
      //MediaDevices
      connectOn: false,
      selMode: "capture",
      SupportedConstraints: undefined,
      videoInput: [],
      videoInputValue: "default",
      audioInput: [],
      audioInputValue: "default",

      //Stream
      MainStream: <MediaStream>undefined,
      VideoSetting: undefined,
      VideoSettingNew: undefined,
      VideoCapabilities: undefined,
      VideoSettingDialog: false,

      AudioSetting: undefined,
      AudioSettingNew: undefined,
      AudioCapabilities: undefined,
      AudioSettingDialog: false,
      AudioDraw: false,

      //Recorder
      Recorder: <MediaRecorder>undefined,
      recorderOptions: {
        audioBitsPerSecond: 64000,
        videoBitsPerSecond: 2500000,
        mimeType: undefined,
      },
      audioBPSOptions: [
        { value: 32000, label: "32kbit/s" },
        { value: 44000, label: "44kbit/s" },
        { value: 48000, label: "48kbit/s" },
        { value: 64000, label: "64kbit/s" },
        { value: 128000, label: "128kbit/s" },
        { value: 196000, label: "196kbit/s" },
        { value: 256000, label: "256kbit/s" },
      ],
      videoBPSOptions: [
        { value: 1000000, label: "1Mbit/s" },
        { value: 2500000, label: "2,5Mbit/s" },
        { value: 6000000, label: "6Mbit/s" },
      ],
      mimeTypeOptions: [],
      // eslint-disable-next-line no-undef
      //RecorderOption: <MediaRecorderOptions>{
      //  mimeType: "",
      //},
      recorderSlices: 10, //in sec
      recorderSlicesOptions: [
        { value: 1, label: "1 sec" },
        { value: 10, label: "10 sec" },
        { value: 30, label: "30 sec" },
        { value: 60, label: "1 min" },
        { value: 600, label: "10 min" },
        { value: 1800, label: "30 min" },
        { value: 3600, label: "1 h" },
        { value: 0, label: " undefine" },
      ],
      recorderAutoStop: 30, //min
      RecorderStartTime: undefined,
      RecorderBlobList: [],
      RecorderInfoList: [],
      RecorderSize: 0,
      RecorderState: "inactive",

      PlayerPosition: 0, //Attention only integer in sec !
      durationPlayer: 0, //Attention only integer in sec !
      PlayerMuted: true,
      playOn: false,
      playbackRate: 1,
      playbackRateOptions: [
        { value: 0.25, label: "1/4" },
        { value: 0.5, label: "1/2" },
        { value: 1, label: "1" },
        { value: 2, label: "2" },
        { value: 4, label: "4" },
      ],
      TimeOutID: undefined,
      fileName: "test",
      OnOffOptions: [],
    };
  },

  created() {
    console.log(`created`);
  },

  async mounted() {
    try {
      console.log(`mounted`);
      this.OnOffOptions = [
        { value: true, label: this.$t("On") },
        { value: false, label: this.$t("Off") },
      ];
      await this.getMediaDevices();
      this.SupportedConstraints = navigator.mediaDevices.getSupportedConstraints();
      console.log(`SupportedConstraints: ${JSON.stringify(this.SupportedConstraints)}`);
      this.getMimeType();
      //recorder data from localStorage
      this.loadRecorderData();
    } catch (e) {
      console.error(e);
    }
  },

  updated() {
    console.log(`updated`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
  },

  methods: {
    //
    saveRecorderData() {
      console.log(`saveRecorderData()`);

      localStorage.setItem("recorderOptions", JSON.stringify(this.recorderOptions));
      localStorage.setItem("recorderSlices", JSON.stringify(this.recorderSlices));
      localStorage.setItem("recorderAutoStop", JSON.stringify(this.recorderAutoStop));
      //localStorage.setItem("fileName", JSON.stringify(this.fileName));
    },
    loadRecorderData() {
      console.log(`loadRecorderData()`);
      try {
        const JSONrecorderOptions = localStorage.getItem("recorderOptions");
        if (JSONrecorderOptions) {
          let recorderOptions = JSON.parse(JSONrecorderOptions);
          if (recorderOptions) {
            this.recorderOptions = recorderOptions;
          }
        }

        const JSONrecorderSlices = localStorage.getItem("recorderSlices");
        if (JSONrecorderOptions) {
          let recorderSlices = JSON.parse(JSONrecorderSlices);
          if (recorderSlices) {
            this.recorderSlices = recorderSlices;
          }
        }

        const JSONrecorderAutoStop = localStorage.getItem("recorderAutoStop");
        if (JSONrecorderAutoStop) {
          let recorderAutoStop = JSON.parse(JSONrecorderAutoStop);
          if (recorderAutoStop) {
            this.recorderAutoStop = recorderAutoStop;
          }
        }

        // const JSONfileName = localStorage.getItem("fileName");
        // if (JSONfileName) {
        //   let fileName = JSON.parse(JSONfileName);
        //   if (fileName) {
        //     this.fileName = fileName;
        //   }
        // }
      } catch (e) {
        console.error(e);
      }
    },
    //
    async connectOnBtn() {
      let that = this;
      console.log(`connectOnBtn(e)`);
      //      if (this.RecorderBlobList && this.RecorderBlobList.length > 0) {
      //        this.$q.notify({ type: "warning", message: this.$t("InfoRecorderMemory"), position: "center", timeout: 5000 });
      //      }

      //this.RecorderBlobList = [];
      //this.RecorderInfoList = [];
      //this.RecorderSize = 0;
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video");
      if (myVideoTag && !this.connectOn) {
        await this.startDisplayMedia();
        if (this.MainStream && this.MainStream.active) {
          this.MainStream.oninactive = () => {
            that.connectOffBtn();
          };
          myVideoTag.srcObject = this.MainStream;
          this.connectOn = true;
        }
      }
    },
    //
    // eslint-disable-next-line no-unused-vars
    connectOffBtn() {
      console.log(`connectOffBtn(e)`);

      this.MainStream.oninactive = undefined;
      this.stopDisplayMedia();
      this.stopRecorder();

      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video");
      if (myVideoTag) {
        myVideoTag.srcObject = undefined;
      }
      this.connectOn = false;
    },

    getMimeType() {
      let trueMimeType = [
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
      let that = this;
      trueMimeType.forEach((item: any) => {
        if (MediaRecorder.isTypeSupported(item.value)) {
          console.log(`codec ${item.value} supported`);

          that.mimeTypeOptions.push(item);
          if (!that.recorderOptions.mimeType) {
            that.recorderOptions.mimeType = item.value;
          }
        } else {
          console.log(`codec ${item.value} not supported`);
        }
      });
    },
    //
    async getMediaDevices() {
      let that = this;
      await navigator.mediaDevices
        .enumerateDevices()
        .then((devices) => {
          that.videoInput = [];
          that.videoInputValue = undefined;
          that.audioInput = [];
          that.audioInputValue = undefined;

          let countCam = 1;
          let countMic = 1;
          let videoinputDefault = undefined;
          let audioinputDefault = undefined;

          this.videoInput.push({ label: `${this.$t("Video_off")}`, value: "off" }); //video off
          this.audioInput.push({ label: `${this.$t("Audio_off")}`, value: "off" }); //audio off

          devices.forEach((device) => {
            if (device.kind == "videoinput") {
              let item = { label: device.label.substring(0, 40) || `Camera ${countCam++}`, value: device.deviceId || "default" };
              if (device.deviceId == "default") {
                item.label = "Video default";
                videoinputDefault = item;
              }
              that.videoInput.push(item);
            }

            if (device.kind == "audioinput") {
              let item = { label: device.label.substring(0, 40) || `Mic ${countMic++}`, value: device.deviceId || "default" };
              if (device.deviceId == "default") {
                item.label = "Audio default";
                audioinputDefault = item;
              }
              that.audioInput.push(item);
            }

            console.log(`${device.kind}: ${device.label} id = ${device.deviceId}`);
          });

          if (videoinputDefault == undefined) {
            let item = { label: `Video default`, value: "default" };
            that.videoInput.push(item);
            videoinputDefault = item;
          }
          if (!that.videoInputValue) {
            that.videoInputValue = videoinputDefault;
          }

          if (audioinputDefault == undefined) {
            let item = { label: `Audio default`, value: "default" };
            that.audioInput.push(item);
            audioinputDefault = item;
          }
          if (!that.audioInputValue) {
            that.audioInputValue = audioinputDefault;
          }
        })
        .catch((err) => {
          console.error(`${err.name}: ${err.message}`);
          this.$q.notify({ type: "negative", message: `${err.name}: ${err.message}`, position: "center", timeout: 5000 });
        });
    },
    //
    async startDisplayMedia() {
      console.log("startDisplayMedia()");
      let that = this;
      try {
        let configuration = { audio: <boolean | any>true, video: <boolean | any>true, systemAudio: "include", surfaceSwitching: "include", selfBrowserSurface: "exclude" };
        if (!this.audioInputValue || this.audioInputValue == "default") {
          configuration.audio = true;
        } else if (this.audioInputValue == "off") {
          configuration.audio = false;
        } else {
          configuration.audio = {
            deviceId: this.audioInputValue || "default",
          };
        }

        if (this.selMode == "capture") {
          //mode capture
          this.MainStream = await navigator.mediaDevices.getDisplayMedia(configuration);
        } else if (this.selMode == "camera") {
          //mode camera
          if (this.videoInputValue == "off") {
            configuration.video = false;
          } else {
            configuration.video = {
              deviceId: this.videoInputValue || "default",
            };
          }
          if (configuration.video == false && configuration.audio == false) {
            throw new Error("audio and video off - no media stream");
          }
          this.MainStream = await navigator.mediaDevices.getUserMedia(configuration);
        } else {
          // error
          return;
        }

        this.newMediaRecorder();
        let lTracks = this.MainStream.getVideoTracks();
        if (lTracks && lTracks.length > 0 && lTracks[0]) {
          this.VideoSetting = lTracks[0].getSettings();
          this.VideoSetting.enabled = lTracks[0].enabled;
          console.log("video getSettings()", JSON.stringify(this.VideoSetting));
          console.log("video getCapabilities()", JSON.stringify(lTracks[0].getCapabilities()));
        }

        lTracks = this.MainStream.getAudioTracks();
        if (lTracks && lTracks.length == 0 && this.audioInputValue && this.audioInputValue != "default" && this.audioInputValue != "off") {
          //audio was not detected via capture - add it manually
          let audioStream = await navigator.mediaDevices.getUserMedia({
            audio: {
              deviceId: this.audioInputValue,
            },
            video: false,
          });
          audioStream.getTracks().forEach((my_track: any) => {
            that.MainStream.addTrack(my_track, audioStream);
          });
          lTracks = this.MainStream.getAudioTracks();
        }
        if (lTracks && lTracks.length > 0 && lTracks[0]) {
          this.AudioSetting = lTracks[0].getSettings();
          this.AudioSetting.muted = lTracks[0].muted || !lTracks[0].enabled;
          console.log("audio getSettings()", JSON.stringify(this.AudioSetting));
          console.log("audio getCapabilities()", JSON.stringify(lTracks[0].getCapabilities()));
          this.drawCanvas("audio");
        }

        //
      } catch (err) {
        this.MainStream = undefined;
        console.error(`${err.name}: ${err.message}`);
        return { error: err };
      }
      return { error: undefined };
    },

    newMediaRecorder() {
      try {
        if (!this.MainStream) {
          throw new Error("no media stream");
        }

        this.Recorder = new MediaRecorder(this.MainStream, this.recorderOptions);
        if (this.RecorderBlobList && this.RecorderBlobList.length > 0) {
          this.$q.notify({ type: "warning", message: this.$t("InfoRecorderMemory"), position: "center", timeout: 5000 });
        }
        // this.RecorderBlobList = [];
        // this.RecorderInfoList = [];
        // this.RecorderSize = 0;
        // this.RecorderStartTime = undefined;
        let that = this;

        this.Recorder.ondataavailable = (e) => {
          console.log(`type: ${e.type} time: ${e.timecode} size: ${e.data.size}`);
          if (e.data && e.data.size > 0) {
            that.RecorderBlobList.push(e.data);
            that.RecorderInfoList.push({ type: e.type, time: Date.now(), size: e.data.size });
            that.RecorderSize += e.data.size;
          }
        };

        //the events no longer function properly after a stop, so the recoder is completely restarted!
        this.Recorder.onstop = () => {
          console.log("onstop()");
          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
        };

        this.Recorder.onerror = (e) => {
          console.log("onerror()", e);
          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
        };
        this.Recorder.onpause = () => {
          console.log("onpause()");

          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
        };
        this.Recorder.onresume = () => {
          console.log("onresume()");

          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
        };
        // eslint-disable-next-line no-unused-vars
        this.Recorder.onstart = () => {
          console.log("onstart()");

          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
        };
      } catch (err) {
        console.error(`${err.name}: ${err.message}`);
      }
    },

    stopDisplayMedia() {
      console.log(`stopDisplayMedia()`);
      if (this.MainStream && this.MainStream.active) {
        this.stopRecorder();
        let tracks = this.MainStream.getTracks();
        tracks.forEach((track) => track.stop());
        this.MainStream = undefined;
        this.VideoSetting = undefined;
        this.AudioSetting = undefined;
        this.Recorder = undefined;
      }
    },
    //
    AudioSettingBtn() {
      let [traks] = this.MainStream.getAudioTracks();
      if (traks) {
        console.log("AudioSettingBtn()");
        this.AudioCapabilities = traks.getCapabilities();
        this.AudioSettingNew = JSON.parse(JSON.stringify(this.AudioSetting));
        this.AudioSettingDialog = true;
      }
    },
    //
    async AudioSettingChangeBtn() {
      let [tracks] = this.MainStream.getAudioTracks();
      if (tracks) {
        console.log("AudioSettingChangeBtn()");
        try {
          await tracks.applyConstraints(
            Object.assign(tracks.getSettings(), {
              autoGainControl: this.AudioSettingNew.autoGainControl,
              echoCancellation: this.AudioSettingNew.echoCancellation,
              noiseSuppression: this.AudioSettingNew.noiseSuppression,
              sampleRate: this.AudioSettingNew.sampleRate,
            })
          );
        } catch (e) {
          console.error(e);
          this.$q.notify({ type: "negative", message: `${e.name}: ${e.message}`, position: "center", timeout: 5000 });
        }
        tracks.enabled = !this.AudioSettingNew.muted;

        this.AudioSetting = tracks.getSettings();
        this.AudioSetting.muted = tracks.muted || !tracks.enabled;
      }
      this.AudioSettingDialog = false;
    },
    //
    VideoSettingBtn() {
      let [tracks] = this.MainStream.getVideoTracks();
      if (tracks) {
        console.log("VideoSettingBtn()");
        this.VideoCapabilities = tracks.getCapabilities();
        this.VideoSettingNew = JSON.parse(JSON.stringify(this.VideoSetting));
        this.VideoSettingDialog = true;
      }
    },
    //
    async VideoSettingChangeBtn() {
      let [tracks] = this.MainStream.getVideoTracks();
      if (tracks) {
        console.log("VideoSettingChangeBtn()");
        try {
          await tracks.applyConstraints(
            Object.assign(tracks.getSettings(), {
              height: this.VideoSettingNew.height,
              width: this.VideoSettingNew.width,
              frameRate: this.VideoSettingNew.frameRate,
            })
          );
        } catch (e) {
          console.error(e);
          this.$q.notify({ type: "negative", message: `${e.name}: ${e.message}`, position: "center", timeout: 5000 });
        }
        tracks.enabled = this.VideoSettingNew.enabled;

        this.VideoSetting = tracks.getSettings();
        this.VideoSetting.enabled = tracks.enabled;
      }
      this.VideoSettingDialog = false;
    },
    //
    startRecorderBtn() {
      console.log(`startRecorderBtn()`);
      let that = this;

      this.newMediaRecorder();

      if (this.Recorder && this.Recorder.state == "inactive") {
        if (this.recorderSlices && this.recorderSlices > 0) {
          this.Recorder.start(1000 * this.recorderSlices); //1 * x sec Slices
        } else {
          this.Recorder.start(); //ohne slices
        }

        if (!this.RecorderStartTime) {
          this.RecorderStartTime = Date.now();
        }

        this.TimeOutID = setTimeout(function () {
          that.stopRecorder();
          that.download();
          that.TimeOutID = undefined;
        }, 1000 * 60 * this.recorderAutoStop); //150 min = 2,5h
      }
    },
    //
    stopRecorderBtn() {
      console.log(`stopRecorderBtn()`);
      this.stopRecorder();
      if (this.TimeOutID) {
        clearTimeout(this.TimeOutID);
        this.TimeOutID = undefined;
      }
    },
    //
    stopRecorder() {
      console.log(`stopRecorder()`);
      if (this.Recorder && this.Recorder.state != "inactive") {
        this.Recorder.requestData();
        this.Recorder.stop();
      }
      this.RecorderState = this.Recorder && this.Recorder.state ? this.Recorder.state : "inactive";
    },
    //
    download() {
      console.log(`download: ${this.RecorderBlobList.length} blobs`);
      if (this.RecorderBlobList && this.RecorderBlobList.length > 0) {
        var blob = new Blob(this.RecorderBlobList, { type: this.recorderOptions.mimeType });
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement("a");
        a.style.display = "none";
        a.href = url;
        a.download = `${this.fileName}.${this.recorderOptions.mimeType.split(";")[0].split("/")[1]}`;
        console.log(`filename ${a.download}`);
        document.body.appendChild(a);
        a.click();
        setTimeout(function () {
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);
        }, 3000);
      }
    },
    //
    upload() {
      console.log(`upload()`);
    },
    //
    clearBuffer() {
      this.RecorderBlobList = [];
      this.RecorderInfoList = [];
      this.RecorderSize = 0;
    },
    //
    uuidv4(): string {
      return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
        var r = (Math.random() * 16) | 0,
          v = c == "x" ? r : (r & 0x3) | 0x8;
        return v.toString(16);
      });
    },
    //
    uuidToColor(id: string) {
      var g = 0,
        b = 0;
      for (var i = 0; i < id.length / 2; i++) {
        var code = id.charCodeAt(i);
        g = g + code;
        code = id.charCodeAt(i * 2);
        b = b + code;
      }
      return [g % 256, b % 256];
    },
    drawCanvas(id: string) {
      let canvas = document.getElementById(`canvas-${id}`);
      if (canvas && this.MainStream && this.MainStream.getTracks().length > 0) {
        this.visualize(id, canvas);
      }
    },
    //
    visualize(id: string, canvas: HTMLCanvasElement) {
      if (this.AudioDraw == true) {
        return;
      }
      var that = this;

      this.AudioDraw = true;
      console.log("visualize " + id);

      var audioCtx = new AudioContext();

      var analyser = audioCtx.createAnalyser();
      analyser.fftSize = 256;
      analyser.minDecibels = -80;
      analyser.maxDecibels = -10;
      analyser.smoothingTimeConstant = 0.85;
      var audioStream = new MediaStream(this.MainStream.getAudioTracks());
      audioCtx.createMediaStreamSource(audioStream).connect(analyser);

      var canvasCtx = canvas.getContext("2d");
      var bufferLength = analyser.frequencyBinCount;
      var dataArray = new Float32Array(bufferLength);
      var gb = this.uuidToColor(this.uuidv4());
      var g = gb[0],
        b = gb[1];
      var MIN = 7;

      function draw() {
        var WIDTH = canvas.width;
        var HEIGHT = canvas.height;

        analyser.getFloatFrequencyData(dataArray);

        canvasCtx.fillStyle = "rgb(0, 0, 0)";
        canvasCtx.fillRect(0, 0, WIDTH, HEIGHT);

        var barWidth = (WIDTH / bufferLength) * 2.5;
        var barHeight,
          point,
          x = 0;

        for (var i = 0; i < bufferLength; i++) {
          point = dataArray[i];
          barHeight = (point + 140) * 2;

          var r = Math.floor(barHeight + 64);
          if (g % 3 === 0) {
            canvasCtx.fillStyle = `rgb(${r},${g},${b})`;
          } else if (g % 3 === 1) {
            canvasCtx.fillStyle = `rgb(${g},${r},${b})`;
          } else {
            canvasCtx.fillStyle = `rgb(${g},${b},${r})`;
          }

          barHeight = HEIGHT / MIN + ((barHeight / 256) * HEIGHT * (MIN - 1)) / MIN;
          if (barHeight < HEIGHT / MIN) {
            barHeight = HEIGHT / MIN;
          }
          canvasCtx.fillRect(x, HEIGHT - barHeight, barWidth, barHeight);

          x += barWidth + 1;
        }

        var el = document.getElementById("peer-" + id);
        if (el && that.AudioSetting) {
          setTimeout(function () {
            requestAnimationFrame(draw);
          }, 50);
        } else {
          that.AudioDraw = false;
          canvasCtx.fillStyle = "rgb(0, 0, 0)";
          canvasCtx.fillRect(0, 0, WIDTH, HEIGHT);
          audioCtx = undefined;
          analyser = undefined;
        }
      }

      draw();
    },
    //
    playOnBtn() {
      console.log(`playOnBtn()`);
      let that = this;
      this.durationPlayer = 0;
      if (this.RecorderInfoList.length > 0) {
        this.durationPlayer = ((this.RecorderInfoList[this.RecorderInfoList.length - 1].time - this.RecorderInfoList[0].time) / 1000).toFixed(0);
      }
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (this.playOn == false && myVideoTag && this.RecorderBlobList && this.RecorderBlobList.length > 0 && this.recorderOptions.mimeType) {
        myVideoTag.pause(); //player stoppen
        //
        myVideoTag.onloadedmetadata = () => {
          console.log("onloadedmetadata()");
          that.durationPlayer = Number.isFinite(myVideoTag.duration) ? myVideoTag.duration.toFixed(0) : that.durationPlayer;
          console.log("Duration change", that.durationPlayer);
        };
        myVideoTag.ondurationchange = () => {
          console.log("ondurationchange()");
          that.durationPlayer = Number.isFinite(myVideoTag.duration) ? myVideoTag.duration.toFixed(0) : that.durationPlayer;
          console.log("Duration change", that.durationPlayer);
        };
        myVideoTag.ontimeupdate = () => {
          console.log(`currentTime ${myVideoTag.currentTime.toFixed(2)} sec`);
          that.PlayerPosition = myVideoTag.currentTime ? myVideoTag.currentTime.toFixed(0) : 0;
          if (that.durationPlayer < that.PlayerPosition) {
            that.durationPlayer = that.PlayerPosition + 1;
          }
        };
        myVideoTag.onended = () => {
          that.playOn = false;
        };
        myVideoTag.onplay = () => {
          that.playOn = true;
        };
        //
        //        var url = window.URL.createObjectURL(this.RecorderBlobList[this.PlayerPosition]);
        let blob = new Blob(this.RecorderBlobList, { type: this.recorderOptions.mimeType });
        let url = window.URL.createObjectURL(blob);
        myVideoTag.src = url;
        myVideoTag.currentTime = this.PlayerPosition;
        myVideoTag.playbackRate = this.playbackRate;
        myVideoTag.play();
      }
    },
    playOffBtn() {
      console.log(`playOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        myVideoTag.pause(); //player stoppen
      }
      this.playOn = false;
    },
    muteOnBtn() {
      console.log(`muteOnBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        this.PlayerMuted = true;
        myVideoTag.muted = true;
      }
    },
    muteOffBtn() {
      console.log(`muteOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        this.PlayerMuted = false;
        myVideoTag.muted = false;
      }
    },
    changePlaybackRate(value: any) {
      console.log(`muteOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        myVideoTag.playbackRate = value;
      }
    },
    computeTime(iDuration: number) {
      let l_time = Math.trunc(iDuration);
      let l_sec = (l_time % 60).toFixed(0);
      l_sec = l_sec.length > 1 ? l_sec : `0${l_sec}`;
      l_time = Math.trunc(iDuration / 60);
      let l_min = (l_time % 60).toFixed(0);
      l_min = l_min.length > 1 ? l_min : `0${l_min}`;
      let l_std = (l_time / 60).toFixed(0);
      l_std = l_std.length > 1 ? l_std : `0${l_std}`;
      return `${l_std}:${l_min}:${l_sec}`;
    },
  },
});
</script>

<template>
  <q-page class="q-pa-md">
    <div class="tw-grid-flow tw-items-start row tw-gap-4">
      <!-- input and mode -->
      <q-card style="min-width: 150px; max-width: 300px">
        <q-card-section>
          <q-select
            v-model="selMode"
            :options="[
              { label: $t('Capture'), value: 'capture' },
              { label: $t('Camera'), value: 'camera' },
              { label: $t('Player'), value: 'player' },
            ]"
            emit-value
            map-options
            :label="$t('Mode')"
            :disable="connectOn"
          />
          <br />
        </q-card-section>
        <q-separator v-if="selMode == 'capture' || selMode == 'camera'" />
        <q-card-section>
          <div v-if="selMode == 'capture' || selMode == 'camera'">
            <h7 class="text-subtitle1">{{ $t("Input_state") }} {{ connectOn ? $t("On") : $t("Off") }}</h7> <br />
          </div>
          <q-select v-if="selMode == 'camera'" v-model="videoInputValue" :options="videoInput" emit-value map-options :label="$t('Video_Input')" :disable="connectOn" />
          <q-select
            v-if="selMode == 'camera' || selMode == 'capture'"
            v-model="audioInputValue"
            :options="audioInput"
            emit-value
            map-options
            :label="$t('Audio_Input')"
            :disable="connectOn"
          />
        </q-card-section>
        <q-separator />

        <q-card-section v-if="AudioSetting && connectOn">
          <div class="row tw-justify-between">
            <h7 class="text-subtitle1">Audio</h7>
            <q-separator />
            <q-btn flat round sizes="sx" padding="none" icon="edit" @click="AudioSettingBtn" />
          </div>
          <h7 class="text-body2">{{ $t("AudioMuted") }}: {{ AudioSetting && AudioSetting.muted == false ? $t("Off") : $t("On") }}</h7> <br />
          <h7 class="text-body2">{{ $t("AutoGainControl") }}: {{ AudioSetting && AudioSetting.autoGainControl == false ? $t("Off") : $t("On") }}</h7> <br />
          <h7 class="text-body2">{{ $t("NoiseSuppression") }}: {{ AudioSetting && AudioSetting.noiseSuppression == false ? $t("Off") : $t("On") }}</h7> <br />
          <h7 class="text-body2">{{ $t("EchoCancellation") }}: {{ AudioSetting && AudioSetting.echoCancellation == false ? $t("Off") : $t("On") }}</h7> <br />
          <h7 class="text-body2">{{ $t("SampleRate") }}: {{ AudioSetting ? (AudioSetting.sampleRate / 1000).toFixed(0) : "0" }} kB/sec</h7> <br />
        </q-card-section>
        <q-separator v-if="AudioSetting && connectOn" />
        <q-card-section v-if="VideoSetting && connectOn">
          <div class="row tw-justify-between">
            <h7 class="text-subtitle1">Video </h7>
            <q-separator />
            <q-btn flat round sizes="sx" padding="none" icon="edit" @click="VideoSettingBtn" />
          </div>
          <h7 class="text-body2">{{ $t("Resolution") }}: {{ VideoSetting ? `${VideoSetting.width}x${VideoSetting.height}` : "0x0" }}</h7> <br />
          <h7 class="text-body2">{{ $t("Frame_rate") }}: {{ VideoSetting ? VideoSetting.frameRate.toFixed(2) : "0" }}</h7> <br />
        </q-card-section>
        <q-separator v-if="VideoSetting && connectOn" />

        <q-card-actions v-if="selMode == 'camera' || selMode == 'capture'" class="tw-justify-end">
          <q-btn
            :label="!connectOn ? $t('Connect') : $t('Disconnect')"
            @click="!connectOn ? connectOnBtn() : connectOffBtn()"
            :class="!connectOn ? 'tw-bg-lime-300' : 'tw-bg-red-300'"
            :icon="!connectOn ? 'link' : 'link_off'"
          />
        </q-card-actions>
      </q-card>

      <!-- input audio / Video-->
      <q-card v-if="selMode == 'capture' || selMode == 'camera'" style="max-width: 300px">
        <q-card-section>
          <h7 class="text-subtitle1">Audio ({{ AudioSetting && !AudioSetting.muted ? $t("On") : $t("Off") }})</h7>

          <div class="peer" id="peer-audio">
            <div class="stat-value">
              <canvas class="peer_canvas rounded" id="canvas-audio" style="background-color: black; width: 268px"> </canvas
              ><!-- the canvas is filled via the programme-->
            </div>
          </div>
        </q-card-section>
        <q-card-section>
          <h7 class="text-subtitle1">Video ({{ VideoSetting && VideoSetting.enabled ? $t("On") : $t("Off") }})</h7> <br />
          <video id="id_video" playsinline autoplay muted style="background-color: black"></video>
        </q-card-section>
      </q-card>

      <!-- play video from file/recorder -->
      <q-card v-if="!connectOn && selMode == 'player'" style="max-width: 300px">
        <q-card-section>
          <h7 class="text-subtitle1">Video Player </h7> <br />
          <video id="id_video_player" playsinline muted style="background-color: black"></video>
        </q-card-section>
        <q-card-section>
          <h7 v-if="RecorderInfoList" class="text-body2">{{ $t("Time") }}: {{ computeTime(PlayerPosition) }}</h7>
          <br />
          <q-slider v-model="PlayerPosition" :min="0" :max="durationPlayer" label style="width: 90%" />
          <br />
          <q-select
            v-model="playbackRate"
            :options="playbackRateOptions"
            emit-value
            map-options
            :label="$t('PlaybackRate')"
            style="width: 50%"
            @update:model-value="changePlaybackRate"
          />
        </q-card-section>
        <q-separator />

        <q-card-actions class="tw-justify-end">
          <!-- -->
          <q-btn :label="!PlayerMuted ? $t('Mute') : $t('Mute_off')" @click="!PlayerMuted ? muteOnBtn() : muteOffBtn()" :icon="!PlayerMuted ? 'volume_up' : 'volume_off'" />
          <q-btn
            :label="!playOn ? $t('Play') : $t('Stop')"
            @click="!playOn ? playOnBtn() : playOffBtn()"
            :icon="!playOn ? 'play_circle' : 'stop_circle'"
            :class="!playOn ? 'tw-bg-lime-300' : 'tw-bg-red-300'"
            :disable="!RecorderBlobList || RecorderBlobList.length == 0"
          />
        </q-card-actions>
      </q-card>

      <!-- recorder -->
      <q-card v-if="selMode == 'capture' || selMode == 'camera' || selMode == 'player' || RecorderSize > 0" style="max-width: 300px">
        <q-card-section>
          <h7 class="text-subtitle1">{{ $t("Recorder_state") }} {{ RecorderState }}</h7> <br />

          <h7 class="text-body2">{{ `${$t("Size")} ${(RecorderSize / 1000000).toFixed(2)}` }} mByte</h7><br />
          <h7 class="text-body2"
            >{{ $t("Time") }}:
            {{
              RecorderStartTime && RecorderInfoList && RecorderInfoList.length > 0
                ? computeTime((RecorderInfoList[RecorderInfoList.length - 1].time - RecorderStartTime) / 1000)
                : "00:00:00"
            }}
          </h7>
          <h7 style="font-size: 10px"> ({{ $t("InfoRefreshPerSlices") }})</h7>

          <q-select
            v-model="recorderOptions.mimeType"
            :options="mimeTypeOptions"
            emit-value
            map-options
            :label="$t('Recorder_mimetype')"
            :disable="mimeTypeOptions.length == 0 || RecorderSize > 0 || RecorderState != 'inactive'"
          />
          <div v-if="selMode == 'capture' || selMode == 'camera'">
            <div class="row">
              <q-select
                v-model="recorderOptions.audioBitsPerSecond"
                :options="audioBPSOptions"
                emit-value
                map-options
                :label="$t('Audio_BPS')"
                :disable="RecorderState != 'inactive'"
                style="width: 50%"
                @update:model-value="saveRecorderData"
              />
              <q-select
                v-model="recorderOptions.videoBitsPerSecond"
                :options="videoBPSOptions"
                emit-value
                map-options
                :label="$t('Video_BPS')"
                :disable="RecorderState != 'inactive'"
                style="width: 50%"
                @update:model-value="saveRecorderData"
              />
            </div>
            <div class="row">
              <q-select
                v-model="recorderSlices"
                :options="recorderSlicesOptions"
                emit-value
                map-options
                :label="$t('Slice_length')"
                :disable="RecorderState != 'inactive'"
                style="width: 50%"
                @update:model-value="saveRecorderData"
              />
              <q-input
                v-model="recorderAutoStop"
                :label="$t('Auto_stop')"
                :disable="RecorderState != 'inactive'"
                type="number"
                mask="###"
                fill-mask="#"
                reverse-fill-mask
                style="width: 50%"
                :rules="[(val) => !!val || '* Required', (val) => val > 0 || 'The minimum is 1 min', (val) => val < 151 || 'The maximum is 150 min']"
                @update:model-value="saveRecorderData"
              />
            </div>
          </div>
          <q-input v-model="fileName" :label="$t('Filename')" />
        </q-card-section>
        <q-separator />

        <q-card-actions class="tw-justify-end tw-gap-2">
          <q-btn
            v-if="RecorderState == 'inactive' && (selMode == 'capture' || selMode == 'camera')"
            label="Rec on"
            icon="videocam"
            @click="startRecorderBtn"
            :disable="RecorderState != 'inactive' || !connectOn || this.recorderAutoStop > 150"
            class="tw-bg-lime-300"
          />
          <q-btn
            v-else-if="selMode == 'capture' || selMode == 'camera'"
            label="Rec off"
            icon="videocam_off"
            @click="stopRecorderBtn"
            :disable="RecorderState == 'inactive' || !connectOn"
            class="tw-bg-red-300"
          />
          <!--div v-if="RecorderSize != 0 && RecorderState == 'inactive'"-->
          <q-btn label="download" @click="download" :disable="RecorderSize == 0 || RecorderState != 'inactive'" icon="download" />
          <!--q-btn v-if="selMode == 'player'" label="upload" @click="upload" icon="upload" /-->
          <q-btn label="clear" @click="clearBuffer" :disable="RecorderSize == 0 || RecorderState != 'inactive'" icon="delete_forever" />
          <!--/div-->
        </q-card-actions>
      </q-card>
      <!--


        Dialog AudioSettingDialog
      -->
      <q-dialog v-model="AudioSettingDialog" no-backdrop-dismiss persistent class="tw-font-sans">
        <q-card style="width: 400px; height: 320px">
          <q-card-section>
            <div class="text-h6">Audio Settings</div>
            <div class="row">
              <q-select v-model="AudioSettingNew.muted" :options="OnOffOptions" emit-value map-options :label="$t('AudioMuted')" style="width: 45%" />
              <q-select v-model="AudioSettingNew.autoGainControl" :options="OnOffOptions" emit-value map-options :label="$t('AutoGainControl')" style="width: 45%" />
              <q-select v-model="AudioSettingNew.noiseSuppression" :options="OnOffOptions" emit-value map-options :label="$t('NoiseSuppression')" style="width: 45%" />
              <q-select v-model="AudioSettingNew.echoCancellation" :options="OnOffOptions" emit-value map-options :label="$t('EchoCancellation')" style="width: 45%" />
            </div>
            <br />
            <h7 class="text-body2">{{ $t("SampleRate") }}: {{ AudioSettingNew ? (AudioSettingNew.sampleRate / 1000).toFixed(0) : "0" }} kB/sec</h7> <br />
            <q-slider v-model="AudioSettingNew.sampleRate" :min="AudioCapabilities.sampleRate.min" :max="AudioCapabilities.sampleRate.max" label style="width: 90%" />
          </q-card-section>
          <q-separator />

          <q-card-actions class="tw-justify-end">
            <q-btn @click="AudioSettingDialog = false" class="tw-bg-red-300" icon="close">{{ $t("Cancel") }}</q-btn>
            <q-btn @click="AudioSettingChangeBtn" class="tw-bg-lime-300" icon="done">{{ $t("Save") }}</q-btn>
          </q-card-actions>
        </q-card>
      </q-dialog>
      <!--


        Dialog VideoSettingDialog
      -->
      <q-dialog v-model="VideoSettingDialog" no-backdrop-dismiss persistent class="tw-font-sans">
        <q-card style="width: 400px; height: 420px">
          <q-card-section>
            <div class="text-h6">Video Settings</div>
            <q-select v-model="VideoSettingNew.enabled" :options="OnOffOptions" emit-value map-options :label="$t('VideoEnabled')" style="width: 45%" />
            <br />
            <h7 class="text-body2">{{ $t("Resolution") }}: {{ VideoSettingNew ? `${VideoSettingNew.width}x${VideoSettingNew.height}` : "0x0" }}</h7
            ><br />
            <br /><h7 class="text-body2">{{ $t("Width") }}: {{ VideoSettingNew ? VideoSettingNew.width : "0" }}</h7
            ><br />
            <q-slider v-model="VideoSettingNew.width" :min="VideoCapabilities.width.min" :max="VideoCapabilities.width.max" label style="width: 80%" />
            <br /><h7 class="text-body2">{{ $t("Height") }}: {{ VideoSettingNew ? VideoSettingNew.height : "0" }}</h7
            ><br />
            <q-slider v-model="VideoSettingNew.height" :min="VideoCapabilities.height.min" :max="VideoCapabilities.height.max" label style="width: 80%" />
            <br /><h7 class="text-body2">{{ $t("SampleRate") }}: {{ VideoSettingNew ? VideoSettingNew.frameRate.toFixed(0) : "0" }}</h7
            ><br />
            <q-slider v-model="VideoSettingNew.frameRate" :min="VideoCapabilities.frameRate.min" :max="VideoCapabilities.frameRate.max" label style="width: 80%" />
          </q-card-section>
          <q-separator />

          <q-card-actions class="tw-justify-end">
            <q-btn @click="VideoSettingDialog = false" class="tw-bg-red-300" icon="close">{{ $t("Cancel") }}</q-btn>
            <q-btn @click="VideoSettingChangeBtn" class="tw-bg-lime-300" icon="done">{{ $t("Save") }}</q-btn>
          </q-card-actions>
        </q-card>
      </q-dialog>
    </div>
  </q-page>
</template>
<style></style>
