import AppLayout from "@/layout/AppLayout.vue";
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: AppLayout,
      children: [
        {
          path: "/pages/home",
          name: "home",
          component: () => import("@/views/pages/Home.vue"),
        },
        {
          path: "/pages/playlist/:playlist_id/:playlist_name?",
          name: "playlist",
          component: () => import("@/views/pages/Playlist.vue"),
          props: true,
        },
        {
          path: "/profile",
          name: "profile",
          component: () => import("@/views/pages/Profile.vue"),
        },
        {
          path: "/pages/search",
          name: "search",
          component: () => import("@/views/pages/Search.vue"),
        },
        {
          path: '/playlist',
          name: 'Playlist',
          component: () => import('@/views/pages/Playlist.vue'), // Adjust path as needed
        },
        {
          path: "/pages/music",
          name: "music",
          component: () => import("@/views/pages/Playlist.vue"),
        },
        {
          path: "/pages/map",
          name: "map",
          component: () => import("@/views/pages/Map.vue"),
        },
      ],
    },
    {
      path: "/pages/notfound",
      name: "notfound",
      component: () => import("@/views/pages/NotFound.vue"),
    },
    {
      path: "/auth/login",
      name: "login",
      component: () => import("@/views/pages/auth/Login.vue"),
    },
    {
      path: "/auth/access",
      name: "accessDenied",
      component: () => import("@/views/pages/auth/Access.vue"),
    },
    {
      path: "/auth/error",
      name: "error",
      component: () => import("@/views/pages/auth/Error.vue"),
    },
    {
      path: "/createaccount",
      name: "createaccount",
      component: () => import("@/views/pages/CreateAccount.vue"),
    },
    {
      path: "/pages/chat/:user_id?",
      name: "chat",
      component: () => import("@/views/pages/Chat.vue"),
    },
  ],
});

export default router;
