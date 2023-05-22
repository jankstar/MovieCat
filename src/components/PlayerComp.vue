<script lang="ts">
import moviecat from "src/lib/moviecat";

import { defineComponent } from "vue";

export default defineComponent({
  name: "PlayerComp",
  props: ["langu", "modelValue", "connectOn", "selMode"],
  emits: ["update:modelValue"],

  components: {},

  data: () => {
    return {
      FileData: {
        Name: "test",
        Extension: "",
        StartTime: undefined,
        EndTime: undefined,
        Size: 0,
        BlobList: [],
        Duration: 0, //Attention only integer in sec !
        Position: 0, //Attention only integer in sec !
      },
      moviecat: moviecat,
      VideoElement: undefined,
      playbackRate: 1,
      playbackRateOptions: moviecat.ConstPlaybackRateOptions,
      PlayerMuted: true,
      PlayerWaiting: false,
      playOn: false,
    };
  },
  mounted() {
    console.log(`PlayerComp mounted`);

    this.FileData = this.modelValue;
  },
  updated() {
    console.log(`PlayerComp updated`);

    this.FileData = this.modelValue;
  },
  //
  methods: {
    updateData() {
      this.$emit("update:modelValue", this.FileData);
    },
    changePlaybackRate(value: any) {
      console.log(`muteOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        myVideoTag.playbackRate = value;
      }
    },
    async changePosBtn(iSec) {
      console.log(`muteOffBtn(${iSec})`);
      if (iSec < 0) {
        this.FileData.Position = 0;
      } else if (iSec > this.FileData.Duration) {
        this.FileData.Position = this.FileData.Duration;
      } else {
        this.FileData.Position = iSec;
      }

      if (this.playOn) {
        this.VideoElement = <HTMLVideoElement>document.getElementById("id_video_player");
        if (this.VideoElement) {
          this.VideoElement.pause();
          this.VideoElement.currentTime = this.FileData.Position;
          try {
            await this.VideoElement.play();
          } catch (e) {
            console.error(e);
          }
        }
      }
    },
    async playOnBtn() {
      console.log(`playOnBtn()`);
      let that = this;

      this.VideoElement = <HTMLVideoElement>document.getElementById("id_video_player");
      if (this.playOn == false && this.VideoElement && this.FileData.BlobList && this.FileData.BlobList.length > 0 && this.FileData.Extension) {
        //this.recorderOptions.mimeType) {
        this.VideoElement.pause(); //player stop
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
          that.updateData();
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
          that.updateData();
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
          that.updateData();
        };
        this.VideoElement.onplay = () => {
          that.playOn = true;
          that.updateData();
        };
        //
        //        var url = window.URL.createObjectURL(this.FileData.BlobList[this.FileData.Position]);
        let blob = new Blob(this.FileData.BlobList, { type: `video/${this.FileData.Extension}` }); //this.recorderOptions.mimeType });
        let url = window.URL.createObjectURL(blob);
        this.VideoElement.src = url;
        this.VideoElement.currentTime = this.FileData.Position;
        this.VideoElement.playbackRate = this.playbackRate;
        try {
          this.PlayerWaiting = true;
          await this.VideoElement.play();
          this.PlayerWaiting = false;
        } catch (e) {
          console.error(e);
          this.$q.notify({ type: "negative", message: `${e.name}: ${e.message}`, position: "center", timeout: 5000 });
          this.PlayerWaiting = false;
          this.playOn = false;
        }
      }
    },
    playOffBtn() {
      console.log(`playOffBtn()`);
      let myVideoTag = <HTMLVideoElement>document.getElementById("id_video_player");
      if (myVideoTag) {
        myVideoTag.pause(); //player stop
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
  },
});
</script>

<template>
  <!-- play video from file/recorder -->
  <q-card style="min-width: 300px; max-width: 60%" class="tw-grow">
    <q-card-section>
      <h7 class="text-subtitle1">Video Player </h7><br />
      <video id="id_video_player" playsinline muted style="background-color: black"></video>
    </q-card-section>
    <q-card-section>
      <h7 v-if="FileData.StartTime != FileData.EndTime" class="text-body2">{{ $t("Time") }}: {{ moviecat.computeTime(FileData.Position) }}</h7>
      <br />
      <q-slider v-model="FileData.Position" :min="0" :max="FileData.Duration" label />
      <br />
      <div class="row tw-gap-4">
        <q-select
          v-model="playbackRate"
          :options="playbackRateOptions"
          emit-value
          map-options
          :label="$t('PlaybackRate')"
          style="width: 50%; max-width: 120px"
          @update:model-value="changePlaybackRate"
        />
        <q-spinner-hourglass v-if="PlayerWaiting" color="primary" size="1em" />
        <h7 v-if="PlayerWaiting" style="font-size: 10px">start playing ...</h7>
      </div>
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
</template>
