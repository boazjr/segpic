<template>
  <div class="images">
    <div v-bind:key="image.id" v-for="image in images">
      <ImageItem v-bind:image="image" />
    </div>
  </div>
</template>

<script>
// import ImageItem from "./ImageItem.vue";
import ImageItem from "./ImageItem.vue";

export default {
  name: "Images",
  components: {
    ImageItem,
  },
  props: ["images"],
};
</script>

<style scoped>
ImageItem {
  margin-left: auto;
  margin-right: auto;
}
</style>