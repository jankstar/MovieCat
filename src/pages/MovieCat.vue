<script lang="ts">
import { defineComponent } from "vue";
import { useQuasar } from "quasar";
// @ts-ignore
import { version, dependencies } from "../../package.json";

export default defineComponent({
  name: "MovieCat",
  props: ["langu"],

  components: {},

  mounted() {
    this.dependenciesList = [];
    for (let p in this.dependencies) {
      this.dependenciesList.push(`${p} : ${this.dependencies[p]}`);
    }
  },

  data: () => {
    const $q = useQuasar();

    return {
      $q,
      $t: undefined,
      // eslint-disable-next-line no-undef
      version: version,
      dependencies: dependencies,
      dependenciesList: [],
    };
  },
});
</script>
<template>
  <q-page class="q-pa-md">
    <div class="tw-flex tw-row tw-gap-4">
      <!-- capture -->
      <q-card>
        <q-card-section>
          <h5><b>MovieCat</b></h5>
          <div v-if="langu == 'de-DE'">
            <br />
            Version: {{ version }} <br />
            Beschreibung: <br />
            Aufzeichnung eines Films von einer Kamera oder Bildschirm und sichern als Datei.
            <br />
          </div>
          <div v-else>
            <br />
            Version: {{ version }}<br />
            Description: <br />
            Record a movie from a camera or screen and save it as a file.
            <br />
          </div>
          <div>
            <br />
            <a href="https://quasar.dev">Quasar: {{ $q.version }} </a><br />
            <a href="https://https://vitejs.dev">Vite: 2.9.15</a><br />
            <a href="https://tailwindcss.com">tailwindcss: 3.2.4</a> <br />
            <a href="https://github.com/webrtcHacks/adapter">adapter.js: 8.0.0</a><br />
            Platform: {{ $q.platform.is.platform }} <br />
            Browser target: es2019, edge88, firefox78, chrome87, safari13.1<br />
          </div>
          <div>
            <br />
            Icons: material-icons <br />
            fontFamily: ChakraPetch, Merriweather<br />
            This Font Software is licensed under the SIL Open Font License, Version 1.1. <br />
            <a href="https://fonts.google.com/" target="_blank" class="text-overline">fonts.google.com</a><br />
          </div>
          <br />

          <div>Dependencies</div>
          <div v-for="item in dependenciesList" v-bind:key="item">{{ item }}</div>
          <br /><br />
          GitHub: <a href="https://github.com/jankstar/MovieCat" target="_blank"
            class="text-overline">https://github.com/jankstar/MovieCat</a><br />
          <a href="https://github.com/jankstar/MovieCat/blob/main/LICENSE" target="_blank">MIT License Copyright (c)
            2024 jankstar</a><br />
        </q-card-section>
        <q-separator />

        <q-card-actions>
          <q-btn :label="`${$t('Home')}`" @click="$router.push('/')" />
        </q-card-actions>
      </q-card>
    </div>
  </q-page>
</template>
<style></style>
