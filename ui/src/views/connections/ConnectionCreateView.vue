<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ProviderInfo } from '@showbiz/sdk'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const api = useApi()

const projectId = route.params.projectId as string

const name = ref('')
const provider = ref('')
const credentialsJson = ref('')
const configJson = ref('')

const providers = ref<ProviderInfo[]>([])
const loadingProviders = ref(false)

const submitting = ref(false)
const error = ref<string | null>(null)

async function fetchProviders() {
  loadingProviders.value = true
  try {
    providers.value = await api.providers.list()
  } catch (e: any) {
    error.value = e.message || 'Failed to load providers'
  } finally {
    loadingProviders.value = false
  }
}

function parseJsonField(value: string): Record<string, unknown> | undefined {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  return JSON.parse(trimmed)
}

async function submit() {
  submitting.value = true
  error.value = null
  try {
    let credentials: Record<string, unknown> | undefined
    let config: Record<string, unknown> | undefined

    try {
      credentials = parseJsonField(credentialsJson.value)
    } catch {
      error.value = 'Credentials must be valid JSON'
      submitting.value = false
      return
    }

    try {
      config = parseJsonField(configJson.value)
    } catch {
      error.value = 'Config must be valid JSON'
      submitting.value = false
      return
    }

    await api.connections.create(projectId, {
      name: name.value.trim(),
      provider: provider.value,
      credentials,
      config,
    })

    router.push({ name: 'connections', params: { projectId } })
  } catch (e: any) {
    error.value = e.message || 'Failed to create connection'
  } finally {
    submitting.value = false
  }
}

function cancel() {
  router.push({ name: 'connections', params: { projectId } })
}

onMounted(fetchProviders)
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
        @click="cancel"
      />
      <h1 class="text-3xl font-bold text-color">Create Connection</h1>
    </div>

    <Card class="max-w-2xl">
      <template #content>
        <div class="space-y-5">
          <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

          <!-- Name -->
          <div class="flex flex-col gap-2">
            <label for="conn-name" class="font-semibold text-color">Name</label>
            <InputText
              id="conn-name"
              v-model="name"
              placeholder="my-connection"
              :disabled="submitting"
            />
          </div>

          <!-- Provider -->
          <div class="flex flex-col gap-2">
            <label for="conn-provider" class="font-semibold text-color">Provider</label>
            <Select
              id="conn-provider"
              v-model="provider"
              :options="providers"
              optionLabel="name"
              optionValue="name"
              placeholder="Select a provider"
              :loading="loadingProviders"
              :disabled="submitting"
              class="w-full"
            />
          </div>

          <!-- Credentials -->
          <div class="flex flex-col gap-2">
            <label for="conn-credentials" class="font-semibold text-color">Credentials</label>
            <p class="text-sm text-muted-color">Provider-specific credentials as JSON</p>
            <Textarea
              id="conn-credentials"
              v-model="credentialsJson"
              rows="5"
              placeholder='{ "apiKey": "..." }'
              :disabled="submitting"
              class="font-mono text-sm"
            />
          </div>

          <!-- Config -->
          <div class="flex flex-col gap-2">
            <label for="conn-config" class="font-semibold text-color">Config (optional)</label>
            <Textarea
              id="conn-config"
              v-model="configJson"
              rows="3"
              placeholder='{ "region": "us-east-1" }'
              :disabled="submitting"
              class="font-mono text-sm"
            />
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-2">
            <Button
              label="Create Connection"
              icon="pi pi-check"
              @click="submit"
              :loading="submitting"
              :disabled="!name.trim() || !provider"
            />
            <Button
              label="Cancel"
              severity="secondary"
              text
              @click="cancel"
              :disabled="submitting"
            />
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>
