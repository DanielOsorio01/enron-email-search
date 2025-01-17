import { createApp } from 'vue';
import App from './App.vue';
import './styles.css'; // Import Tailwind CSS

document.addEventListener("DOMContentLoaded", () => {
  //console.log("Config Loaded:", window.APP_CONFIG); // Debugging

  //const backend_URL = window.APP_CONFIG?.backend_URL || 'fallback_backend_ip';
  //console.log("Backend URL:", backend_URL); // Verify if it's set

  // Initialize Vue App after config is available
  createApp(App).mount("#app");
});
