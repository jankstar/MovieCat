<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "PlayerComp",
  props: ["langu", "modelValue"],
  emits: ["update:modelValue"],

  components: {},

  data: () => {
    return {
      FileData: {
        Name: "test",
        StartTime: undefined,
        EndTime: undefined,
        Size: 0,
        BlobList: [],
        Duration: 0, //Attention only integer in sec !
        Position: 0, //Attention only integer in sec !
      },
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
  },
});
</script>

<template>
  <div></div>
</template>
