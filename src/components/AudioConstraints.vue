<script lang="ts">
import { defineComponent } from "vue";
import moviecat from "src/lib/moviecat.js";

export default defineComponent({
  name: "AudioContraints",
  props: ["langu", "SupportedConstraints", "modelValue"],
  emits: ["update:modelValue"],

  components: {},

  data: () => {
    return {
      Constraints: {
        deviceId: "default",
        autoGainControl: false,
        echoCancellation: false,
        noiseSuppression: false,
        sampleRate: 48000,
      },
      OnOffOptions: [],
      audioBPSOptions: moviecat.ConstAudioBPSOptions,
    };
  },
  mounted() {
    console.log(`AudioContraints mounted`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Constraints = this.modelValue;
  },
  updated() {
    console.log(`AudioContraints updated`);
    this.OnOffOptions = [
      { value: true, label: this.$t("On") },
      { value: false, label: this.$t("Off") },
    ];
    this.Constraints = this.modelValue;
  },
  methods: {
    updateData() {
      this.$emit("update:modelValue", this.Constraints);
    },
  },
});
</script>

<template>
  <div>
    <q-card-section style="width: 300px">
      <h7 class="text-subtitle1">Audio</h7> <br />
      <div class="row">
        <q-select
          v-if="SupportedConstraints.sampleRate"
          v-model="Constraints.sampleRate"
          :options="audioBPSOptions"
          emit-value
          map-options
          :label="$t('Audio_BPS')"
          style="width: 50%"
          @update:model-value="updateData"
        />
        <q-select
          v-if="SupportedConstraints.autoGainControl"
          v-model="Constraints.autoGainControl"
          :options="OnOffOptions"
          emit-value
          map-options
          :label="$t('AutoGainControl')"
          @update:model-value="updateData"
          style="width: 50%"
        />
        <q-select
          v-if="SupportedConstraints.noiseSuppression"
          v-model="Constraints.noiseSuppression"
          :options="OnOffOptions"
          emit-value
          map-options
          :label="$t('NoiseSuppression')"
          @update:model-value="updateData"
          style="width: 50%"
        />
        <q-select
          v-if="SupportedConstraints.echoCancellation"
          v-model="Constraints.echoCancellation"
          :options="OnOffOptions"
          emit-value
          map-options
          :label="$t('EchoCancellation')"
          @update:model-value="updateData"
          style="width: 50%"
        />
      </div>
    </q-card-section>
    <q-separator />
  </div>
</template>
