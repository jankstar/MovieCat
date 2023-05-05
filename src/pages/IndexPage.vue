<script lang="ts">
import { defineComponent } from "vue";
import AudioConstraints from "components/AudioConstraints.vue";
import AudioSettingComp from "components/AudioSettingComp.vue";
import VideoSettingComp from "components/VideoSettingComp.vue";
import moviecat from "src/lib/moviecat.js";
import * as ebml from "ts-ebml";
import type { EBMLElementDetail } from "ts-ebml";

export default defineComponent({
  name: "IndexPage",
  props: ["langu"],

  components: { AudioConstraints, AudioSettingComp, VideoSettingComp },

  data: () => {
    return {
      $t: undefined,
      //MediaDevices
      SpinnerOn: false,
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

      AudioSetting: undefined,

      AudioDraw: false,
      CaptureAndMic: false,
      AudioConstraints: {
        deviceId: "default",
        autoGainControl: false,
        echoCancellation: false,
        noiseSuppression: false,
        sampleRate: 48000,
      },

      //Recorder
      Recorder: <MediaRecorder>undefined,
      RecorderCounter: 0,
      RecorderCounterJob: undefined,
      recorderOptions: {
        audioBitsPerSecond: 64000,
        videoBitsPerSecond: 2500000,
        mimeType: undefined,
      },
      audioBPSOptions: moviecat.ConstAudioBPSOptions,

      videoBPSOptions: moviecat.ConstVideoBPSOptions,
      mimeTypeOptions: [],

      RecorderSlices: 10, //in sec
      RecorderSlicesOptions: moviecat.ConstRecorderSlicesOptions,
      recorderAutoStop: 30, //in min
      RecorderState: "inactive",

      FileData: {
        Name: "test",
        StartTime: undefined,
        EndTime: undefined,
        Size: 0,
        BlobList: [],
        Duration: 0, //Attention only integer in sec !
        Position: 0, //Attention only integer in sec !
      },

      LocalStopFunc: undefined,

      FileApi: false,
      FileApiDir: "",
      FileApiFileHandler: undefined,
      FileApiFileEntries: [],
      FileApiPickerID: 0,

      VideoElement: undefined,
      PlayerMuted: true,
      calculateDuration: false,
      playOn: false,
      playbackRate: 1,
      playbackRateOptions: moviecat.ConstPlaybackRateOptions,
      TimeOutID: undefined,
      OnOffOptions: [],
    };
  },

  created() {
    console.log(`IndexPage created`);
    this.SupportedConstraints = navigator.mediaDevices.getSupportedConstraints();
    console.log(`SupportedConstraints: ${JSON.stringify(this.SupportedConstraints)}`);
    // @ts-ignore
    if (window.showDirectoryPicker) {
      console.log(`file api is true`);
      this.FileApi = true;
    }
  },

  async mounted() {
    try {
      console.log(`IndexPage mounted`);
      this.OnOffOptions = [
        { value: true, label: this.$t("On") },
        { value: false, label: this.$t("Off") },
      ];

      this.getMediaDevices();

      this.getMimeType();
      //recorder data from localStorage
      this.loadRecorderData();
    } catch (e) {
      console.error(e);
    }
  },

  updated() {
    console.log(`IndexPage updated`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
  },

  methods: {
    //
    saveRecorderData() {
      console.log(`saveRecorderData()`);

      localStorage.setItem("selMode", JSON.stringify(this.selMode));
      localStorage.setItem("recorderOptions", JSON.stringify(this.recorderOptions));
      localStorage.setItem("RecorderSlices", JSON.stringify(this.RecorderSlices));
      localStorage.setItem("recorderAutoStop", JSON.stringify(this.recorderAutoStop));
    },
    loadRecorderData() {
      console.log(`loadRecorderData()`);
      try {
        const JSONrselMode = localStorage.getItem("selMode");
        if (JSONrselMode) {
          let selMode = JSON.parse(JSONrselMode);
          if (selMode) {
            this.selMode = selMode;
          }
        }

        const JSONrecorderOptions = localStorage.getItem("recorderOptions");
        if (JSONrecorderOptions) {
          let recorderOptions = JSON.parse(JSONrecorderOptions);
          if (recorderOptions) {
            this.recorderOptions = recorderOptions;
          }
        }

        const JSONRecorderSlices = localStorage.getItem("RecorderSlices");
        if (JSONrecorderOptions) {
          let RecorderSlices = JSON.parse(JSONRecorderSlices);
          if (RecorderSlices) {
            this.RecorderSlices = RecorderSlices;
          }
        }

        const JSONrecorderAutoStop = localStorage.getItem("recorderAutoStop");
        if (JSONrecorderAutoStop) {
          let recorderAutoStop = JSON.parse(JSONrecorderAutoStop);
          if (recorderAutoStop) {
            this.recorderAutoStop = recorderAutoStop;
          }
        }
      } catch (e) {
        console.error(e);
      }
    },
    //
    async connectOnBtn() {
      let that = this;
      console.log(`connectOnBtn(e)`);

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
      let that = this;
      moviecat.ConstTrueMimeType.forEach((item: any) => {
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
              if (device.deviceId == "default" || device.deviceId == "") {
                item.label = "Video default";
                videoinputDefault = item;
              }
              that.videoInput.push(item);
            }

            if (device.kind == "audioinput") {
              let item = { label: device.label.substring(0, 40) || `Mic ${countMic++}`, value: device.deviceId || "default" };
              if (device.deviceId == "default" || device.deviceId == "") {
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
            that.videoInputValue = videoinputDefault.value;
          }

          if (audioinputDefault == undefined) {
            let item = { label: `Audio default`, value: "default" };
            that.audioInput.push(item);
            audioinputDefault = item;
          }
          if (!that.audioInputValue) {
            that.audioInputValue = audioinputDefault.value;
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
        this.CaptureAndMic = false;
        let configuration = { audio: <boolean | any>true, video: <boolean | any>true, systemAudio: "include", surfaceSwitching: "include", selfBrowserSurface: "exclude" };

        if (this.audioInputValue == "off") {
          configuration.audio = false;
        } else {
          configuration.audio = this.AudioConstraints;
          configuration.audio.deviceId = this.audioInputValue;
        }

        if (this.selMode == "capture") {
          //mode capture

          if (this.audioInputValue != "default" && this.audioInputValue != "off") {
            configuration.audio = false; //wenn capture und ein mic device, dann audio separat
            this.CaptureAndMic = true;
          }

          this.MainStream = await navigator.mediaDevices.getDisplayMedia(configuration);

          if (this.CaptureAndMic) {
            //noch ein mic audio kanal dazu
            let myAudio = await navigator.mediaDevices.getUserMedia({
              video: false,
              audio: {
                deviceId: this.audioInputValue.value,
              },
            });
            let tracks = myAudio.getAudioTracks();
            tracks.forEach((track) => {
              that.MainStream.addTrack(track);
            });
          }
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
            throw new Error(this.$t("InfoAudiOffVideoOff"));
          }
          this.MainStream = await navigator.mediaDevices.getUserMedia(configuration);
        } else {
          // error
          return;
        }

        await this.newMediaRecorder();

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

    async newMediaRecorder() {
      console.log(`newMediaRecorder()`);
      let that = this;
      try {
        if (!this.MainStream) {
          throw new Error(this.$t("InfoNoMediaStream"));
        }

        if (this.FileData.BlobList && this.FileData.BlobList.length > 0) {
          if ((await this.confirmDialog(this.$t("QMemory"), this.$t("Yes"), this.$t("No"))) == true) {
            this.clearBuffer();
            this.FileData.StartTime = Date.now();
          }
        }
        this.Recorder = new MediaRecorder(this.MainStream, this.recorderOptions);

        this.Recorder.ondataavailable = (e) => {
          console.log(`type: ${e.type} time: ${e.timecode} size: ${e.data.size}`);
          if (e.data && e.data.size > 0) {
            that.FileData.BlobList.push(e.data);
            that.FileData.EndTime = Date.now();
            that.FileData.Size += e.data.size;
            that.RecorderCounter = 0;
          }
        };

        //the events no longer function properly after a stop, so the recoder is completely restarted!
        this.Recorder.onstop = () => {
          console.log("onstop()");
          that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
          if (that.RecorderCounterJob) {
            that.stopTimer();
          }
          if (that.LocalStopFunc) {
            that.LocalStopFunc();
          }
        };

        this.Recorder.onerror = (e: any) => {
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
          if (that.RecorderSlices != 1) {
            //only if no 1 sec slices
            that.startTimer();
          }
        };
      } catch (err) {
        console.error(`${err.name}: ${err.message}`);
      }
    },

    //
    stopTimer() {
      console.log(`stopTimer()`);
      if (this.RecorderCounterJob) {
        clearTimeout(this.RecorderCounterJob);
        this.RecorderCounter = 0;
      }
    },

    startTimer() {
      console.log(`startTimer()`);
      let that = this;
      if (this.RecorderState != "inactive") {
        this.RecorderCounterJob = setTimeout(() => {
          that.RecorderCounter += 1;
          that.startTimer();
        }, 1000);
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
    async startRecorderBtn() {
      console.log(`startRecorderBtn()`);
      let that = this;

      await this.newMediaRecorder();

      if (this.Recorder && this.Recorder.state == "inactive") {
        if (this.RecorderSlices && this.RecorderSlices > 0) {
          this.Recorder.start(1000 * this.RecorderSlices); //1 * x sec Slices
        } else {
          this.Recorder.start(); //ohne slices
        }

        if (!this.FileData.StartTime) {
          this.FileData.StartTime = Date.now();
        }

        this.LocalStopFunc = undefined;
        this.TimeOutID = setTimeout(function () {
          that.LocalStopFunc = that.download;
          that.stopRecorder();
          //that.download();
          that.TimeOutID = undefined;
        }, 1000 * 60 * this.recorderAutoStop); //180 min = 3h
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
    async download() {
      console.log(`download: ${this.FileData.BlobList.length} blobs`);
      this.SpinnerOn = true;
      if (this.FileData.BlobList && this.FileData.BlobList.length > 0) {
        if (this.FileApiFileHandler && this.FileApiFileHandler.name) {
          //
          let lFielname = `${this.FileData.Name}.${this.recorderOptions.mimeType.split(";")[0].split("/")[1]}`;
          try {
            let lFileHandler = await this.FileApiFileHandler.getFileHandle(lFielname, { create: true });
            const lWritable = await lFileHandler.createWritable();

            const blob = new Blob(this.FileData.BlobList, { type: this.recorderOptions.mimeType });

            await lWritable.write(blob);
            await lWritable.close();
          } catch (err) {
            console.error(err);
            this.$q.notify({ type: "negative", message: `${err.name}: ${err.message}`, position: "center", timeout: 5000 });
          }
        } else {
          //standard download without file api
          const blob = new Blob(this.FileData.BlobList, { type: this.recorderOptions.mimeType });
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.style.display = "none";
          a.href = url;
          a.download = `${this.FileData.Name}.${this.recorderOptions.mimeType.split(";")[0].split("/")[1]}`;
          console.log(`filename ${a.download}`);
          document.body.appendChild(a);
          a.click();
          setTimeout(function () {
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
          }, 3000);
        }
      }
      this.SpinnerOn = false;
    },
    //
    async openDirBtn() {
      console.log(`openDirBtn()`);
      let that = this;
      try {
        // @ts-ignore
        if (!window.showDirectoryPicker) {
          throw Error("file api not supported");
        }
        if (!this.FileApiPickerID) {
          this.FileApiPickerID = (Math.random() * 100).toFixed(0);
        }
        if (!this.FileApiFileHandler) {
          // @ts-ignore
          this.FileApiFileHandler = await window.showDirectoryPicker({
            id: this.FileApiPickerID,
            mode: "readwrite",
            startIn: "downloads",
          });
        } else {
          // @ts-ignore
          this.FileApiFileHandler = await window.showDirectoryPicker({
            id: this.FileApiPickerID,
            mode: "readwrite",
          });
        }
        if (this.FileApiFileHandler) {
          console.log(`openDirBtn() ${this.FileApiFileHandler.name}`);
          this.FileApiPickerID = 0;
        }
        this.FileApiFileEntries = [];
        let lExtension = `.${this.recorderOptions.mimeType.split(";")[0].split("/")[1]}`;
        for await (const [key, value] of this.FileApiFileHandler.entries()) {
          if (value.kind == "file" && value.name.endsWith(lExtension)) {
            that.FileApiFileEntries.push({ key, value });
            console.log(`file: ${value.name}`);
          }
        }
      } catch (err) {
        console.error(err);
        this.$q.notify({ type: "negative", message: `${err.name}: ${err.message}`, position: "center", timeout: 5000 });
      }
    },
    //
    async loadFileBtn(iKey: string) {
      //file is loading in slices
      console.log(`loadFileBtn()`, iKey);
      if (this.FileData.BlobList.length > 0) {
        if ((await this.confirmDialog(this.$t("QMemory"), this.$t("Yes"), this.$t("Cancel"))) != true) {
          //exit and return
          return;
        }
      }
      this.SpinnerOn = true;
      try {
        let lFileHandler = await this.FileApiFileHandler.getFileHandle(iKey);
        const file = await lFileHandler.getFile();
        let lSize = file.size;
        this.clearBuffer();
        this.FileData.StartTime = Date.now();
        let lRecorderTime = this.FileData.StartTime;
        let lStart = 0;

        while (lSize > 0) {
          let content = undefined;
          let lContenSize = 0;
          let lSizeLoad = lSize;
          if (lSizeLoad > 2000000 * 1024) {
            //max memory 2 gByte
            lSizeLoad = 2000000 * 1024;
            lSize -= 2000000 * 1024;
            content = await file.slice(lStart, lStart + lSizeLoad); //await file.arrayBuffer();
            lStart += lSizeLoad;
          } else if (lStart != 0) {
            lSize = 0;
            content = await file.slice(lStart); //await file.arrayBuffer();
          } else {
            lSize = 0;
            content = await file.arrayBuffer();
          }

          this.FileData.BlobList.push(content);
          if (content.byteLength) {
            lContenSize = content.byteLength;
          } else {
            lContenSize = content.size;
          }
          console.log(`time: ${lRecorderTime} + ${Math.trunc(lContenSize / (0.25 * 1024))}`); //2,5 mByte /msec
          lRecorderTime += Math.trunc(lContenSize / (0.25 * 1024));
          this.FileData.EndTime = lRecorderTime;
          this.FileData.Size += lContenSize;
        }
        this.FileData.Name = iKey.substring(0, iKey.indexOf(`.${this.recorderOptions.mimeType.split(";")[0].split("/")[1]}`));

        if (this.FileData.BlobList.length > 0 && this.recorderOptions.mimeType.includes("webm")) {
          //if webm then decode
          try {
            this.calculateDuration = true;
            let that = this;
            const decoder = new ebml.Decoder();
            const reader = new ebml.Reader();
            //decode only the first blob due to memory limitation
            let blob = <Blob>{};
            if (this.FileData.BlobList.length == 1) {
              blob = new Blob(this.FileData.BlobList, { type: this.recorderOptions.mimeType });
            } else {
              blob = this.FileData.BlobList[0];
            }
            blob
              .arrayBuffer()
              .then((abuffer) => {
                const ebmlElms = <EBMLElementDetail[]>decoder.decode(abuffer);
                ebmlElms.forEach((elm) => {
                  reader.read(elm);
                });
                reader.stop();
                const sec = (reader.duration * reader.timecodeScale) / 1000 / 1000; // / 1000;
                console.log(`duration ${that.computeTime(sec)}`);
                if (that.FileData.StartTime == undefined) {
                  that.FileData.StartTime = 0;
                }
                that.FileData.EndTime = that.FileData.StartTime + sec * 1000;
                that.FileData.Duration = sec;
                that.calculateDuration = false;
              })
              .catch((e) => {
                console.error("blob larger than 2 GB: ", e);
                this.calculateDuration = false;
              });
          } catch (e) {
            console.error(e);
            this.calculateDuration = false;
          }
        }

        this.selMode = "player";
      } catch (err) {
        console.error(err);
        this.$q.notify({ type: "negative", message: `${err.name}: ${err.message}`, position: "center", timeout: 5000 });
      }
      this.SpinnerOn = false;
    },
    //
    clearBuffer() {
      this.FileData.BlobList = [];
      this.RecorderCounter = 0;
      this.FileData.Position = 0;
      this.FileData.StartTime = undefined;
      this.FileData.EndTime = undefined;
      this.FileData.Size = 0;
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

      this.VideoElement = <HTMLVideoElement>document.getElementById("id_video_player");
      if (this.playOn == false && this.VideoElement && this.FileData.BlobList && this.FileData.BlobList.length > 0 && this.recorderOptions.mimeType) {
        this.VideoElement.pause(); //player stoppen
        //
        this.VideoElement.onloadedmetadata = () => {
          console.log("onloadedmetadata()");
          let sec = Number.isFinite(that.VideoElement.duration) ? Math.trunc(that.VideoElement.duration) : that.FileData.Duration;
          if (sec && that.FileData.Duration < sec) {
            that.FileData.Duration = sec;
            if (that.FileData.StartTime == undefined) {
              that.FileData.StartTime = 0;
            }
            that.FileData.EndTime = that.FileData.StartTime + Math.trunc(that.FileData.Duration * 1000);
          }
          console.log("Duration change", that.FileData.Duration);
        };
        this.VideoElement.ondurationchange = () => {
          console.log("ondurationchange()");
          let sec = Number.isFinite(that.VideoElement.duration) ? Math.trunc(that.VideoElement.duration) : that.FileData.Duration;
          if (sec && that.FileData.Duration < sec) {
            that.FileData.Duration = sec;
            if (that.FileData.StartTime == undefined) {
              that.FileData.StartTime = 0;
            }
            that.FileData.EndTime = that.FileData.StartTime + Math.trunc(that.FileData.Duration * 1000);
          }
          console.log("Duration change", that.FileData.Duration);
        };
        this.VideoElement.ontimeupdate = () => {
          console.log(`currentTime ${that.VideoElement.currentTime.toFixed(2)} sec`);
          that.FileData.Position =
            that.VideoElement.currentTime && Number.isFinite(that.VideoElement.currentTime) ? Math.trunc(that.VideoElement.currentTime) : that.FileData.Position;
          if (that.FileData.Position && that.FileData.Duration < that.FileData.Position) {
            that.FileData.Duration = that.FileData.Position + 1;
            if (that.FileData.StartTime == undefined) {
              that.FileData.StartTime = 0;
            }
            that.FileData.EndTime = that.FileData.StartTime + Math.trunc(that.FileData.Duration * 1000);
          }
        };
        this.VideoElement.onended = () => {
          that.playOn = false;
        };
        this.VideoElement.onplay = () => {
          that.playOn = true;
        };
        //
        //        var url = window.URL.createObjectURL(this.FileData.BlobList[this.FileData.Position]);
        let blob = new Blob(this.FileData.BlobList, { type: this.recorderOptions.mimeType });
        let url = window.URL.createObjectURL(blob);
        this.VideoElement.src = url;
        this.VideoElement.currentTime = this.FileData.Position;
        this.VideoElement.playbackRate = this.playbackRate;
        this.VideoElement.play();
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
    changePosBtn(iSec) {
      console.log(`muteOffBtn(${iSec})`);
    },
    changePlaybackRate(value: any) {
      console.log(`muteOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        myVideoTag.playbackRate = value;
      }
    },
    //
    computeTime(iDuration: number): string {
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
    //
    confirmDialog(iMessage: string, iYes: string, iNo: string) {
      return new Promise((resolve) =>
        this.$q
          .dialog({
            title: this.$t("Confirm"),
            message: iMessage,
            ok: {
              label: iYes,
              push: iYes != "",
              class: "tw-bg-lime-300",
            },
            cancel: {
              label: iNo,
              push: iNo != "",
              class: "tw-bg-red-300",
            },
            persistent: true,
            class: "tw-font-sans",
          })
          .onOk(() => resolve(true))
          .onCancel(() => resolve(false))
      );
    },
  },
});
</script>

<template>
  <q-page class="q-pa-md">
    <div class="tw-flex tw-items-start tw-flex-wrap tw-gap-4">
      <!-- input and mode -->
      <q-card style="max-width: 300px">
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
            @update:model-value="saveRecorderData"
          />
          <br />
        </q-card-section>
        <q-separator v-if="selMode == 'capture' || selMode == 'camera'" />
        <q-card-section v-if="selMode != 'player'">
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
        <q-separator v-if="selMode != 'player'" />

        <AudioConstraints v-if="selMode != 'player' && !connectOn" v-model="AudioConstraints" :langu="langu" :SupportedConstraints="SupportedConstraints" />

        <AudioSettingComp
          v-if="AudioSetting && connectOn"
          v-model="AudioSetting"
          :langu="langu"
          :selMode="selMode"
          :MainStream="MainStream"
          :SupportedConstraints="SupportedConstraints"
        />

        <VideoSettingComp
          v-if="VideoSetting && connectOn"
          v-model="VideoSetting"
          :langu="langu"
          :selMode="selMode"
          :MainStream="MainStream"
          :SupportedConstraints="SupportedConstraints"
        />

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
          <h7 class="text-subtitle1">Audio ({{ AudioSetting && !AudioSetting.muted ? $t("On") : $t("Off") }}) </h7>

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
          <h7 class="text-subtitle1">Video Player </h7><br />
          <video id="id_video_player" playsinline muted style="background-color: black"></video>
        </q-card-section>
        <q-card-section>
          <h7 v-if="FileData.StartTime != FileData.EndTime" class="text-body2">{{ $t("Time") }}: {{ computeTime(FileData.Position) }}</h7>
          <br />
          <q-slider v-model="FileData.Position" :min="0" :max="FileData.Duration" label style="width: 90%" />
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

        <q-card-actions class="tw-grid tw-grid-rows-2 tw-gap-4 tw-justify-center">
          <div>
            <q-btn icon="skip_previous" @click="changePosBtn(0)"
              ><q-tooltip>{{ $t("Start") }}</q-tooltip></q-btn
            >
            <q-btn icon="replay_10" @click="changePosBtn(FileData.Position - 10)"
              ><q-tooltip>{{ $t("10_sec_replay") }}</q-tooltip></q-btn
            >
            <q-btn
              :label="!playOn ? $t('Play') : $t('Stop')"
              @click="!playOn ? playOnBtn() : playOffBtn()"
              :icon="!playOn ? 'play_circle' : 'stop_circle'"
              :class="!playOn ? 'tw-bg-lime-300' : 'tw-bg-red-300'"
              :disable="!FileData.BlobList || FileData.BlobList.length == 0"
            />
            <q-btn icon="forward_10" @click="changePosBtn(FileData.Position + 10)"
              ><q-tooltip>{{ $t("10_sec_forward") }}</q-tooltip></q-btn
            >
            <q-btn icon="skip_next" @click="changePosBtn(FileData.Duration)"
              ><q-tooltip>{{ $t("End") }}</q-tooltip></q-btn
            >
          </div>
          <!-- -->
          <div>
            <q-btn @click="$q.fullscreen.toggle(VideoElement)" :icon="$q.fullscreen.isActive ? 'fullscreen_exit' : 'fullscreen'" />
            <q-btn :label="!PlayerMuted ? $t('Mute') : $t('Mute_off')" @click="!PlayerMuted ? muteOnBtn() : muteOffBtn()" :icon="!PlayerMuted ? 'volume_up' : 'volume_off'" />
          </div>
        </q-card-actions>
      </q-card>

      <!-- recorder -->
      <q-card v-if="selMode == 'capture' || selMode == 'camera' || selMode == 'player' || FileData.Size > 0" style="max-width: 300px">
        <q-card-section>
          <h7 class="text-subtitle1">{{ $t("Recorder_state") }} {{ RecorderState }}</h7> <br />

          <h7 class="text-body2">{{ `${$t("Size")} ${(FileData.Size / 1000000).toFixed(2)}` }} mByte </h7>
          <h7 style="font-size: 10px"> ({{ $t("InfoRefreshPerSlices") }})</h7>
          <br />
          <div class="row">
            <h7 class="text-body2"
              >{{ $t("Time") }}:
              {{ FileData.StartTime && FileData.EndTime ? computeTime((FileData.EndTime - FileData.StartTime) / 1000 + RecorderCounter) : computeTime(RecorderCounter) }}
            </h7>
            <q-spinner-hourglass v-if="calculateDuration" color="primary" size="1em" />
            <h7 v-if="calculateDuration" style="font-size: 10px">calculate duration</h7>
          </div>

          <q-select
            v-model="recorderOptions.mimeType"
            :options="mimeTypeOptions"
            emit-value
            map-options
            :label="$t('Recorder_mimetype')"
            :disable="mimeTypeOptions.length == 0 || FileData.Size > 0 || RecorderState != 'inactive'"
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
                v-model="RecorderSlices"
                :options="RecorderSlicesOptions"
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
                :rules="[(val) => !!val || this.$t('Required'), (val) => val > 0 || this.$t('InfoMin1'), (val) => val < 181 || this.$t('InfoMax180')]"
                @update:model-value="saveRecorderData"
              />
            </div>
          </div>
          <q-input v-model="FileData.Name" :label="$t('Filename')" />
        </q-card-section>
        <q-separator />

        <q-card-actions class="tw-justify-end tw-gap-2">
          <q-btn
            v-if="RecorderState == 'inactive' && (selMode == 'capture' || selMode == 'camera')"
            label="Rec on"
            icon="videocam"
            @click="startRecorderBtn"
            :disable="RecorderState != 'inactive' || !connectOn || this.recorderAutoStop > 180"
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
          <q-btn label="download" @click="download" :disable="FileData.Size == 0 || RecorderState != 'inactive'" icon="download" />
          <q-btn label="clear" @click="clearBuffer" :disable="FileData.Size == 0 || RecorderState != 'inactive'" icon="delete_forever" />
        </q-card-actions>
      </q-card>

      <q-card v-if="FileApi" class="" style="width: 600px">
        <q-card-section v-if="FileApi">
          <div class="row tw-justify-between">
            <h7 class="text-subtitle1"> Files {{ `*/${FileApiFileHandler && FileApiFileHandler.name ? FileApiFileHandler.name : "?"}: ${FileApiFileEntries.length || 0}` }}</h7>
            <q-btn label="DIR" icon="folder" @click="openDirBtn" />
          </div>
          <br />

          <q-virtual-scroll style="max-height: 280px" :items="FileApiFileEntries" separator v-slot="{ item, index }">
            <q-item :key="index" dense>
              <q-item-section>
                <q-item-label>
                  <q-btn icon="upload_file" sizes="sx" padding="none" @click="loadFileBtn(`${item.key}`)" :disable="RecorderState != 'inactive' || connectOn" />
                  {{ item.key }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-virtual-scroll>
        </q-card-section>
        <q-separator v-if="FileApi" />
      </q-card>
      <q-dialog v-model="SpinnerOn" no-backdrop-dismiss persistent class="tw-font-sans">
        <div v-if="SpinnerOn" class="q-gutter-md justify-center"><q-spinner-hourglass color="light-green" size="xl" /></div>
      </q-dialog>
    </div>
  </q-page>
</template>
<style></style>
