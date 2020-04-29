<template>
    <div>
        <ion-tabs>
            <ion-tab tab="home">
                <Home />
            </ion-tab>
            <ion-tab tab="hashs">
                <Hashs />
            </ion-tab>
            <ion-tab-bar slot="bottom">
                <ion-tab-button tab="home">
                    <ion-icon name="add-circle-outline"></ion-icon>
                    <ion-label>New URL</ion-label>
                </ion-tab-button>
                <ion-tab-button tab="hashs">
                    <ion-icon name="list"></ion-icon>
                    <ion-label>URL</ion-label>
                </ion-tab-button>
            </ion-tab-bar>
        </ion-tabs>
    </div>
</template>

<script>
const Home = () => import("@/views/Home.vue");
const Hashs = () => import("@/views/Hashs.vue")
export default {
    name: "tab",
    components: { Home, Hashs }
};
</script>