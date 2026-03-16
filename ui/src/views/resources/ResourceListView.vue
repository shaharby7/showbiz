<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Resource, Connection, ResourceTypeInfo } from '@showbiz/sdk'
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
const filterType = computed(() => (route.params.resourceType as string) || null)

const resources = ref<Resource[]>([])
const connections = ref<Connection[]>([])
const resourceTypes = ref<ResourceTypeInfo[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showDeleteDialog = ref(false)
const resourceToDelete = ref<Resource | null>(null)
const deleting = ref(false)
const deleteError = ref<string | null>(null)

const currentTypeInfo = computed(() =>
  filterType.value ? resourceTypes.value.find((rt) => rt.name === filterType.value) || null : null,
)

const filteredResources = computed(() => {
  if (!filterType.value) return resources.value
  return resources.value.filter((r) => r.resourceType === filterType.value)
})

function connectionName(connectionId: string | null): string {
  if (!connectionId) return '—'
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

function fieldValue(resource: Resource, fieldName: string): string {
  const val = resource.values?.[fieldName]
  return val !== undefined && val !== null ? String(val) : '—'
}

async function fetchData() {
  if (!projectId.value) return
  loading.value = true
  error.value = null
  try {
    const [resourceResult, connectionResult, typeResult] = await Promise.all([
      api.resources.list(projectId.value),
      api.connections.list(projectId.value),
      api.resourceTypes.list(),
    ])
    resources.value = resourceResult.data
    connections.value = connectionResult.data
    resourceTypes.value = typeResult
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

function typeLabel(name: string): string {
  return name.charAt(0).toUpperCase() + name.slice(1) + 's'
}

function typeIcon(name: string): string {
  switch (name) {
    case 'machine':
      return 'pi pi-desktop'
    case 'network':
      return 'pi pi-sitemap'
    default:
      return 'pi pi-box'
  }
}

onMounted(fetchData)
watch(projectId, fetchData)
watch(filterType, fetchData)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-color">
          <template v-if="filterType">{{ typeLabel(filterType) }}</template>
          <template v-else>Resources</template>
        </h1>
        <p v-if="projectStore.currentProject" class="text-muted-color mt-1">
          <template v-if="filterType">{{ typeLabel(filterType) }} in {{ projectStore.currentProject.name }}</template>
          <template v-else>Resources in {{ projectStore.currentProject.name }}</template>
        </p>
      </div>
      <Button
        v-if="projectId"
        :label="filterType ? `Create ${filterType.charAt(0).toUpperCase() + filterType.slice(1)}` : 'Create Resource'"
        icon="pi pi-plus"
        @click="filterType
          ? router.push({ name: 'resource-create-typed', params: { projectId, resourceType: filterType } })
          : router.push({ name: 'resource-create', params: { projectId } })"
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

    <!-- Data table -->
    <DataTable
      v-if="projectId && !loading"
      :value="filteredResources"
      :rowHover="true"
      stripedRows
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i :class="filterType ? typeIcon(filterType) : 'pi pi-box'" class="text-4xl mb-3 block"></i>
          <p>No {{ filterType || '' }} resources yet. Create one to get started.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.name }}</span>
        </template>
      </Column>
      <!-- Type column (only when showing all types) -->
      <Column v-if="!filterType" field="resourceType" header="Type" sortable>
        <template #body="{ data }">
          <span class="text-color">{{ data.resourceType }}</span>
        </template>
      </Column>
      <!-- Type-specific input columns -->
      <Column
        v-if="currentTypeInfo"
        v-for="field in currentTypeInfo.inputSchema"
        :key="field.name"
        :header="field.name"
        sortable
      >
        <template #body="{ data }">
          <span class="text-color">{{ fieldValue(data, field.name) }}</span>
        </template>
      </Column>
      <!-- Type-specific output columns -->
      <Column
        v-if="currentTypeInfo"
        v-for="field in currentTypeInfo.outputSchema"
        :key="'out-' + field.name"
        :header="field.name"
        sortable
      >
        <template #body="{ data }">
          <span class="text-muted-color">{{ fieldValue(data, field.name) }}</span>
        </template>
      </Column>
      <!-- Connection (only for types that require it) -->
      <Column v-if="!currentTypeInfo || currentTypeInfo.requiresConnection" header="Connection" sortable>
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
