import { fileURLToPath, URL } from "node:url";
import { PrimeVueResolver } from "@primevue/auto-import-resolver";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/vite";
import { defineConfig, loadEnv } from "vite";
import path from "path";
import dotenv from "dotenv";

dotenv.config({ path: path.resolve(__dirname, ".env") });

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, ".."), "");

  return {
    optimizeDeps: {
      noDiscovery: true,
      include: [
        '@supabase/supabase-js', // Existing optimization
        'node:fs', // Example Node.js package for optimization
        'node:path', // Example Node.js package
        // Add other Node.js dependencies if needed
      ],

    },
    plugins: [
      vue(),
      Components({
        resolvers: [PrimeVueResolver()],
      }),
    ],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    define: {
      "process.env": env, // Ensure variables are accessible
    },
    build: {
      target: 'esnext', // Ensure compatibility with top-level await and modern JS features
    },
  };
});
