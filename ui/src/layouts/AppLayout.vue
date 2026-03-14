<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Select from 'primevue/select'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import { useAuthStore } from '@/stores/auth'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'
import type { Organization, Project } from '@showbiz/sdk'

const router = useRouter()
const authStore = useAuthStore()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()

const userMenu = ref<InstanceType<typeof Menu> | null>(null)

const selectedOrg = computed({
  get: () => orgStore.currentOrg,
  set: (val: Organization | null) => {
    orgStore.selectOrganization(val)
    if (val) {
      projectStore.fetchProjects(val.id)
      projectStore.selectProject(null)
    }
  },
})

const selectedProject = computed({
  get: () => projectStore.currentProject,
  set: (val: Project | null) => projectStore.selectProject(val),
})

const hasProject = computed(() => !!projectStore.currentProject)

const sidebarItems = computed(() => {
  const items = [
    { label: 'Dashboard', icon: 'pi pi-home', to: '/' },
    { label: 'Organizations', icon: 'pi pi-building', to: '/organizations' },
    { label: 'Projects', icon: 'pi pi-folder', to: '/projects' },
  ]
  if (hasProject.value) {
    const pid = projectStore.currentProject!.id
    items.push(
      { label: 'Connections', icon: 'pi pi-link', to: `/projects/${pid}/connections` },
      { label: 'Resources', icon: 'pi pi-box', to: `/projects/${pid}/resources` },
      { label: 'IAM', icon: 'pi pi-shield', to: `/projects/${pid}/iam` },
    )
  }
  items.push({ label: 'Providers', icon: 'pi pi-server', to: '/providers' })
  return items
})

const userMenuItems = ref([
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: () => {
      authStore.logout()
      router.push({ name: 'login' })
    },
  },
])

function toggleUserMenu(event: Event) {
  userMenu.value?.toggle(event)
}

onMounted(async () => {
  try {
    await orgStore.initialize()
    if (orgStore.currentOrg) {
      await projectStore.fetchProjects(orgStore.currentOrg.id)
      await projectStore.initialize()
    }
  } catch {
    // org/project load may fail if not set up yet
  }
})

watch(() => orgStore.currentOrg, async (newOrg) => {
  if (newOrg) {
    await projectStore.fetchProjects(newOrg.id)
  }
})
</script>

<template>
  <div class="flex flex-col h-screen">
    <!-- Top navbar -->
    <header class="flex items-center justify-between px-4 py-2 border-b border-gray-200 bg-white">
      <div class="flex items-center gap-4">
        <span class="text-xl font-bold text-primary">Showbiz</span>
        <Select
          v-model="selectedOrg"
          :options="orgStore.organizations"
          optionLabel="displayName"
          placeholder="Select Organization"
          class="w-56"
        />
        <Select
          v-model="selectedProject"
          :options="projectStore.projects"
          optionLabel="name"
          placeholder="Select Project"
          class="w-56"
        />
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm text-gray-600">{{ authStore.user?.email }}</span>
        <Button
          icon="pi pi-user"
          severity="secondary"
          text
          rounded
          @click="toggleUserMenu"
        />
        <Menu ref="userMenu" :model="userMenuItems" :popup="true" />
      </div>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar -->
      <nav class="w-56 border-r border-gray-200 bg-gray-50 overflow-y-auto">
        <ul class="py-2">
          <li v-for="item in sidebarItems" :key="item.to">
            <router-link
              :to="item.to"
              class="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
              active-class="bg-gray-200 font-semibold"
            >
              <i :class="item.icon"></i>
              {{ item.label }}
            </router-link>
          </li>
        </ul>
      </nav>

      <!-- Main content -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>
