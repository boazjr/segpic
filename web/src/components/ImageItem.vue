<template>
  <div>
    <div>
      <img
        :src="image.download_url"
        v-bind:class="{ flag: image.flag }"
        v-on:click="flagImage"
      />
      <p>(By:{{ image.author }})</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "ImageItem",
  props: ["image"],
  methods: {
    flagImage() {
      fetch(`http://localhost:8081/api/images/${this.image.id}`, {
        method: "PATCH",
        body: JSON.stringify({
          flag: !this.image.flag,
        }),
        headers: {
          "Content-type": "application/json; charset=UTF-8",
        },
      })
        .then((res) => {
          return res.json();
        })
        .then((img) => {
          this.image.flag = img.flag;
        })
        .catch((err) => console.error(err));
    },
  },
};
</script>

<style scoped>
.flag {
  border: 5px solid rgb(0, 233, 70);
}
div {
  padding-left: 5px;
  padding-right: 5px;
  width: 500px;
  float: left;
}
p {
  margin: 0px;
  text-align: right;
}
</style>