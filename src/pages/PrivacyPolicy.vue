<script lang="ts">
import { defineComponent } from "vue";
import { useQuasar } from "quasar";

export default defineComponent({
  name: "PrivacyPolicy",
  props: ["langu"],

  components: {},

  data: () => {
    const $q = useQuasar();
    return {
      $q,
      $t: undefined,
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
          <div v-if="langu == 'de-DE'">
            <h5><b>Datenschutz</b></h5>

            <br />
            A) Keine persönlichen Daten <br />
            Diese Software speichert keine persönlichen Daten vom Anwender. Es werden keine User oder Account Daten abgefragt oder gespeichert. Es erfolgt keine Speicherung von
            Zugriffen auf den Server, keine Performance Trace oder Werbungen. Der Quellcode des Programms ist öffentlich zugänglich unter GitHub:
            <a href="https://github.com/jankstar/MovieCat" target="_blank" class="text-overline">https://github.com/jankstar/MovieCat</a>. <br />
            B) Lokalestore <br />
            Folgende technische Daten werden im localstore gespeichert und dienen der Voreinstellung: Modus, Audio-Bits per Secunde, video Bits Per Sekunde, mime Type, recorder
            Slices, recorder Auto Stop Time, Sprache, dark-flag.
          </div>
          <div v-else>
            <h5><b>Privacy Policy</b></h5>

            <br />
            A) No personal data <br />
            This software does not store any personal data of the user. No user or account data is requested or stored. There is no storage of accesses to the server, no
            performance trace or advertisements. The source code of the programme is publicly accessible: GitHub:
            <a href="https://github.com/jankstar/MovieCat" target="_blank" class="text-overline">https://github.com/jankstar/MovieCat</a>. <br />
            B) Localstore <br />
            The following technical data are stored in the localstore and serve as default settings: mode, audio bits per second, video bits per second, mime type, recorder slices,
            recorder auto stop time, language, dark-flag.
          </div>
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
