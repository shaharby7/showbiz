<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import type { Policy, PolicyAttachment } from '@showbiz/sdk'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'

const route = useRoute()
const api = useApi()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()

const orgId = computed(() => orgStore.currentOrg?.id)
const projectId = computed(() => (route.params.projectId as string) || projectStore.currentProject?.id)

// --- Policies state ---
const globalPolicies = ref<Policy[]>([])
const orgPolicies = ref<Policy[]>([])
const loadingPolicies = ref(false)
const policiesError = ref<string | null>(null)

const showCreatePolicyDialog = ref(false)
const policyForm = ref({ name: '', permissions: '' })
const creatingPolicy = ref(false)
const createPolicyError = ref<string | null>(null)

const showDeletePolicyDialog = ref(false)
const policyToDelete = ref<Policy | null>(null)
const deletingPolicy = ref(false)
const deletePolicyError = ref<string | null>(null)

// --- Attachments state ---
const attachments = ref<PolicyAttachment[]>([])
const loadingAttachments = ref(false)
const attachmentsError = ref<string | null>(null)

const showAttachDialog = ref(false)
const attachForm = ref({ userEmail: '', policyId: '' })
const attaching = ref(false)
const attachError = ref<string | null>(null)

const showDetachDialog = ref(false)
const attachmentToDetach = ref<PolicyAttachment | null>(null)
const detaching = ref(false)
const detachError = ref<string | null>(null)

// --- Computed ---
const allPolicies = computed(() => [...globalPolicies.value, ...orgPolicies.value])

function policyName(policyId: string): string {
  const policy = allPolicies.value.find((p) => p.id === policyId)
  return policy?.name || policyId
}

function policyScope(policy: Policy): 'success' | 'info' {
  return policy.scope === 'global' ? 'info' : 'success'
}

// --- Fetch ---
async function fetchPolicies() {
  loadingPolicies.value = true
  policiesError.value = null
  try {
    const promises: Promise<any>[] = [api.iam.listGlobalPolicies()]
    if (orgId.value) {
      promises.push(api.iam.listOrgPolicies(orgId.value))
    }
    const results = await Promise.all(promises)
    globalPolicies.value = results[0]
    orgPolicies.value = results[1] || []
  } catch (e: any) {
    policiesError.value = e.message || 'Failed to load policies'
  } finally {
    loadingPolicies.value = false
  }
}

async function fetchAttachments() {
  if (!orgId.value || !projectId.value) return
  loadingAttachments.value = true
  attachmentsError.value = null
  try {
    attachments.value = await api.iam.listAttachments(orgId.value, projectId.value)
  } catch (e: any) {
    attachmentsError.value = e.message || 'Failed to load attachments'
  } finally {
    loadingAttachments.value = false
  }
}

async function fetchAll() {
  await Promise.all([fetchPolicies(), fetchAttachments()])
}

// --- Policy CRUD ---
function openCreatePolicyDialog() {
  policyForm.value = { name: '', permissions: '' }
  createPolicyError.value = null
  showCreatePolicyDialog.value = true
}

async function createOrgPolicy() {
  if (!orgId.value) return
  creatingPolicy.value = true
  createPolicyError.value = null
  try {
    const permissions = policyForm.value.permissions
      .split(',')
      .map((p) => p.trim())
      .filter((p) => p.length > 0)
    await api.iam.createOrgPolicy(orgId.value, {
      name: policyForm.value.name.trim(),
      permissions,
    })
    showCreatePolicyDialog.value = false
    await fetchPolicies()
  } catch (e: any) {
    createPolicyError.value = e.message || 'Failed to create policy'
  } finally {
    creatingPolicy.value = false
  }
}

function confirmDeletePolicy(policy: Policy) {
  policyToDelete.value = policy
  deletePolicyError.value = null
  showDeletePolicyDialog.value = true
}

async function deleteOrgPolicy() {
  if (!orgId.value || !policyToDelete.value) return
  deletingPolicy.value = true
  deletePolicyError.value = null
  try {
    await api.iam.deleteOrgPolicy(orgId.value, policyToDelete.value.id)
    showDeletePolicyDialog.value = false
    policyToDelete.value = null
    await fetchPolicies()
  } catch (e: any) {
    deletePolicyError.value = e.message || 'Failed to delete policy'
  } finally {
    deletingPolicy.value = false
  }
}

// --- Attachment CRUD ---
function openAttachDialog() {
  attachForm.value = { userEmail: '', policyId: '' }
  attachError.value = null
  showAttachDialog.value = true
}

async function attachPolicy() {
  if (!orgId.value || !projectId.value) return
  attaching.value = true
  attachError.value = null
  try {
    await api.iam.attachPolicy(orgId.value, projectId.value, {
      userEmail: attachForm.value.userEmail.trim(),
      policyId: attachForm.value.policyId,
    })
    showAttachDialog.value = false
    await fetchAttachments()
  } catch (e: any) {
    attachError.value = e.message || 'Failed to attach policy'
  } finally {
    attaching.value = false
  }
}

