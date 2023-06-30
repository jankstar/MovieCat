<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "AudioSettingComp",
  props: ["langu", "selMode", "MainStream", "SupportedConstraints", "modelValue"],
  emits: ["update:modelValue"],

  components: {},

  data: () => {
    return {
      Setting: {
        muted: false,
        autoGainControl: false,
        noiseSuppression: false,
        echoCancellation: false,
        sampleRate: 48000,
      },
      AudioSettingNew: undefined,
      AudioCapabilities: undefined,
      AudioSettingDialog: false,
      OnOffOptions: [],
    };
  },
  mounted() {
    console.log(`AudioSettingComp mounted`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Setting = this.modelValue;
  },
  updated() {
    console.log(`AudioSettingComp updated`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Setting = this.modelValue;
  },
  //
  methods: {
    //
    AudioSettingBtn() {
      let [traks] = this.MainStream.getAudioTracks();
      if (traks) {
        console.log("AudioSettingBtn()");
        this.AudioCapabilities = traks.getCapabilities();
        this.AudioSettingNew = JSON.parse(JSON.stringify(this.Setting));
        this.AudioSettingDialog = true;
      }
    },
    //
    async AudioSettingChangeBtn() {
      let that = this;
      let [tracks] = this.MainStream.getAudioTracks();
      if (tracks) {
        console.log("AudioSettingChangeBtn()");

        tracks
          .applyConstraints(
            Object.assign(tracks.getSettings(), {
              autoGainControl: this.AudioSettingNew.autoGainControl,
              echoCancellation: this.AudioSettingNew.echoCancellation,
              noiseSuppression: this.AudioSettingNew.noiseSuppression,
              sampleRate: this.AudioSettingNew.sampleRate,
            })
          )
          .then(() => {
            tracks.enabled = !this.AudioSettingNew.muted;

            that.Setting = tracks.getSettings();
            that.Setting.muted = tracks.muted || !tracks.enabled;
            that.$emit("update:modelValue", that.Setting);
          })
          .catch((e: any) => {
            console.error(e);
            that.$q.notify({ type: "negative", message: `${e.name}: ${e.message}`, position: "center", timeout: 5000 });
          });
      }
      this.AudioSettingDialog = false;
    },
  },
});
</script>

<template>
  <div>
    <q-card-section>
      <div class="row tw-justify-between">
        <div class="text-subtitle1">Audio</div>
        <q-separator />
        <q-btn flat round sizes="sx" padding="none" icon="edit" @click="AudioSettingBtn" />
      </div>
      <div class="text-body2">{{ $t("AudioMuted") }}: {{ Setting && Setting.muted == false ? $t("Off") : $t("On") }}</div>
      <div v-if="SupportedConstraints.autoGainControl" class="text-body2">
        {{ $t("AutoGainControl") }}: {{ Setting && Setting.autoGainControl == false ? $t("Off") : $t("On") }}
      </div>
      <div v-if="SupportedConstraints.noiseSuppression" class="text-body2">
        {{ $t("NoiseSuppression") }}: {{ Setting && Setting.noiseSuppression == false ? $t("Off") : $t("On") }}
      </div>
      <div v-if="SupportedConstraints.echoCancellation" class="text-body2">
        {{ $t("EchoCancellation") }}: {{ Setting && Setting.echoCancellation == false ? $t("Off") : $t("On") }}
      </div>
      <div v-if="SupportedConstraints.sampleRate" class="text-body2">{{ $t("SampleRate") }}: {{ Setting ? (Setting.sampleRate / 1000).toFixed(0) : "0" }} kB/sec</div>
    </q-card-section>
    <q-separator />
  </div>
  <!-- Dialog AudioSettingDialog -->
  <q-dialog v-model="AudioSettingDialog" no-backdrop-dismiss persistent class="tw-font-sans">
    <q-card style="width: 400px; height: 330px">
      <q-card-section>
        <div class="text-h6">Audio Settings</div>
        <div class="row">
          <q-select v-model="AudioSettingNew.muted" :options="OnOffOptions" emit-value map-options :label="$t('AudioMuted')" style="width: 45%" />
          <q-select
            v-if="SupportedConstraints.autoGainControl"
            v-model="AudioSettingNew.autoGainControl"
            :options="OnOffOptions"
            emit-value
            map-options
            :label="$t('AutoGainControl')"
            style="width: 45%"
          />
          <q-select
            v-if="SupportedConstraints.noiseSuppression"
            v-model="AudioSettingNew.noiseSuppression"
            :options="OnOffOptions"
            emit-value
            map-options
            :label="$t('NoiseSuppression')"
            style="width: 45%"
          />
          <q-select
            v-if="SupportedConstraints.echoCancellation"
            v-model="AudioSettingNew.echoCancellation"
            :options="OnOffOptions"
            emit-value
            map-options
            :label="$t('EchoCancellation')"
            style="width: 45%"
          />
        </div>
        <br />
        <div v-if="SupportedConstraints.sampleRate" class="text-body2">
          {{ $t("SampleRate") }}: {{ AudioSettingNew ? (AudioSettingNew.sampleRate / 1000).toFixed(0) : "0" }} kB/sec
        </div>
        <br />
        <q-slider
          v-if="SupportedConstraints.sampleRate"
          v-model="AudioSettingNew.sampleRate"
          :min="AudioCapabilities.sampleRate.min"
          :max="AudioCapabilities.sampleRate.max"
          label
          style="width: 90%"
        />
      </q-card-section>
      <q-separator />

      <q-card-actions class="tw-justify-end">
        <q-btn @click="AudioSettingDialog = false" class="tw-bg-red-300" icon="close">{{ $t("Cancel") }}</q-btn>
        <q-btn @click="AudioSettingChangeBtn" class="tw-bg-lime-300" icon="done">{{ $t("Save") }}</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
