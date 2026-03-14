<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import type { Organization, Policy } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'
import { useOrganizationStore } from '@/stores/organization'

const route = useRoute()
const router = useRouter()
const api = useApi()
const orgStore = useOrganizationStore()

const orgId = route.params.orgId as string

const org = ref<Organization | null>(null)
const loading = ref(false)
const error = ref('')

// Edit state
const editDisplayName = ref('')
const editing = ref(false)
const editLoading = ref(false)
const editError = ref('')

// Status toggle
const statusLoading = ref(false)

// Policies
const policies = ref<Policy[]>([])
const policiesLoading = ref(false)
const showPolicyDialog = ref(false)
const policyForm = ref({ name: '', permissions: '' })
const policyCreateLoading = ref(false)
const policyCreateError = ref('')

async function fetchOrg() {
  loading.value = true
  error.value = ''
  try {
    org.value = await api.organizations.get(orgId)
    editDisplayName.value = org.value.displayName
  } catch (e: any) {
    error.value = e?.message || 'Failed to load organization.'
  } finally {
    loading.value = false
  }
}

async function fetchPolicies() {
  policiesLoading.value = true
  try {
    policies.value = await api.iam.listOrgPolicies(orgId)
  } catch {
    // Policies might not be available
    policies.value = []
  } finally {
    policiesLoading.value = false
  }
}

async function saveDisplayName() {
  if (!editDisplayName.value.trim()) {
    editError.value = 'Display name is required.'
    return
  }
  editLoading.value = true
  editError.value = ''
  try {
    org.value = await api.organizations.update(orgId, {
      displayName: editDisplayName.value.trim(),
    })
    editing.value = false
    await orgStore.fetchOrganizations()
  } catch (e: any) {
    editError.value = e?.message || 'Failed to update organization.'
  } finally {
    editLoading.value = false
  }
}

async function toggleStatus() {
  if (!org.value) return
  statusLoading.value = true
  try {
    if (org.value.active) {
      await api.organizations.deactivate(orgId)
    } else {
      await api.organizations.activate(orgId)
    }
    await fetchOrg()
    await orgStore.fetchOrganizations()
  } catch (e: any) {
    error.value = e?.message || 'Failed to update status.'
  } finally {
    statusLoading.value = false
  }
}

function openPolicyDialog() {
  policyForm.value = { name: '', permissions: '' }
  policyCreateError.value = ''
  showPolicyDialog.value = true
}

async function createPolicy() {
  if (!policyForm.value.name.trim()) {
    policyCreateError.value = 'Policy name is required.'
    return
  }
  if (!policyForm.value.permissions.trim()) {
    policyCreateError.value = 'At least one permission is required.'
    return
  }
  policyCreateLoading.value = true
  policyCreateError.value = ''
  try {
    const permissions = policyForm.value.permissions
      .split(',')
      .map((p) => p.trim())
      .filter(Boolean)
    await api.iam.createOrgPolicy(orgId, {
      name: policyForm.value.name.trim(),
      permissions,
    })
    showPolicyDialog.value = false
    await fetchPolicies()
  } catch (e: any) {
    policyCreateError.value = e?.message || 'Failed to create policy.'
  } finally {
    policyCreateLoading.value = false
  }
}

