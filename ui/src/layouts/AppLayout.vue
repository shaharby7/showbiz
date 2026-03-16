<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Select from 'primevue/select'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import { useAuthStore } from '@/stores/auth'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'
import { useDarkMode } from '@/composables/useDarkMode'
import { useApi } from '@/composables/useApi'
import type { Organization, Project, ResourceTypeInfo } from '@showbiz/sdk'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()
const api = useApi()
const { isDark, toggle: toggleDark } = useDarkMode()

const userMenu = ref<InstanceType<typeof Menu> | null>(null)
const resourceTypes = ref<ResourceTypeInfo[]>([])

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

function typeIcon(name: string): string {
  switch (name) {
    case 'machine': return 'pi pi-desktop'
    case 'network': return 'pi pi-sitemap'
    default: return 'pi pi-box'
  }
}

function typeLabel(name: string): string {
  return name.charAt(0).toUpperCase() + name.slice(1) + 's'
}

function isActiveRoute(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

const userMenuItems = computed(() => [
  {
    label: 'Organizations',
    icon: 'pi pi-building',
    command: () => router.push({ name: 'organizations' }),
  },
  {
    label: 'Projects',
    icon: 'pi pi-folder',
    command: () => router.push({ name: 'projects' }),
  },
  { separator: true },
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

async function fetchResourceTypes() {
  try {
    resourceTypes.value = await api.resourceTypes.list()
  } catch {
    // resource types may not be available yet
  }
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
  await fetchResourceTypes()
})

watch(() => orgStore.currentOrg, async (newOrg) => {
  if (newOrg) {
    await projectStore.fetchProjects(newOrg.id)
  }
})
</script>

<template>
  <div class="flex flex-col h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <!-- Top navbar -->
    <header class="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="flex items-center gap-4">
        <router-link to="/" class="flex items-center gap-2 hover:opacity-80 transition-opacity">
          <img src="/showbiz.svg" alt="Showbiz" class="h-7 w-7" />
          <span class="text-xl font-bold text-primary">Showbiz</span>
        </router-link>
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
        <Button
          :icon="isDark ? 'pi pi-sun' : 'pi pi-moon'"
          severity="secondary"
          text
          rounded
          @click="toggleDark"
          v-tooltip.bottom="isDark ? 'Light mode' : 'Dark mode'"
        />
        <span class="text-sm text-gray-600 dark:text-gray-300">{{ authStore.user?.email }}</span>
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
      <nav class="w-56 border-r border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 overflow-y-auto flex flex-col">
        <div v-if="hasProject" class="flex-1 py-3">
          <!-- Resources section -->
          <div class="px-4 py-1 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
            Resources
          </div>
          <ul>
            <li v-for="rt in resourceTypes" :key="rt.name">
              <router-link
                :to="`/projects/${projectStore.currentProject!.id}/resources/${rt.name}`"
                class="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                :class="{ 'bg-gray-200 dark:bg-gray-600 font-semibold': isActiveRoute(`/projects/${projectStore.currentProject!.id}/resources/${rt.name}`) }"
              >
                <i :class="typeIcon(rt.name)"></i>
                {{ typeLabel(rt.name) }}
              </router-link>
            </li>
          </ul>

          <!-- Configuration section -->
          <div class="px-4 py-1 mt-4 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
            Configuration
          </div>
          <ul>
            <li>
              <router-link
                :to="`/projects/${projectStore.currentProject!.id}/connections`"
                class="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                active-class="bg-gray-200 dark:bg-gray-600 font-semibold"
              >
                <i class="pi pi-link"></i>
                Connections
              </router-link>
            </li>
            <li>
              <router-link
                :to="`/projects/${projectStore.currentProject!.id}/iam`"
                class="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                active-class="bg-gray-200 dark:bg-gray-600 font-semibold"
              >
                <i class="pi pi-shield"></i>
                IAM
              </router-link>
            </li>
          </ul>
        </div>

        <div v-else class="flex-1 flex items-center justify-center p-4">
          <p class="text-sm text-center text-muted-color">Select a project to view resources</p>
        </div>
      </nav>

      <!-- Main content -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>
