<template>
  <div class="ion-page">
    <ion-header>
      <ion-toolbar>
        <ion-title>URL</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <ion-list v-bind:key="hash.ID" v-for="hash in hashs">
        <ion-item>
          <ion-label>
            <ion-router-link
              target="_blank"
              :href="`${baseUrl}/url/${hash.Key}`"
              >{{ hash.Key }}</ion-router-link
            >
          </ion-label>
          <ion-label>{{ hash.Value }}</ion-label>
          <ion-button @click="deleteHash(hash.ID)" full>delete</ion-button>
        </ion-item>
      </ion-list>
      <ion-button @click="onClick()" full>CSV</ion-button>
    </ion-content>
  </div>
</template>

<script>
import axios from "axios";
const baseUrl = process.env.VUE_APP_API_BASE_URL;
export default {
  name: "hashs",
  data() {
    return {
      hashs: null,
      baseUrl: baseUrl
    };
  },
  watch: {
    $route: "reload"
  },
  async created() {
    await this.reload();
  },
  methods: {
    async reload() {
      console.log(baseUrl);
      if (this.$route.fullPath == "/hashs") {
        try {
          const response = await axios.get(`${baseUrl}/hashs`);
          this.hashs = response.data;
        } catch (e) {
          console.log(e);
        }
      }
    },
    async deleteHash(id) {
      try {
        await axios.delete(`${baseUrl}/hashs/${id}`);
      } catch (e) {
        console.log(e);
      }
      await this.reload();
    },
    async onClick() {
      try {
        let csv = "短縮URL,URL\n";
        this.hashs.forEach(element => {
          let line;
          line = baseUrl + "url/" + element["Key"];
          line += "," + element["Value"] + "\n";
          csv += line;
        });
        let blob = new Blob([csv], { type: "text/csv" });
        let link = document.createElement("a");
        let d = new Date();
        link.href = window.URL.createObjectURL(blob);
        link.download = "URL_" + d.toJSON() + ".csv";
        link.click();
      } catch (e) {
        console.log(e);
      }
    }
  }
};
</script>
