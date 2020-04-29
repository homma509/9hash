import Vue from "vue";
import { IonicVueRouter } from "@ionic/vue";

Vue.use(IonicVueRouter);

const routes = [
  {
    path: "/",
    redirect: "/home",
    component: () => import("@/views/Tab.vue"),
    children: [
      {
        path: "/home",
        name: "home",
        component: () => import("@/views/Home.vue")
      },
      {
        path: "/hashs",
        name: "hashs",
        component: () => import("@/views/Hashs.vue")
      }
    ]
  }
];

const router = new IonicVueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
