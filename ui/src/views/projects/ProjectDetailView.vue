<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Project } from '@showbiz/sdk'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Textarea from 'primevue/textarea'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'

const route = useRoute()
const router = useRouter()
const api = useApi()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()

const project = ref<Project | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const editing = ref(false)
const editDescription = ref('')
const saving = ref(false)
const saveError = ref<string | null>(null)

const showDeleteDialog = ref(false)
const deleting = ref(false)
const deleteError = ref<string | null>(null)

const projectId = route.params.projectId as string

async function fetchProject() {
  const orgId = orgStore.currentOrg?.id
  if (!orgId) {
    error.value = 'No organization selected'
    return
  }
  loading.value = true
  error.value = null
  try {
    project.value = await api.projects.get(orgId, projectId)
  } catch (e: any) {
    error.value = e.message || 'Failed to load project'
  } finally {
    loading.value = false
  }
}

function startEdit() {
  editDescription.value = project.value?.description || ''
  saveError.value = null
  editing.value = true
}

function cancelEdit() {
  editing.value = false
  saveError.value = null
}

async function saveDescription() {
  const orgId = orgStore.currentOrg?.id
  if (!orgId || !project.value) return
  saving.value = true
  saveError.value = null
  try {
    project.value = await api.projects.update(orgId, projectId, {
      description: editDescription.value,
    })
    editing.value = false
  } catch (e: any) {
    saveError.value = e.message || 'Failed to update project'
  } finally {
    saving.value = false
  }
}

async function deleteProject() {
  const orgId = orgStore.currentOrg?.id
  if (!orgId) return
  deleting.value = true
  deleteError.value = null
  try {
    await api.projects.delete(orgId, projectId)
    if (projectStore.currentProject?.id === projectId) {
      projectStore.selectProject(null)
    }
    router.push({ name: 'projects' })
  } catch (e: any) {
    deleteError.value = e.message || 'Failed to delete project'
  } finally {
    deleting.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(fetchProject)
</script>

<template>
  <div class="space-y-6">
    <!-- Back + Header -->
    <div class="flex items-center gap-3">
      <Button
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        rounded
        @click="router.push({ name: 'projects' })"
      />
      <div class="flex-1">
        <h1 class="text-3xl font-bold text-color">
          {{ project?.name || 'Project Detail' }}
        </h1>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- Error -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <template v-if="project && !loading">
      <!-- Details Card -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-folder text-primary"></i>
            <span>Project Details</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p class="text-sm font-semibold text-muted-color">Name</p>
                <p class="text-color text-lg">{{ project.name }}</p>
              </div>
              <div>
                <p class="text-sm font-semibold text-muted-color">Organization ID</p>
                <p class="text-color">{{ project.organizationId }}</p>
              </div>
              <div>
                <p class="text-sm font-semibold text-muted-color">Created</p>
                <p class="text-color">{{ formatDate(project.createdAt) }}</p>
              </div>
              <div>
                <p class="text-sm font-semibold text-muted-color">Updated</p>
                <p class="text-color">{{ formatDate(project.updatedAt) }}</p>
              </div>
            </div>

            <!-- Description (inline edit) -->
            <div>
              <div class="flex items-center gap-2 mb-1">
                <p class="text-sm font-semibold text-muted-color">Description</p>
                <Button
                  v-if="!editing"
                  icon="pi pi-pencil"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  @click="startEdit"
                />
              </div>
              <template v-if="!editing">
                <p class="text-color">{{ project.description || 'No description' }}</p>
              </template>
              <template v-else>
                <Message v-if="saveError" severity="error" :closable="false" class="mb-2">
                  {{ saveError }}
                </Message>
                <Textarea
                  v-model="editDescription"
                  rows="3"
                  class="w-full"
                  :disabled="saving"
                />
                <div class="flex gap-2 mt-2">
                  <Button
                    label="Save"
                    icon="pi pi-check"
                    size="small"
                    @click="saveDescription"
                    :loading="saving"
                  />
                  <Button
                    label="Cancel"
                    severity="secondary"
                    text
                    size="small"
                    @click="cancelEdit"
                    :disabled="saving"
                  />
                </div>
              </template>
            </div>
          </div>
        </template>
      </Card>

      <!-- Navigation Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card
          class="cursor-pointer hover:shadow-lg transition-shadow"
          @click="router.push({ name: 'connections', params: { projectId: project.id } })"
        >
          <template #content>
            <div class="flex items-center gap-4">
              <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-blue-100 text-blue-600">
                <i class="pi pi-link text-xl"></i>
              </div>
              <div>
                <p class="font-semibold text-color">Connections</p>
                <p class="text-sm text-muted-color">Manage provider connections</p>
              </div>
              <i class="pi pi-chevron-right ml-auto text-muted-color"></i>
            </div>
          </template>
        </Card>

        <Card
          class="cursor-pointer hover:shadow-lg transition-shadow"
          @click="router.push({ name: 'resources', params: { projectId: project.id } })"
        >
          <template #content>
            <div class="flex items-center gap-4">
              <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-green-100 text-green-600">
                <i class="pi pi-box text-xl"></i>
              </div>
              <div>
                <p class="font-semibold text-color">Resources</p>
                <p class="text-sm text-muted-color">Manage project resources</p>
              </div>
              <i class="pi pi-chevron-right ml-auto text-muted-color"></i>
            </div>
          </template>
        </Card>

        <Card
          class="cursor-pointer hover:shadow-lg transition-shadow"
          @click="router.push({ name: 'iam', params: { projectId: project.id } })"
        >
          <template #content>
            <div class="flex items-center gap-4">
              <div class="flex items-center justify-center w-12 h-12 rounded-lg bg-purple-100 text-purple-600">
                <i class="pi pi-shield text-xl"></i>
              </div>
              <div>
                <p class="font-semibold text-color">IAM</p>
                <p class="text-sm text-muted-color">Access control & permissions</p>
              </div>
              <i class="pi pi-chevron-right ml-auto text-muted-color"></i>
            </div>
          </template>
        </Card>
      </div>

      <!-- Danger Zone -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2 text-red-500">
            <i class="pi pi-exclamation-triangle"></i>
            <span>Danger Zone</span>
          </div>
        </template>
        <template #content>
          <div class="flex items-center justify-between">
            <div>
              <p class="font-semibold text-color">Delete this project</p>
              <p class="text-sm text-muted-color">
                This action cannot be undone. All connections and resources will be deleted.
              </p>
            </div>
            <Button
              label="Delete Project"
              icon="pi pi-trash"
              severity="danger"
              outlined
              @click="showDeleteDialog = true"
            />
          </div>
        </template>
      </Card>
    </template>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Delete Project"
      :modal="true"
      :style="{ width: '26rem' }"
    >
      <div class="space-y-3">
        <Message v-if="deleteError" severity="error" :closable="false">{{ deleteError }}</Message>
        <p class="text-color">
          Are you sure you want to delete <strong>{{ project?.name }}</strong>?
        </p>
        <p class="text-sm text-muted-color">
          This will permanently delete the project and all associated data.
        </p>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showDeleteDialog = false"
          :disabled="deleting"
        />
        <Button
          label="Delete"
          icon="pi pi-trash"
          severity="danger"
          @click="deleteProject"
          :loading="deleting"
        />
      </template>
    </Dialog>
  </div>
</template>
