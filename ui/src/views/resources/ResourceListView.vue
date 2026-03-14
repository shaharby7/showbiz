<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Resource, Connection } from '@showbiz/sdk'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'
import { useProjectStore } from '@/stores/project'

const route = useRoute()
const router = useRouter()
const api = useApi()
const projectStore = useProjectStore()

const projectId = computed(() => (route.params.projectId as string) || projectStore.currentProject?.id)

const resources = ref<Resource[]>([])
const connections = ref<Connection[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showDeleteDialog = ref(false)
const resourceToDelete = ref<Resource | null>(null)
const deleting = ref(false)
const deleteError = ref<string | null>(null)

function connectionName(connectionId: string): string {
  const conn = connections.value.find((c) => c.id === connectionId)
  return conn?.name || connectionId
}

function statusSeverity(status: string): 'success' | 'warn' | 'danger' | 'info' {
  switch (status) {
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
      return 'warn'
    case 'error':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

async function fetchData() {
  if (!projectId.value) return
  loading.value = true
  error.value = null
  try {
    const [resourceResult, connectionResult] = await Promise.all([
      api.resources.list(projectId.value),
      api.connections.list(projectId.value),
    ])
    resources.value = resourceResult.data
    connections.value = connectionResult.data
  } catch (e: any) {
    error.value = e.message || 'Failed to load resources'
  } finally {
    loading.value = false
  }
}

function confirmDelete(resource: Resource) {
  resourceToDelete.value = resource
  deleteError.value = null
  showDeleteDialog.value = true
}

async function deleteResource() {
  if (!projectId.value || !resourceToDelete.value) return
  deleting.value = true
  deleteError.value = null
  try {
    await api.resources.delete(projectId.value, resourceToDelete.value.id)
    showDeleteDialog.value = false
    resourceToDelete.value = null
    await fetchData()
  } catch (e: any) {
    deleteError.value = e.message || 'Failed to delete resource'
  } finally {
    deleting.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-color">Resources</h1>
        <p v-if="projectStore.currentProject" class="text-muted-color mt-1">
          Resources in {{ projectStore.currentProject.name }}
        </p>
      </div>
      <Button
        v-if="projectId"
        label="Create Resource"
        icon="pi pi-plus"
        @click="router.push({ name: 'resource-create', params: { projectId: projectId } })"
      />
    </div>

    <!-- No project selected -->
    <Message v-if="!projectId" severity="warn" :closable="false">
      Please select a project to view resources.
    </Message>

    <!-- Error -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- DataTable -->
    <DataTable
      v-if="projectId && !loading"
      :value="resources"
      :rowHover="true"
      stripedRows
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-box text-4xl mb-3 block"></i>
          <p>No resources yet. Create your first resource to get started.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="resourceType" header="Type" sortable>
        <template #body="{ data }">
          <span class="text-color">{{ data.resourceType }}</span>
        </template>
      </Column>
      <Column field="connectionId" header="Connection" sortable>
        <template #body="{ data }">
          <span class="text-muted-color">{{ connectionName(data.connectionId) }}</span>
        </template>
      </Column>
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :value="data.status" :severity="statusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="createdAt" header="Created" sortable>
        <template #body="{ data }">
          <span class="text-muted-color">{{ formatDate(data.createdAt) }}</span>
        </template>
      </Column>
      <Column header="Actions" :style="{ width: '6rem' }">
        <template #body="{ data }">
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            size="small"
            @click="confirmDelete(data)"
          />
        </template>
      </Column>
    </DataTable>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Delete Resource"
      :modal="true"
      :style="{ width: '26rem' }"
    >
      <div class="space-y-3">
        <Message v-if="deleteError" severity="error" :closable="false">{{ deleteError }}</Message>
        <p class="text-color">
          Are you sure you want to delete <strong>{{ resourceToDelete?.name }}</strong>?
        </p>
        <p class="text-sm text-muted-color">This action cannot be undone.</p>
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
          @click="deleteResource"
          :loading="deleting"
        />
      </template>
    </Dialog>
  </div>
</template>
