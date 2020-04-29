<template>
  <div class="ion-page">
    <ion-header>
      <ion-toolbar>
        <ion-title>New URL</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <ion-item>
        <ion-textarea
          :value="URLs"
          @ionInput="URLs = $event.target.value"
          placeholder="Enter URLs"
          rows="20"
        >
        </ion-textarea>
      </ion-item>
      <ion-button @click="onClick()" full>New</ion-button>
    </ion-content>
  </div>
</template>

<script>
import axios from "axios";
const baseUrl = process.env.VUE_APP_API_BASE_URL;
export default {
  name: "home",
  data: function() {
    return {
      URLs: ""
    };
  },
  methods: {
    async onClick() {
      try {
        const response = await axios.post(`${baseUrl}hashs`, {
          values: this.URLs.split('\n').filter(URL => URL.length > 0)
        });
        console.log(response.data);
        this.toCSV(response.data);
        this.URLs = "";
      } catch (e) {
        console.log(e);
      }
    },
    toCSV: function(data) {
        try {
            var csv = "短縮URL,URL\n";
            data.forEach(element => {
                var line;
                line = baseUrl + "url/" + element["Key"];
                line += "," + element["Value"] + "\n";
                csv += line;
            });
            let blob = new Blob([csv], { type: "text/csv" });
            let link = document.createElement("a");
            link.href = window.URL.createObjectURL(blob);
            link.download = "URL.csv";
            link.click();
        } catch (e) {
            console.log(e);
        }
    }
  }
};
</script>
