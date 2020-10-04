<template>
  <div id="app">
    <h1>SegPic</h1>
    <p>(please click the images)</p>
    <Images v-bind:images="images" />
  </div>
</template>

<script>
import Images from "./components/Images.vue";

export default {
  name: "App",
  components: {
    Images,
  },
  data() {
    return { images: [] };
  },
  created() {
    fetch("http://localhost:8081/api/images")
      .then((res) => {
        return res.json();
      })
      .then((json) => {
        this.images = json;
      })
      .catch((err) => console.error(err));
  },
};
</script>

<style>
h1 {
  font-size: xxx-large;
  justify-content: center;
  text-align: center;
}
p {
  font-size: large;
  justify-content: center;
  text-align: center;
}
</style>
