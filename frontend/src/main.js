import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Asegúrate de importar el enrutador

const app = createApp(App);
app.use(router); // Usar el enrutador
app.mount('#app');

