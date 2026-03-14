<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import type { Organization } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'
import { useOrganizationStore } from '@/stores/organization'

const router = useRouter()
const api = useApi()
const orgStore = useOrganizationStore()

const organizations = ref<Organization[]>([])
const loading = ref(false)
const error = ref('')

const showCreateDialog = ref(false)
const createForm = ref({ name: '', displayName: '' })
const createLoading = ref(false)
const createError = ref('')

async function fetchOrganizations() {
  loading.value = true
  error.value = ''
  try {
    const result = await api.organizations.list()
    organizations.value = result.data
  } catch (e: any) {
    error.value = e?.message || 'Failed to load organizations.'
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  createForm.value = { name: '', displayName: '' }
  createError.value = ''
  showCreateDialog.value = true
}

async function createOrganization() {
  if (!createForm.value.name.trim()) {
    createError.value = 'Name is required.'
    return
  }
  createLoading.value = true
  createError.value = ''
  try {
    await api.organizations.create({
      name: createForm.value.name.trim(),
      displayName: createForm.value.displayName.trim() || undefined,
    })
    showCreateDialog.value = false
    await fetchOrganizations()
    await orgStore.fetchOrganizations()
  } catch (e: any) {
    createError.value = e?.message || 'Failed to create organization.'
  } finally {
    createLoading.value = false
  }
}

function onRowClick(event: { data: Organization }) {
  router.push({ name: 'organization-detail', params: { orgId: event.data.id } })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

onMounted(fetchOrganizations)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-color">Organizations</h1>
      <Button
        label="Create Organization"
        icon="pi pi-plus"
        @click="openCreateDialog"
      />
    </div>

    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <DataTable
      v-else
      :value="organizations"
      :rowHover="true"
      class="cursor-pointer"
      @row-click="onRowClick"
      stripedRows
      :loading="loading"
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-building text-4xl mb-3 block"></i>
          <p>No organizations found. Create one to get started.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="displayName" header="Display Name" sortable />
      <Column field="active" header="Status" sortable>
        <template #body="{ data }">
          <Tag
            :value="data.active ? 'Active' : 'Inactive'"
            :severity="data.active ? 'success' : 'danger'"
          />
        </template>
      </Column>
      <Column field="createdAt" header="Created" sortable>
        <template #body="{ data }">
          {{ formatDate(data.createdAt) }}
        </template>
      </Column>
    </DataTable>

    <!-- Create Organization Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      header="Create Organization"
      :modal="true"
      :style="{ width: '28rem' }"
    >
      <div class="space-y-4">
        <Message v-if="createError" severity="error" :closable="false">{{ createError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="org-name" class="font-semibold text-color">Name *</label>
          <InputText
            id="org-name"
            v-model="createForm.name"
            placeholder="my-organization"
            :disabled="createLoading"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="org-display-name" class="font-semibold text-color">Display Name</label>
          <InputText
            id="org-display-name"
            v-model="createForm.displayName"
            placeholder="My Organization"
            :disabled="createLoading"
          />
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showCreateDialog = false"
          :disabled="createLoading"
        />
        <Button
          label="Create"
          icon="pi pi-check"
          @click="createOrganization"
          :loading="createLoading"
        />
      </template>
    </Dialog>
  </div>
</template>
