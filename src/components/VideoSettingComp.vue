<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "VideoSettingComp",
  props: ["langu", "selMode", "MainStream", "SupportedConstraints", "modelValue"],
  emits: ["update:modelValue"],

  components: {},

  data: () => {
    return {
      Setting: {
        width: 0,
        height: 0,
        frameRate: 0,
        enabled: false,
      },
      VideoSettingNew: undefined,
      VideoCapabilities: undefined,
      VideoSettingDialog: false,
      OnOffOptions: [],
    };
  },
  mounted() {
    console.log(`VideoSettingComp mounted`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Setting = this.modelValue;
  },
  updated() {
    console.log(`VideoSettingComp updated`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Setting = this.modelValue;
  },
  //
  methods: {
    //
    updateData() {
      this.$emit("update:modelValue", this.Setting);
    },
    //
    VideoSettingBtn() {
      let [tracks] = this.MainStream.getVideoTracks();
      if (tracks) {
        console.log("VideoSettingBtn()");
        this.VideoCapabilities = tracks.getCapabilities();
        this.VideoSettingNew = JSON.parse(JSON.stringify(this.Setting));
        this.VideoSettingDialog = true;
      }
    },
    //
    async VideoSettingChangeBtn() {
      let that = this;
      let [tracks] = this.MainStream.getVideoTracks();
      if (tracks) {
        console.log("VideoSettingChangeBtn()");
        tracks
          .applyConstraints(
            Object.assign(tracks.getSettings(), {
              height: this.VideoSettingNew.height,
              width: this.VideoSettingNew.width,
              frameRate: this.VideoSettingNew.frameRate,
            })
          )
          .then(() => {
            tracks.enabled = this.VideoSettingNew.enabled;

            that.Setting = tracks.getSettings();
            that.Setting.enabled = tracks.enabled;
            that.$emit("update:modelValue", that.Setting);
          })
          .catch((e: any) => {
            console.error(e);
            that.$q.notify({ type: "negative", message: `${e.name}: ${e.message}`, position: "center", timeout: 5000 });
          });

        this.VideoSettingDialog = false;
      }
    },
  },
});
</script>

<template>
  <div>
    <q-card-section>
      <div class="row tw-justify-between">
        <div class="text-subtitle1">Video</div>
        <q-separator />
        <q-btn flat round sizes="sx" padding="none" icon="edit" @click="VideoSettingBtn" />
      </div>
      <div class="text-body2">{{ $t("Resolution") }}: {{ Setting ? `${Setting.width}x${Setting.height}` : "0x0" }}</div>
      <div v-if="SupportedConstraints.frameRate" class="text-body2">{{ $t("Frame_rate") }}: {{ Setting ? Setting.frameRate.toFixed(2) : "0" }}</div>
    </q-card-section>
    <q-separator />
  </div>
  <!-- Dialog VideoSettingDialog  -->
  <q-dialog v-model="VideoSettingDialog" no-backdrop-dismiss persistent class="tw-font-sans">
    <q-card style="width: 400px; height: 400px">
      <q-card-section>
        <div class="text-h6">Video Settings</div>
        <q-select v-model="VideoSettingNew.enabled" :options="OnOffOptions" emit-value map-options :label="$t('VideoEnabled')" style="width: 45%" />
        <br />

        <div v-if="SupportedConstraints.width" class="text-body2">{{ $t("Resolution") }}: {{ VideoSettingNew ? `${VideoSettingNew.width}x${VideoSettingNew.height}` : "0x0" }}</div>
        <br v-if="SupportedConstraints.width" />

        <div v-if="SupportedConstraints.width" class="text-body2">{{ $t("Width") }}: {{ VideoSettingNew ? VideoSettingNew.width : "0" }}</div>
        <q-slider
          v-if="SupportedConstraints.width"
          v-model="VideoSettingNew.width"
          :min="VideoCapabilities.width.min"
          :max="VideoCapabilities.width.max"
          label
          style="width: 80%"
        />
        <div v-if="SupportedConstraints.height" class="text-body2">{{ $t("Height") }}: {{ VideoSettingNew ? VideoSettingNew.height : "0" }}</div>
        <q-slider
          v-if="SupportedConstraints.height"
          v-model="VideoSettingNew.height"
          :min="VideoCapabilities.height.min"
          :max="VideoCapabilities.height.max"
          label
          style="width: 80%"
        />
        <div class="text-body2">{{ $t("SampleRate") }}: {{ VideoSettingNew ? VideoSettingNew.frameRate.toFixed(0) : "0" }}</div>
        <q-slider
          v-if="SupportedConstraints.frameRate"
          v-model="VideoSettingNew.frameRate"
          :min="VideoCapabilities.frameRate.min"
          :max="VideoCapabilities.frameRate.max"
          label
          style="width: 80%"
        />
      </q-card-section>
      <q-separator />

      <q-card-actions class="tw-justify-end">
        <q-btn @click="VideoSettingDialog = false" class="tw-bg-red-300" icon="close">{{ $t("Cancel") }}</q-btn>
        <q-btn @click="VideoSettingChangeBtn" class="tw-bg-lime-300" icon="done">{{ $t("Save") }}</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
