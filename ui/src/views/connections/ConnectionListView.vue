<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Connection } from '@showbiz/sdk'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'
import { useProjectStore } from '@/stores/project'

const route = useRoute()
const router = useRouter()
const api = useApi()
const projectStore = useProjectStore()

const connections = ref<Connection[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showDeleteDialog = ref(false)
const connectionToDelete = ref<Connection | null>(null)
const deleting = ref(false)
const deleteError = ref<string | null>(null)

const projectId = computed(() => (route.params.projectId as string) || projectStore.currentProject?.id)

async function fetchConnections() {
  if (!projectId.value) return
  loading.value = true
  error.value = null
  try {
    const result = await api.connections.list(projectId.value)
    connections.value = result.data
  } catch (e: any) {
    error.value = e.message || 'Failed to load connections'
  } finally {
    loading.value = false
  }
}

function navigateToCreate() {
  router.push({ name: 'connection-create', params: { projectId: projectId.value } })
}

function confirmDelete(connection: Connection) {
  connectionToDelete.value = connection
  deleteError.value = null
  showDeleteDialog.value = true
}

async function deleteConnection() {
  if (!projectId.value || !connectionToDelete.value) return
  deleting.value = true
  deleteError.value = null
  try {
    await api.connections.delete(projectId.value, connectionToDelete.value.id)
    showDeleteDialog.value = false
    connectionToDelete.value = null
    await fetchConnections()
  } catch (e: any) {
    deleteError.value = e.message || 'Failed to delete connection'
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

onMounted(fetchConnections)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Button
          icon="pi pi-arrow-left"
          severity="secondary"
          text
          rounded
          @click="router.push({ name: 'project-detail', params: { projectId: projectId } })"
        />
        <div>
          <h1 class="text-3xl font-bold text-color">Connections</h1>
          <p v-if="projectStore.currentProject" class="text-muted-color mt-1">
            Connections in {{ projectStore.currentProject.name }}
          </p>
        </div>
      </div>
      <Button
        v-if="projectId"
        label="Create Connection"
        icon="pi pi-plus"
        @click="navigateToCreate"
      />
    </div>

    <!-- No project selected -->
    <Message v-if="!projectId" severity="warn" :closable="false">
      Please select a project to view connections.
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
      :value="connections"
      :rowHover="true"
      stripedRows
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-link text-4xl mb-3 block"></i>
          <p>No connections yet. Create your first connection to get started.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="provider" header="Provider" sortable>
        <template #body="{ data }">
          <span class="text-color">{{ data.provider }}</span>
        </template>
      </Column>
      <Column field="createdAt" header="Created" sortable>
        <template #body="{ data }">
          <span class="text-muted-color">{{ formatDate(data.createdAt) }}</span>
        </template>
      </Column>
      <Column header="" style="width: 5rem">
        <template #body="{ data }">
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            size="small"
            @click.stop="confirmDelete(data)"
          />
        </template>
      </Column>
    </DataTable>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="showDeleteDialog"
      header="Delete Connection"
      :modal="true"
      :style="{ width: '26rem' }"
    >
      <div class="space-y-3">
        <Message v-if="deleteError" severity="error" :closable="false">{{ deleteError }}</Message>
        <p class="text-color">
          Are you sure you want to delete <strong>{{ connectionToDelete?.name }}</strong>?
        </p>
        <p class="text-sm text-muted-color">
          This will fail if the connection has associated resources.
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
          @click="deleteConnection"
          :loading="deleting"
        />
      </template>
    </Dialog>
  </div>
</template>