function confirmDetach(attachment: PolicyAttachment) {
  attachmentToDetach.value = attachment
  detachError.value = null
  showDetachDialog.value = true
}

async function detachPolicy() {
  if (!orgId.value || !projectId.value || !attachmentToDetach.value) return
  detaching.value = true
  detachError.value = null
  try {
    await api.iam.detachPolicy(orgId.value, projectId.value, {
      userEmail: attachmentToDetach.value.userEmail,
      policyId: attachmentToDetach.value.policyId,
    })
    showDetachDialog.value = false
    attachmentToDetach.value = null
    await fetchAttachments()
  } catch (e: any) {
    detachError.value = e.message || 'Failed to detach policy'
  } finally {
    detaching.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

onMounted(fetchAll)
</script>

<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold text-color">IAM</h1>
      <p class="text-muted-color mt-1">Manage policies and access control</p>
    </div>

    <Message v-if="!orgId" severity="warn" :closable="false">
      Please select an organization to manage IAM.
    </Message>

    <!-- ============ SECTION 1: Policies ============ -->
    <template v-if="orgId">
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-color">Policies</h2>
          <Button
            label="Create Org Policy"
            icon="pi pi-plus"
            size="small"
            @click="openCreatePolicyDialog"
          />
        </div>

        <Message v-if="policiesError" severity="error" :closable="false">{{ policiesError }}</Message>

        <div v-if="loadingPolicies" class="flex items-center justify-center py-8">
          <i class="pi pi-spinner pi-spin text-3xl text-primary"></i>
        </div>

        <!-- Global Policies -->
        <div v-if="!loadingPolicies">
          <h3 class="text-sm font-semibold text-muted-color uppercase tracking-wide mb-2">Global Policies</h3>
          <DataTable :value="globalPolicies" :rowHover="true" stripedRows>
            <template #empty>
              <div class="text-center py-4 text-muted-color">No global policies.</div>
            </template>
            <Column field="name" header="Name" sortable>
              <template #body="{ data }">
                <span class="font-semibold text-color">{{ data.name }}</span>
              </template>
            </Column>
            <Column field="scope" header="Scope">
              <template #body="{ data }">
                <Tag :value="data.scope" :severity="policyScope(data)" />
              </template>
            </Column>
            <Column field="permissions" header="Permissions">
              <template #body="{ data }">
                <div class="flex flex-wrap gap-1">
                  <Tag
                    v-for="perm in data.permissions"
                    :key="perm"
                    :value="perm"
                    severity="secondary"
                  />
                </div>
              </template>
            </Column>
            <Column field="createdAt" header="Created" sortable>
              <template #body="{ data }">
                <span class="text-muted-color">{{ formatDate(data.createdAt) }}</span>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- Org Policies -->
        <div v-if="!loadingPolicies">
          <h3 class="text-sm font-semibold text-muted-color uppercase tracking-wide mb-2">Organization Policies</h3>
          <DataTable :value="orgPolicies" :rowHover="true" stripedRows>
            <template #empty>
              <div class="text-center py-4 text-muted-color">No organization policies yet.</div>
            </template>
            <Column field="name" header="Name" sortable>
              <template #body="{ data }">
                <span class="font-semibold text-color">{{ data.name }}</span>
              </template>
            </Column>
            <Column field="scope" header="Scope">
              <template #body="{ data }">
                <Tag :value="data.scope" :severity="policyScope(data)" />
              </template>
            </Column>
            <Column field="permissions" header="Permissions">
              <template #body="{ data }">
                <div class="flex flex-wrap gap-1">
                  <Tag
                    v-for="perm in data.permissions"
                    :key="perm"
                    :value="perm"
                    severity="secondary"
                  />
                </div>
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
                  @click="confirmDeletePolicy(data)"
                />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- ============ SECTION 2: Policy Attachments ============ -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-color">Policy Attachments</h2>
          <Button
            v-if="projectId"
            label="Attach Policy"
            icon="pi pi-plus"
            size="small"
            @click="openAttachDialog"
          />
        </div>

        <Message v-if="!projectId" severity="warn" :closable="false">
          Please select a project to manage policy attachments.
        </Message>

        <Message v-if="attachmentsError" severity="error" :closable="false">{{ attachmentsError }}</Message>

        <div v-if="loadingAttachments && projectId" class="flex items-center justify-center py-8">
          <i class="pi pi-spinner pi-spin text-3xl text-primary"></i>
        </div>

        <DataTable
          v-if="projectId && !loadingAttachments"
          :value="attachments"
          :rowHover="true"
          stripedRows
        >
          <template #empty>
            <div class="text-center py-4 text-muted-color">
              <i class="pi pi-shield text-3xl mb-2 block"></i>
              <p>No policy attachments yet.</p>
            </div>
          </template>
          <Column field="userEmail" header="User Email" sortable>
            <template #body="{ data }">
              <span class="font-semibold text-color">{{ data.userEmail }}</span>
            </template>
          </Column>
          <Column field="policyId" header="Policy" sortable>
            <template #body="{ data }">
              <span class="text-color">{{ policyName(data.policyId) }}</span>
            </template>
          </Column>
          <Column field="createdAt" header="Attached" sortable>
            <template #body="{ data }">
              <span class="text-muted-color">{{ formatDate(data.createdAt) }}</span>
            </template>
          </Column>
          <Column header="Actions" :style="{ width: '6rem' }">
            <template #body="{ data }">
              <Button
                icon="pi pi-times"
                severity="danger"
                text
                rounded
                size="small"
                @click="confirmDetach(data)"
              />
            </template>
          </Column>
        </DataTable>
      </div>
    </template>

    <!-- ============ DIALOGS ============ -->

    <!-- Create Org Policy Dialog -->
    <Dialog
      v-model:visible="showCreatePolicyDialog"
      header="Create Organization Policy"
      :modal="true"
      :style="{ width: '30rem' }"
    >
      <div class="space-y-4">
        <Message v-if="createPolicyError" severity="error" :closable="false">{{ createPolicyError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="policy-name" class="font-semibold text-color">Name</label>
          <InputText
            id="policy-name"
            v-model="policyForm.name"
            placeholder="my-policy"
            :disabled="creatingPolicy"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="policy-permissions" class="font-semibold text-color">Permissions</label>
          <InputText
            id="policy-permissions"
            v-model="policyForm.permissions"
            placeholder="resource:create, resource:read, resource:delete"
            :disabled="creatingPolicy"
          />
          <small class="text-muted-color">Comma-separated list of permissions</small>
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showCreatePolicyDialog = false"
          :disabled="creatingPolicy"
        />
        <Button
          label="Create"
          icon="pi pi-check"
          @click="createOrgPolicy"
          :loading="creatingPolicy"
          :disabled="!policyForm.name.trim() || !policyForm.permissions.trim()"
        />
      </template>
    </Dialog>

    <!-- Delete Policy Dialog -->
    <Dialog
      v-model:visible="showDeletePolicyDialog"
      header="Delete Policy"
      :modal="true"
      :style="{ width: '26rem' }"
    >
      <div class="space-y-3">
        <Message v-if="deletePolicyError" severity="error" :closable="false">{{ deletePolicyError }}</Message>
        <p class="text-color">
          Are you sure you want to delete <strong>{{ policyToDelete?.name }}</strong>?
        </p>
        <p class="text-sm text-muted-color">This will remove the policy and all its attachments.</p>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showDeletePolicyDialog = false"
          :disabled="deletingPolicy"
        />
        <Button
          label="Delete"
          icon="pi pi-trash"
          severity="danger"
          @click="deleteOrgPolicy"
          :loading="deletingPolicy"
        />
      </template>
    </Dialog>

    <!-- Attach Policy Dialog -->
    <Dialog
      v-model:visible="showAttachDialog"
      header="Attach Policy"
      :modal="true"
      :style="{ width: '30rem' }"
    >
      <div class="space-y-4">
        <Message v-if="attachError" severity="error" :closable="false">{{ attachError }}</Message>
        <div class="flex flex-col gap-2">
          <label for="attach-email" class="font-semibold text-color">User Email</label>
          <InputText
            id="attach-email"
            v-model="attachForm.userEmail"
            placeholder="user@example.com"
            :disabled="attaching"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="attach-policy" class="font-semibold text-color">Policy</label>
          <Select
            id="attach-policy"
            v-model="attachForm.policyId"
            :options="allPolicies"
            optionLabel="name"
            optionValue="id"
            placeholder="Select a policy"
            :disabled="attaching"
            class="w-full"
          />
        </div>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showAttachDialog = false"
          :disabled="attaching"
        />
        <Button
          label="Attach"
          icon="pi pi-check"
          @click="attachPolicy"
          :loading="attaching"
          :disabled="!attachForm.userEmail.trim() || !attachForm.policyId"
        />
      </template>
    </Dialog>

    <!-- Detach Policy Dialog -->
    <Dialog
      v-model:visible="showDetachDialog"
      header="Detach Policy"
      :modal="true"
      :style="{ width: '26rem' }"
    >
      <div class="space-y-3">
        <Message v-if="detachError" severity="error" :closable="false">{{ detachError }}</Message>
        <p class="text-color">
          Detach policy <strong>{{ policyName(attachmentToDetach?.policyId || '') }}</strong>
          from <strong>{{ attachmentToDetach?.userEmail }}</strong>?
        </p>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          text
          @click="showDetachDialog = false"
          :disabled="detaching"
        />
        <Button
          label="Detach"
          icon="pi pi-times"
          severity="danger"
          @click="detachPolicy"
          :loading="detaching"
        />
      </template>
    </Dialog>
  </div>
</template>
