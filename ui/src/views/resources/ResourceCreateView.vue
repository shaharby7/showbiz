<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Connection, ResourceTypeInfo, CreateResourceInput } from '@showbiz/sdk'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Message from 'primevue/message'
import Tag from 'primevue/tag'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const api = useApi()

const projectId = route.params.projectId as string

const connections = ref<Connection[]>([])
const resourceTypes = ref<ResourceTypeInfo[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const name = ref('')
const selectedType = ref<string>('')
const selectedConnectionId = ref<string>('')
const fieldValues = ref<Record<string, string | number | null>>({})
const creating = ref(false)
const createError = ref<string | null>(null)

const currentTypeInfo = computed(() =>
  resourceTypes.value.find((rt) => rt.name === selectedType.value) || null,
)

const requiresConnection = computed(() => currentTypeInfo.value?.requiresConnection ?? false)

const filteredConnections = computed(() => {
  if (!requiresConnection.value) return []
  return connections.value
})

// Reset fields when type changes
watch(selectedType, () => {
  fieldValues.value = {}
  selectedConnectionId.value = ''
})

async function fetchData() {
  loading.value = true
  error.value = null
  try {
    const [connResult, typeResult] = await Promise.all([
      api.connections.list(projectId),
      api.resourceTypes.list(),
    ])
    connections.value = connResult.data
    resourceTypes.value = typeResult

    // Pre-select type from route param if present
    const routeType = route.params.resourceType as string
    if (routeType && typeResult.some((rt: ResourceTypeInfo) => rt.name === routeType)) {
      selectedType.value = routeType
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load data'
  } finally {
    loading.value = false
  }
}

const routeType = computed(() => (route.params.resourceType as string) || null)

function goBack() {
  if (routeType.value) {
    router.push({ name: 'resources-by-type', params: { projectId, resourceType: routeType.value } })
  } else {
    router.push({ name: 'resources', params: { projectId } })
  }
}

async function createResource() {
  creating.value = true
  createError.value = null
  try {
    const input: CreateResourceInput = {
      name: name.value,
      resourceType: selectedType.value,
      values: { ...fieldValues.value },
    }
    if (requiresConnection.value && selectedConnectionId.value) {
      input.connectionId = selectedConnectionId.value
    }
    await api.resources.create(projectId, input)
    goBack()
  } catch (e: any) {
    createError.value = e.message || 'Failed to create resource'
  } finally {
    creating.value = false
  }
}

const canSubmit = computed(() => {
  if (!name.value.trim() || !selectedType.value || creating.value) return false
  if (requiresConnection.value && !selectedConnectionId.value) return false
  // Check required fields
  const info = currentTypeInfo.value
  if (info) {
    for (const field of info.inputSchema) {
      if (field.required) {
        const val = fieldValues.value[field.name]
        if (val === undefined || val === null || val === '') return false
      }
    }
  }
  return true
})

function typeLabel(name: string): string {
  return name.charAt(0).toUpperCase() + name.slice(1)
}

onMounted(fetchData)
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
        @click="goBack()"
      />
      <h1 class="text-3xl font-bold text-color">Create Resource</h1>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <!-- Error loading data -->
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <!-- Form -->
    <div v-if="!loading" class="max-w-xl space-y-5">
      <Message v-if="createError" severity="error" :closable="false">{{ createError }}</Message>

      <!-- Name -->
      <div class="flex flex-col gap-2">
        <label for="resource-name" class="font-semibold text-color">Name</label>
        <InputText
          id="resource-name"
          v-model="name"
          placeholder="my-resource"
          :disabled="creating"
        />
      </div>

      <!-- Resource Type -->
      <div class="flex flex-col gap-2">
        <label for="resource-type" class="font-semibold text-color">Resource Type</label>
        <Select
          id="resource-type"
          v-model="selectedType"
          :options="resourceTypes"
          optionLabel="name"
          optionValue="name"
          placeholder="Select resource type"
          :disabled="creating"
          class="w-full"
        >
          <template #option="{ option }">
            <div class="flex items-center gap-2">
              <span>{{ typeLabel(option.name) }}</span>
              <Tag
                v-if="!option.requiresConnection"
                value="Showbiz-managed"
                severity="info"
                class="text-xs"
              />
            </div>
          </template>
        </Select>
      </div>

      <!-- Connection (only for types that require it) -->
      <div v-if="requiresConnection" class="flex flex-col gap-2">
        <label for="resource-connection" class="font-semibold text-color">Connection</label>
        <Select
          id="resource-connection"
          v-model="selectedConnectionId"
          :options="filteredConnections"
          optionLabel="name"
          optionValue="id"
          placeholder="Select a connection"
          :disabled="creating"
          class="w-full"
        />
        <small v-if="filteredConnections.length === 0" class="text-muted-color">
          No connections available. Create a connection first.
        </small>
      </div>

      <div v-if="!requiresConnection && selectedType" class="flex items-center gap-2 text-sm text-muted-color">
        <i class="pi pi-info-circle"></i>
        <span>This resource type is managed by Showbiz — no provider connection needed.</span>
      </div>

      <!-- Dynamic fields based on resource type schema -->
      <template v-if="currentTypeInfo">
        <div
          v-for="field in currentTypeInfo.inputSchema"
          :key="field.name"
          class="flex flex-col gap-2"
        >
          <label :for="'field-' + field.name" class="font-semibold text-color">
            {{ field.name }}
            <span v-if="field.required" class="text-red-500">*</span>
          </label>
          <InputNumber
            v-if="field.type === 'number'"
            :id="'field-' + field.name"
            :modelValue="fieldValues[field.name] as number | null"
            @update:modelValue="(v: number | null) => fieldValues[field.name] = v"
            :placeholder="field.description"
            :disabled="creating"
            class="w-full"
          />
          <InputText
            v-else
            :id="'field-' + field.name"
            :modelValue="(fieldValues[field.name] as string) ?? ''"
            @update:modelValue="(v: string | undefined) => fieldValues[field.name] = v ?? ''"
            :placeholder="field.description"
            :disabled="creating"
          />
          <small class="text-muted-color">{{ field.description }}</small>
        </div>
      </template>

      <!-- Actions -->
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
          @click="goBack()"
          :disabled="creating"
        />
      </div>
    </div>
  </div>
</template>
