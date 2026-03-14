<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from 'primevue/card'
import Button from 'primevue/button'
import { useAuthStore } from '@/stores/auth'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'
import { useApi } from '@/composables/useApi'

const router = useRouter()
const authStore = useAuthStore()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()
const api = useApi()

const connectionCount = ref<number | null>(null)
const resourceCount = ref<number | null>(null)
const statsLoading = ref(false)

async function fetchProjectStats() {
  const project = projectStore.currentProject
  if (!project) {
    connectionCount.value = null
    resourceCount.value = null
    return
  }
  statsLoading.value = true
  try {
    const [conns, resources] = await Promise.all([
      api.connections.list(project.id),
      api.resources.list(project.id),
    ])
    connectionCount.value = conns.data.length
    resourceCount.value = resources.data.length
  } catch {
    connectionCount.value = null
    resourceCount.value = null
  } finally {
    statsLoading.value = false
  }
}

onMounted(fetchProjectStats)

watch(() => projectStore.currentProject, fetchProjectStats)
</script>

<template>
  <div class="space-y-8">
    <!-- Welcome -->
    <div>
      <h1 class="text-3xl font-bold text-color">
        Welcome{{ authStore.user?.displayName ? `, ${authStore.user.displayName}` : '' }}
      </h1>
      <p class="text-muted-color mt-1">Here's an overview of your workspace.</p>
    </div>

    <!-- Org & Project cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Organization card -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-building text-primary"></i>
            <span>Organization</span>
          </div>
        </template>
        <template #content>
          <div v-if="orgStore.currentOrg" class="space-y-2">
            <p class="text-lg font-semibold text-color">
              {{ orgStore.currentOrg.displayName }}
            </p>
            <p class="text-sm text-muted-color">ID: {{ orgStore.currentOrg.id }}</p>
            <p class="text-sm text-muted-color">
              Total: {{ orgStore.organizations.length }} organization{{
                orgStore.organizations.length !== 1 ? 's' : ''
              }}
            </p>
          </div>
          <div v-else class="space-y-3">
            <p class="text-muted-color">No organization selected.</p>
            <Button
              label="Create Organization"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              @click="router.push({ name: 'organizations' })"
            />
          </div>
        </template>
      </Card>

      <!-- Project card -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-folder text-primary"></i>
            <span>Project</span>
          </div>
        </template>
        <template #content>
          <div v-if="projectStore.currentProject" class="space-y-2">
            <p class="text-lg font-semibold text-color">
              {{ projectStore.currentProject.name }}
            </p>
            <p v-if="projectStore.currentProject.description" class="text-sm text-muted-color">
              {{ projectStore.currentProject.description }}
            </p>
            <p class="text-sm text-muted-color">
              Total: {{ projectStore.projects.length }} project{{
                projectStore.projects.length !== 1 ? 's' : ''
              }}
            </p>
          </div>
          <div v-else class="space-y-3">
            <p class="text-muted-color">
              {{ orgStore.currentOrg ? 'No project selected.' : 'Select an organization first.' }}
            </p>
            <Button
              v-if="orgStore.currentOrg"
              label="Create Project"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              @click="router.push({ name: 'projects' })"
            />
          </div>
        </template>
      </Card>
    </div>

    <!-- Stats cards (visible when project is selected) -->
    <div v-if="projectStore.currentProject" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card>
        <template #content>
          <div class="flex items-center gap-4">
            <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-blue-100 text-blue-600">
              <i class="pi pi-link text-xl"></i>
            </div>
            <div>
              <p class="text-sm text-muted-color">Connections</p>
              <p class="text-2xl font-bold text-color">
                <i v-if="statsLoading" class="pi pi-spinner pi-spin text-lg"></i>
                <span v-else>{{ connectionCount ?? '—' }}</span>
              </p>
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #content>
          <div class="flex items-center gap-4">
            <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-green-100 text-green-600">
              <i class="pi pi-box text-xl"></i>
            </div>
            <div>
              <p class="text-sm text-muted-color">Resources</p>
              <p class="text-2xl font-bold text-color">
                <i v-if="statsLoading" class="pi pi-spinner pi-spin text-lg"></i>
                <span v-else>{{ resourceCount ?? '—' }}</span>
              </p>
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #content>
          <div class="flex items-center gap-4">
            <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-purple-100 text-purple-600">
              <i class="pi pi-folder text-xl"></i>
            </div>
            <div>
              <p class="text-sm text-muted-color">Projects</p>
              <p class="text-2xl font-bold text-color">{{ projectStore.projects.length }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card>
        <template #content>
          <div class="flex items-center gap-4">
            <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-amber-100 text-amber-600">
              <i class="pi pi-building text-xl"></i>
            </div>
            <div>
              <p class="text-sm text-muted-color">Organizations</p>
              <p class="text-2xl font-bold text-color">{{ orgStore.organizations.length }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Quick actions -->
    <Card>
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-bolt text-primary"></i>
          <span>Quick Actions</span>
        </div>
      </template>
      <template #content>
        <div class="flex flex-wrap gap-3">
          <Button
            label="Organizations"
            icon="pi pi-building"
            severity="secondary"
            outlined
            @click="router.push({ name: 'organizations' })"
          />
          <Button
            label="Projects"
            icon="pi pi-folder"
            severity="secondary"
            outlined
            @click="router.push({ name: 'projects' })"
          />
          <Button
            v-if="projectStore.currentProject"
            label="Connections"
            icon="pi pi-link"
            severity="secondary"
            outlined
            @click="
              router.push({
                name: 'connections',
                params: { projectId: projectStore.currentProject!.id },
              })
            "
          />
          <Button
            v-if="projectStore.currentProject"
            label="Resources"
            icon="pi pi-box"
            severity="secondary"
            outlined
            @click="
              router.push({
                name: 'resources',
                params: { projectId: projectStore.currentProject!.id },
              })
            "
          />
          <Button
            label="Providers"
            icon="pi pi-server"
            severity="secondary"
            outlined
            @click="router.push({ name: 'providers' })"
          />
        </div>
      </template>
    </Card>
  </div>
</template>
