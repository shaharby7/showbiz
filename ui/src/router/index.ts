import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/views/DashboardView.vue'),
      },
      {
        path: 'organizations',
        name: 'organizations',
        component: () => import('@/views/organizations/OrganizationListView.vue'),
      },
      {
        path: 'organizations/:orgId',
        name: 'organization-detail',
        component: () => import('@/views/organizations/OrganizationDetailView.vue'),
      },
      {
        path: 'organizations/:orgId/members',
        name: 'organization-members',
        component: () => import('@/views/organizations/OrganizationMembersView.vue'),
      },
      {
        path: 'projects',
        name: 'projects',
        component: () => import('@/views/projects/ProjectListView.vue'),
      },
      {
        path: 'projects/:projectId',
        name: 'project-detail',
        component: () => import('@/views/projects/ProjectDetailView.vue'),
      },
      {
        path: 'projects/:projectId/connections',
        name: 'connections',
        component: () => import('@/views/connections/ConnectionListView.vue'),
      },
      {
        path: 'projects/:projectId/connections/new',
        name: 'connection-create',
        component: () => import('@/views/connections/ConnectionCreateView.vue'),
      },
      {
        path: 'projects/:projectId/resources',
        name: 'resources',
        component: () => import('@/views/resources/ResourceListView.vue'),
      },
      {
        path: 'projects/:projectId/resources/new',
        name: 'resource-create',
        component: () => import('@/views/resources/ResourceCreateView.vue'),
      },
      {
        path: 'projects/:projectId/iam',
        name: 'iam',
        component: () => import('@/views/iam/IAMView.vue'),
      },
      {
        path: 'providers',
        name: 'providers',
        component: () => import('@/views/providers/ProviderListView.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const authStore = useAuthStore()
  if (!to.meta.public && !authStore.isAuthenticated) {
    return { name: 'login' }
  }
})

export default router
