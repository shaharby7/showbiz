<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Connection, CreateResourceInput } from '@showbiz/sdk'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const api = useApi()

const projectId = route.params.projectId as string

const connections = ref<Connection[]>([])
const resourceTypes = ref<string[]>(['machine', 'network'])
const loading = ref(false)
const error = ref<string | null>(null)

const form = ref<CreateResourceInput>({
  name: '',
  connectionId: '',
  resourceType: '',
  values: {},
})
const valuesJson = ref('{}')
const creating = ref(false)
const createError = ref<string | null>(null)
const jsonError = ref<string | null>(null)

const selectedConnection = computed(() =>
  connections.value.find((c) => c.id === form.value.connectionId) || null,
)

async function fetchConnections() {
  loading.value = true
  error.value = null
  try {
    const result = await api.connections.list(projectId)
    connections.value = result.data
  } catch (e: any) {
    error.value = e.message || 'Failed to load connections'
  } finally {
    loading.value = false
  }
}

watch(
  () => form.value.connectionId,
  async (connectionId) => {
    if (!connectionId) {
      resourceTypes.value = ['machine', 'network']
      return
    }
    const conn = connections.value.find((c) => c.id === connectionId)
    if (conn?.provider) {
      try {
        const providerInfo = await api.providers.get(conn.provider)
        if (providerInfo.resourceTypes?.length) {
          resourceTypes.value = providerInfo.resourceTypes
          return
        }
      } catch {
        // fall back to defaults
      }
    }
    resourceTypes.value = ['machine', 'network']
  },
)

function validateJson() {
  jsonError.value = null
  try {
    JSON.parse(valuesJson.value)
    return true
  } catch {
    jsonError.value = 'Invalid JSON'
    return false
  }
}

async function createResource() {
  if (!validateJson()) return
  creating.value = true
  createError.value = null
  try {
    const input: CreateResourceInput = {
      name: form.value.name,
      connectionId: form.value.connectionId,
      resourceType: form.value.resourceType,
      values: JSON.parse(valuesJson.value),
    }
    await api.resources.create(projectId, input)
    router.push({ name: 'resources', params: { projectId } })
  } catch (e: any) {
    createError.value = e.message || 'Failed to create resource'
  } finally {
    creating.value = false
  }
}

const canSubmit = computed(
  () =>
    form.value.name.trim() !== '' &&
    form.value.connectionId !== '' &&
    form.value.resourceType !== '' &&
    !creating.value,
)

onMounted(fetchConnections)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-3">
      <Button
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        rounded
        @click="router.push({ name: 'resources', params: { projectId } })"
      />
      <h1 class="text-3xl font-bold text-color">Create Resource</h1>
    </div>

    <!-- Loading connections -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- Error loading connections -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <!-- Form -->
    <div v-if="!loading" class="max-w-xl space-y-5">
      <Message v-if="createError" severity="error" :closable="false">{{ createError }}</Message>

      <div class="flex flex-col gap-2">
        <label for="resource-name" class="font-semibold text-color">Name</label>
        <InputText
          id="resource-name"
          v-model="form.name"
          placeholder="my-resource"
          :disabled="creating"
        />
      </div>

      <div class="flex flex-col gap-2">
        <label for="resource-connection" class="font-semibold text-color">Connection</label>
        <Select
          id="resource-connection"
          v-model="form.connectionId"
          :options="connections"
          optionLabel="name"
          optionValue="id"
          placeholder="Select a connection"
          :disabled="creating"
          class="w-full"
        />
        <small v-if="selectedConnection" class="text-muted-color">
          Provider: {{ selectedConnection.provider }}
        </small>
      </div>

      <div class="flex flex-col gap-2">
        <label for="resource-type" class="font-semibold text-color">Resource Type</label>
        <Select
          id="resource-type"
          v-model="form.resourceType"
          :options="resourceTypes"
          placeholder="Select resource type"
          :disabled="creating"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-2">
        <label for="resource-values" class="font-semibold text-color">Values (JSON)</label>
        <Textarea
          id="resource-values"
          v-model="valuesJson"
          rows="6"
          placeholder='{"key": "value"}'
          :disabled="creating"
          class="font-mono text-sm"
        />
        <small v-if="jsonError" class="text-red-500">{{ jsonError }}</small>
      </div>

      <div class="flex gap-3">
        <Button
          label="Create"
          icon="pi pi-check"
          @click="createResource"
          :loading="creating"
          :disabled="!canSubmit"
        />
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="router.push({ name: 'resources', params: { projectId } })"
          :disabled="creating"
        />
      </div>
    </div>
  </div>
</template>
