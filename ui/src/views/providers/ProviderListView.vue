<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ProviderInfo } from '@showbiz/sdk'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'

const api = useApi()

const providers = ref<ProviderInfo[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

async function fetchProviders() {
  loading.value = true
  error.value = null
  try {
    providers.value = await api.providers.list()
  } catch (e: any) {
    error.value = e.message || 'Failed to load providers'
  } finally {
    loading.value = false
  }
}

onMounted(fetchProviders)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold text-color">Providers</h1>
      <p class="text-muted-color mt-1">Available infrastructure providers and their supported resource types.</p>
    </div>

    <!-- Error -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- DataTable -->
    <DataTable
      v-if="!loading"
      :value="providers"
      :rowHover="true"
      stripedRows
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-server text-4xl mb-3 block"></i>
          <p>No providers available.</p>
        </div>
      </template>
      <Column field="name" header="Provider" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-3">
            <div class="flex items-center justify-center w-10 h-10 rounded-lg bg-blue-100 text-blue-600">
              <i class="pi pi-server text-lg"></i>
            </div>
            <span class="font-semibold text-color">{{ data.name }}</span>
          </div>
        </template>
      </Column>
      <Column header="Resource Types">
        <template #body="{ data }">
          <div class="flex flex-wrap gap-2">
            <Tag
              v-for="rt in data.resourceTypes"
              :key="rt"
              :value="rt"
              severity="info"
            />
            <span v-if="!data.resourceTypes?.length" class="text-muted-color">None</span>
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
