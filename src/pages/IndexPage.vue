<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "IndexPage",
  props: ["langu"],

  components: {},

  data: () => {
    return {
      //MediaDevices
      captureOn: false,
      selCapture: true,
      videoInput: [],
      videoInputValue: "default",
      audioInput: [],
      audioInputValue: "default",

      //Stream
      MainStream: <MediaStream>undefined,
      VideoSetting: undefined,
      AudioSetting: undefined,
      AudioDraw: false,

      //Recorder
      Recorder: <MediaRecorder>undefined,
      // eslint-disable-next-line no-undef
      RecorderOption: <MediaRecorderOptions>{
        mimeType: "",
      },
      RecoderBlobList: [],
      RecoderInfoList: [],
      RecorderSize: 0,
      RecorderState: "inactive",
      TimeOutID: undefined,
      fileName: "test",
    };
  },

  created() {
    console.log(`created`);
  },

  async mounted() {
    console.log(`mounted`);
    await this.getMediaDevices();
  },

  methods: {
    //
    async startBtn() {
      let that = this;
      console.log(`startBtn(e)`);
      this.RecoderBlobList = [];
      this.RecoderInfoList = [];
      this.RecorderSize = 0;
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video");
      if (myVideoTag && !this.captureOn) {
        await this.startDisplayMedia();
        if (this.MainStream && this.MainStream.active) {
          this.MainStream.oninactive = () => {
            that.stopBtn();
          };
          myVideoTag.srcObject = this.MainStream;
          this.captureOn = true;
        }
      }
    },
    //
    // eslint-disable-next-line no-unused-vars
    stopBtn() {
      console.log(`stopBtn(e)`);

      this.MainStream.oninactive = undefined;
      this.stopDisplayMedia();

      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video");
      if (myVideoTag) {
        myVideoTag.srcObject = undefined;
      }
      this.captureOn = false;
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
          //$q.notify({ type: "negative", message: `${err.name}: ${err.message}`, position: "center", timeout: 5000 });
        });
    },
    //
    async startDisplayMedia() {
      console.log("startDisplayMedia()");
      let that = this;
      try {
        let configuration = { audio: <boolean | any>true, video: <boolean | any>true, systemAudio: "include", surfaceSwitching: "include", selfBrowserSurface: "exclude" };
        if (this.audioInputValue && this.audioInputValue != "default") {
          configuration.audio = {
            deviceId: this.audioInputValue || "default",
          };
        }
        if (this.audioInputValue && this.audioInputValue == "off") {
          configuration.audio = false;
        }
        if (this.selCapture) {
          //device capture
          this.MainStream = await navigator.mediaDevices.getDisplayMedia(configuration);
        } else {
          //device camera
          if (this.videoInputValue == "off") {
            configuration.video = false;
          }
          configuration.video = {
            deviceId: this.videoInputValue || "default",
          };
          this.MainStream = await navigator.mediaDevices.getUserMedia(configuration);
        }

        this.newMediaRecorder();
        let lTracks = this.MainStream.getVideoTracks();
        if (lTracks && lTracks.length > 0 && lTracks[0]) {
          this.VideoSetting = lTracks[0].getSettings();
        }

        lTracks = this.MainStream.getAudioTracks();
        if (lTracks && lTracks.length == 0 && configuration.audio.deviceId != "default" && configuration.audio.deviceId != "off") {
          //
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
          this.AudioSetting.muted = lTracks[0].muted;

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
      if (!this.MainStream) {
        throw new Error("no media stream");
      }
      this.Recorder = new MediaRecorder(this.MainStream);
      this.RecoderBlobList = [];
      this.RecoderInfoList = [];
      this.RecorderSize = 0;
      let that = this;

      this.Recorder.ondataavailable = (e) => {
        console.log(`type: ${e.type} time: ${e.timecode} size: ${e.data.size}`);
        if (e.data && e.data.size > 0) {
          that.RecoderBlobList.push(e.data);
          that.RecoderInfoList.push({ type: e.type, time: e.timecode, size: e.data.size });
          that.RecorderSize += e.data.size;
        }
      };

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
      this.Recorder.onstart = (e) => {
        console.log("onstart()");

        that.RecorderState = that.Recorder && that.Recorder.state ? that.Recorder.state : "inactive";
      };
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
    startRecorderBtn() {
      console.log(`startRecorderBtn()`);
      let that = this;

      this.newMediaRecorder();

      if (this.Recorder && this.Recorder.state == "inactive") {
        this.Recorder.start(1000); //1 sec Blöcke

        this.TimeOutID = setTimeout(function () {
          that.stopRecorder();
          that.TimeOutID = undefined;
        }, 1000 * 60 * 6 * 25); //2,5 Stunden
      }
    },
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
      if (!(this.Recorder && this.Recorder.state == "inactive")) {
        this.Recorder.requestData();
        this.Recorder.stop();
      }
      this.RecorderState = this.Recorder && this.Recorder.state ? this.Recorder.state : "inactive";
    },
    //
    download() {
      console.log(`download: ${this.RecoderBlobList.length} blobs`);
      if (this.RecoderBlobList && this.RecoderBlobList.length > 0) {
        var blob = new Blob(this.RecoderBlobList, { type: "video/webm" });
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement("a");
        a.style.display = "none";
        a.href = url;
        a.download = `${this.fileName}.webm`;
        document.body.appendChild(a);
        a.click();
        setTimeout(function () {
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);
        }, 3000);
      }
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
        }
      }

      draw();
    },
    //
  },
});
</script>
<template>
  <q-page class="q-pa-md">
    <div class="tw-flex tw-row tw-gap-4">
      <!-- capture -->
      <q-card style="max-width: 300px">
        <q-card-section>
          <h7 class="text-body2">{{ $t("Capture_state") }} {{ captureOn ? $t("On") : $t("Off") }}</h7> <br />
          <q-select
            v-model="selCapture"
            :options="[
              { label: $t('Capture'), value: true },
              { label: $t('Camera'), value: false },
            ]"
            emit-value
            map-options
            :label="$t('Device')"
            :disable="captureOn"
          />

          <q-select v-if="!selCapture" v-model="videoInputValue" :options="videoInput" emit-value map-options :label="$t('Video_Input')" :disable="captureOn" />
          <q-select v-model="audioInputValue" :options="audioInput" emit-value map-options :label="$t('Audio_Input')" :disable="captureOn" />
        </q-card-section>

        <q-card-section>
          <h7 class="text-subtitle1">Audio ({{ AudioSetting ? $t("On") : $t("Off") }})</h7> <br />
          <h7 class="text-body2">{{ $t("AudioMuted") }}: {{ AudioSetting && AudioSetting.muted == false ? $t("Off") : $t("On") }}</h7> <br />

          <div class="peer" id="peer-audio">
            <div class="stat-value">
              <canvas class="peer_canvas rounded" id="canvas-audio" style="background-color: black; width: 268px"> </canvas
              ><!-- das canvas wird über Programm gefüllt-->
            </div>
          </div>

          <h7 class="text-subtitle1">Video ({{ VideoSetting ? $t("On") : $t("Off") }})</h7> <br />
          <h7 class="text-body2">{{ $t("Resolution") }}: {{ VideoSetting ? `${VideoSetting.width}x${VideoSetting.height}` : "0x0" }}</h7> <br />
          <h7 class="text-body2">{{ $t("Frame_rate") }}: {{ VideoSetting ? VideoSetting.frameRate.toFixed(2) : "0" }}</h7> <br />

          <video id="id_video" playsinline autoplay muted style="background-color: black"></video>
        </q-card-section>
        <q-separator />

        <q-card-actions>
          <q-btn v-if="!captureOn" :label="$t('Start')" @click="startBtn" />
          <q-btn v-if="captureOn" :label="$t('Stop')" @click="stopBtn" />
        </q-card-actions>
      </q-card>

      <!-- recorder -->
      <q-card v-if="captureOn" style="max-width: 300px">
        <q-card-section>
          <h7 class="text-body2">{{ $t("Recorder_state") }} {{ RecorderState }}</h7> <br />
          <h7 class="text-body2">{{ `${$t("Size")} ${(RecorderSize / 1000000).toFixed(2)}` }} mByte</h7> <br />
          <q-input v-model="fileName" :label="$t('Filename')" />
        </q-card-section>
        <q-card-actions>
          <q-btn label="Rec on" @click="startRecorderBtn" :disable="RecorderState != 'inactive'" />
          <q-btn label="Rec off" @click="stopRecorderBtn" :disable="RecorderState == 'inactive'" />
          <q-btn label="download" @click="download" :disable="RecorderSize == 0" />
        </q-card-actions>
      </q-card>
    </div>
  </q-page>
</template>
<style></style>
