<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import type { User } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const api = useApi()

const orgId = route.params.orgId as string

const members = ref<User[]>([])
const loading = ref(false)
const error = ref('')

const showAddDialog = ref(false)
const addEmail = ref('')
const addLoading = ref(false)
const addError = ref('')

async function fetchMembers() {
  loading.value = true
  error.value = ''
  try {
    members.value = await api.organizations.listMembers(orgId)
  } catch (e: any) {
    error.value = e?.message || 'Failed to load members.'
  } finally {
    loading.value = false
  }
}

function openAddDialog() {
  addEmail.value = ''
  addError.value = ''
  showAddDialog.value = true
}

async function addMember() {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!addEmail.value.trim()) {
    addError.value = 'Email is required.'
    return
  }
  if (!emailRegex.test(addEmail.value.trim())) {
    addError.value = 'Please enter a valid email address.'
    return
  }
  addLoading.value = true
  addError.value = ''
  try {
    await api.organizations.addMember(orgId, addEmail.value.trim())
    showAddDialog.value = false
    await fetchMembers()
  } catch (e: any) {
    addError.value = e?.message || 'Failed to add member.'
  } finally {
    addLoading.value = false
  }
}

async function removeMember(email: string) {
  if (!confirm(`Remove member "${email}" from this organization?`)) return
  try {
    await api.organizations.removeMember(orgId, email)
    await fetchMembers()
  } catch (e: any) {
    error.value = e?.message || 'Failed to remove member.'
  }
}

onMounted(fetchMembers)
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
        @click="router.push({ name: 'organization-detail', params: { orgId } })"
      />
      <h1 class="text-3xl font-bold text-color">Organization Members</h1>
    </div>

    <div class="flex justify-end">
      <Button
        label="Add Member"
        icon="pi pi-user-plus"
        @click="openAddDialog"
      />
    </div>

    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <DataTable
      v-else
      :value="members"
      stripedRows
      :loading="loading"
    >
      <template #empty>
        <div class="text-center py-8 text-muted-color">
          <i class="pi pi-users text-4xl mb-3 block"></i>
          <p>No members found. Add a member to get started.</p>
        </div>
      </template>
      <Column field="email" header="Email" sortable>
        <template #body="{ data }">
          <span class="font-semibold text-color">{{ data.email }}</span>
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
      <Column header="" style="width: 4rem">
        <template #body="{ data }">
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            @click="removeMember(data.email)"
          />
        </template>
      </Column>
    </DataTable>

    <!-- Add Member Dialog -->
    <Dialog
      v-model:visible="showAddDialog"
      header="Add Member"
      :modal="true"
      :style="{ width: '28rem' }"
    >
      <div class="space-y-4">
        <Message v-if="addError" severity="error" :closable="false">{{ addError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="member-email" class="font-semibold text-color">Email *</label>
          <InputText
            id="member-email"
            v-model="addEmail"
            placeholder="user@example.com"
            :disabled="addLoading"
          />
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showAddDialog = false"
          :disabled="addLoading"
        />
        <Button
          label="Add"
          icon="pi pi-check"
          @click="addMember"
          :loading="addLoading"
        />
      </template>
    </Dialog>
  </div>
</template>
