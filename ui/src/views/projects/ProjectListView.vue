<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { Project, CreateProjectInput } from '@showbiz/sdk'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'

const router = useRouter()
const api = useApi()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()

const projects = ref<Project[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showCreateDialog = ref(false)
const createForm = ref<CreateProjectInput>({ name: '', description: '' })
const creating = ref(false)
const createError = ref<string | null>(null)

async function fetchProjects() {
  const orgId = orgStore.currentOrg?.id
  if (!orgId) return
  loading.value = true
  error.value = null
  try {
    const result = await api.projects.list(orgId)
    projects.value = result.data
  } catch (e: any) {
    error.value = e.message || 'Failed to load projects'
  } finally {
    loading.value = false
  }
}

async function createProject() {
  const orgId = orgStore.currentOrg?.id
  if (!orgId) return
  creating.value = true
  createError.value = null
  try {
    const project = await api.projects.create(orgId, createForm.value)
    showCreateDialog.value = false
    createForm.value = { name: '', description: '' }
    await projectStore.fetchProjects(orgId)
    projectStore.selectProject(project)
    await fetchProjects()
  } catch (e: any) {
    createError.value = e.message || 'Failed to create project'
  } finally {
    creating.value = false
  }
}

function openCreateDialog() {
  createForm.value = { name: '', description: '' }
  createError.value = null
  showCreateDialog.value = true
}

function onRowClick(event: any) {
  const project = event.data as Project
  projectStore.selectProject(project)
  router.push({ name: 'project-detail', params: { projectId: project.id } })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

onMounted(fetchProjects)
watch(() => orgStore.currentOrg, fetchProjects)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-color">Projects</h1>
        <p v-if="orgStore.currentOrg" class="text-muted-color mt-1">
          Projects in {{ orgStore.currentOrg.displayName }}
        </p>
      </div>
      <Button
        v-if="orgStore.currentOrg"
        label="Create Project"
        icon="pi pi-plus"
        @click="openCreateDialog"
      />
    </div>

    <!-- No org selected -->
    <Message v-if="!orgStore.currentOrg" severity="warn" :closable="false">
      Please select an organization to view projects.
    </Message>

    <!-- Error -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- DataTable -->
    <DataTable
      v-if="orgStore.currentOrg && !loading"
      :value="projects"
      :rowHover="true"
      stripedRows
      @row-click="onRowClick"
      class="cursor-pointer"
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-folder text-4xl mb-3 block"></i>
          <p>No projects yet. Create your first project to get started.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="description" header="Description">
        <template #body="{ data }">
          <span class="text-muted-color">{{ data.description || '—' }}</span>
        </template>
      </Column>
      <Column field="createdAt" header="Created" sortable>
        <template #body="{ data }">
          <span class="text-muted-color">{{ formatDate(data.createdAt) }}</span>
        </template>
      </Column>
    </DataTable>

    <!-- Create Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      header="Create Project"
      :modal="true"
      :style="{ width: '28rem' }"
    >
      <div class="space-y-4">
        <Message v-if="createError" severity="error" :closable="false">{{ createError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="project-name" class="font-semibold text-color">Name</label>
          <InputText
            id="project-name"
            v-model="createForm.name"
            placeholder="my-project"
            :disabled="creating"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="project-desc" class="font-semibold text-color">Description</label>
          <Textarea
            id="project-desc"
            v-model="createForm.description"
            rows="3"
            placeholder="Optional description"
            :disabled="creating"
          />
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showCreateDialog = false"
          :disabled="creating"
        />
        <Button
          label="Create"
          icon="pi pi-check"
          @click="createProject"
          :loading="creating"
          :disabled="!createForm.name.trim()"
        />
      </template>
    </Dialog>
  </div>
</template>
