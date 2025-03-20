import { createRouter, createWebHistory } from 'vue-router';
import MenuPrincipal from '@/components/MenuPrincipal.vue';

const routes = [
  {
    path: '/',
    name: 'MenuPrincipal',
    component: MenuPrincipal
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;

