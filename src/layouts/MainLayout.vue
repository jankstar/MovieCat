<script>
import { defineComponent } from "vue";

export default defineComponent({
  name: "MainLayout",

  components: {},
  data: () => {
    return {
      localeOptions: [
        { value: "en-US", label: "English" },
        { value: "de-DE", label: "Deutsch" },
      ],
    };
  },
  mounted() {
    const JSONi18n_locale = localStorage.getItem("i18n_locale");
    if (JSONi18n_locale) {
      let i18n_locale = JSON.parse(JSONi18n_locale);
      if (i18n_locale) {
        this.$i18n.locale = i18n_locale;
      }
    }
    const JSONdark = localStorage.getItem("dark");
    if (JSONdark) {
      let dark = JSON.parse(JSONdark);
      if (dark) {
        if (this.$q.dark.isActive != dark) {
          this.$q.dark.toggle();
        }
      }
    }
  },
  computed: {
    isValid() {
      return true;
    },
  },
  methods: {
    info() {
      console.log(`info()`);
      if (window.location.hash.includes("moviecat")) {
        this.$router.push("/");
      } else {
        this.$router.push("moviecat");
      }
    },
    //
    policy() {
      console.log(`policy()`);
      if (window.location.hash.includes("privacypolicy")) {
        this.$router.push("/");
      } else {
        this.$router.push("privacypolicy");
      }
    },
    saveLanguData() {
      console.log(`saveLanguData()`);
      localStorage.setItem("i18n_locale", JSON.stringify(this.$i18n.locale));
    },
    saveDarkData() {
      console.log(`saveDarkData()`);
      localStorage.setItem("dark", JSON.stringify(this.$q.dark.isActive));
    },
  },
});
</script>

<template>
  <q-layout view="hHh lpR fFf" class="tw-font-sans">
    <q-header elevated>
      <q-toolbar class="bg-primary text-white rounded-borders tw-bg-blue-600">
        <q-btn flat dense round icon="info" aria-label="Menu" @click="info">
          <q-tooltip class="tw-bg-blue-400">{{ $t("InfoMovieCat") }}</q-tooltip></q-btn
        >
        <q-btn flat dense round icon="policy" aria-label="Policy" @click="policy">
          <q-tooltip class="tw-bg-blue-400">{{ $t("InfoPrivacyPolicy") }}</q-tooltip></q-btn
        >

        <q-icon name="img:icons/icon-256x256.png" size="24px" />
        <q-toolbar-title> MovieCat 0.2.1</q-toolbar-title>
        <q-select
          v-model="$i18n.locale"
          :options="localeOptions"
          dark
          :label="$t('Language')"
          dense
          emit-value
          map-options
          options-dense
          style="min-width: 100px"
          :popup-content-style="{ backgroundColor: '#99ccff' }"
          @update:model-value="saveLanguData"
        />
        <q-btn
          flat
          dense
          round
          :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'"
          aria-label="Dark"
          @click="
            {
              $q.dark.toggle();
              saveDarkData();
            }
          "
        >
          <q-tooltip class="tw-bg-blue-400">{{ $t("InfoToggle") }}</q-tooltip></q-btn
        >
      </q-toolbar>
    </q-header>

    <!-- wird benötigt für layout; nicht löschen!-->
    <q-footer elevated> </q-footer>
    <q-page-container>
      <router-view :langu="$i18n.locale" />
    </q-page-container>
  </q-layout>
</template>

<style></style>