async function deletePolicy(policyId: string) {
  if (!confirm('Are you sure you want to delete this policy?')) return
  try {
    await api.iam.deleteOrgPolicy(orgId, policyId)
    await fetchPolicies()
  } catch (e: any) {
    error.value = e?.message || 'Failed to delete policy.'
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

onMounted(() => {
  fetchOrg()
  fetchPolicies()
})
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
        @click="router.push({ name: 'organizations' })"
      />
      <h1 class="text-3xl font-bold text-color">Organization Detail</h1>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spinner pi-spin text-4xl text-primary"></i>
    </div>

    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>

    <template v-if="org && !loading">
      <!-- Organization Info Card -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-building text-primary"></i>
            <span>{{ org.displayName || org.name }}</span>
            <Tag
              :value="org.active ? 'Active' : 'Inactive'"
              :severity="org.active ? 'success' : 'danger'"
              class="ml-2"
            />
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-muted-color">Name</p>
                <p class="font-semibold text-color">{{ org.name }}</p>
              </div>
              <div>
                <p class="text-sm text-muted-color">Display Name</p>
                <div v-if="editing" class="flex items-center gap-2">
                  <InputText
                    v-model="editDisplayName"
                    :disabled="editLoading"
                    class="flex-1"
                  />
                  <Button
                    icon="pi pi-check"
                    severity="success"
                    size="small"
                    @click="saveDisplayName"
                    :loading="editLoading"
                  />
                  <Button
                    icon="pi pi-times"
                    severity="secondary"
                    text
                    size="small"
                    @click="editing = false; editDisplayName = org!.displayName; editError = ''"
                    :disabled="editLoading"
                  />
                </div>
                <div v-else class="flex items-center gap-2">
                  <p class="font-semibold text-color">{{ org.displayName }}</p>
                  <Button
                    icon="pi pi-pencil"
                    severity="secondary"
                    text
                    size="small"
                    @click="editing = true"
                  />
                </div>
                <Message v-if="editError" severity="error" :closable="false" class="mt-1">{{ editError }}</Message>
              </div>
              <div>
                <p class="text-sm text-muted-color">Created</p>
                <p class="font-semibold text-color">{{ formatDate(org.createdAt) }}</p>
              </div>
              <div>
                <p class="text-sm text-muted-color">Updated</p>
                <p class="font-semibold text-color">{{ formatDate(org.updatedAt) }}</p>
              </div>
            </div>

            <div class="flex flex-wrap gap-3 pt-4 border-t border-surface">
              <Button
                :label="org.active ? 'Deactivate' : 'Activate'"
                :icon="org.active ? 'pi pi-ban' : 'pi pi-check-circle'"
                :severity="org.active ? 'danger' : 'success'"
                outlined
                :loading="statusLoading"
                @click="toggleStatus"
              />
              <Button
                label="Members"
                icon="pi pi-users"
                severity="secondary"
                outlined
                @click="router.push({ name: 'organization-members', params: { orgId: org!.id } })"
              />
              <Button
                label="Projects"
                icon="pi pi-folder"
                severity="secondary"
                outlined
                @click="router.push({ name: 'projects' })"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- Policies Section -->
      <Card>
        <template #title>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <i class="pi pi-shield text-primary"></i>
              <span>Policies</span>
            </div>
            <Button
              label="Create Policy"
              icon="pi pi-plus"
              size="small"
              @click="openPolicyDialog"
            />
          </div>
        </template>
        <template #content>
          <div v-if="policiesLoading" class="flex items-center justify-center py-6">
            <i class="pi pi-spinner pi-spin text-2xl text-primary"></i>
          </div>
          <DataTable v-else :value="policies" stripedRows>
            <template #empty>
              <div class="text-center py-6 text-muted-color">
                <p>No policies defined for this organization.</p>
              </div>
            </template>
            <Column field="name" header="Name" sortable>
              <template #body="{ data }">
                <span class="font-semibold text-color">{{ data.name }}</span>
              </template>
            </Column>
            <Column field="permissions" header="Permissions">
              <template #body="{ data }">
                <div class="flex flex-wrap gap-1">
                  <Tag
                    v-for="perm in data.permissions"
                    :key="perm"
                    :value="perm"
                    severity="info"
                  />
                </div>
              </template>
            </Column>
            <Column field="createdAt" header="Created" sortable>
              <template #body="{ data }">
                {{ formatDate(data.createdAt) }}
              </template>
            </Column>
            <Column header="" style="width: 4rem">
              <template #body="{ data }">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  size="small"
                  @click="deletePolicy(data.id)"
                />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </template>

    <!-- Create Policy Dialog -->
    <Dialog
      v-model:visible="showPolicyDialog"
      header="Create Policy"
      :modal="true"
      :style="{ width: '28rem' }"
    >
      <div class="space-y-4">
        <Message v-if="policyCreateError" severity="error" :closable="false">{{ policyCreateError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="policy-name" class="font-semibold text-color">Name *</label>
          <InputText
            id="policy-name"
            v-model="policyForm.name"
            placeholder="read-only-policy"
            :disabled="policyCreateLoading"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="policy-permissions" class="font-semibold text-color">Permissions *</label>
          <InputText
            id="policy-permissions"
            v-model="policyForm.permissions"
            placeholder="read, write, admin"
            :disabled="policyCreateLoading"
          />
          <small class="text-muted-color">Comma-separated list of permissions</small>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showPolicyDialog = false"
          :disabled="policyCreateLoading"
        />
        <Button
          label="Create"
          icon="pi pi-check"
          @click="createPolicy"
          :loading="policyCreateLoading"
        />
      </template>
    </Dialog>
  </div>
</template>
