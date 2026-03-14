import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Organization } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'

export const useOrganizationStore = defineStore('organization', () => {
  const api = useApi()
  const currentOrg = ref<Organization | null>(null)
  const organizations = ref<Organization[]>([])

  async function fetchOrganizations() {
    const result = await api.organizations.list()
    organizations.value = result.data
  }

  function selectOrganization(org: Organization | null) {
    currentOrg.value = org
    if (org) {
      localStorage.setItem('currentOrgId', org.id)
    } else {
      localStorage.removeItem('currentOrgId')
    }
  }

  async function initialize() {
    await fetchOrganizations()
    const savedOrgId = localStorage.getItem('currentOrgId')
    if (savedOrgId) {
      const found = organizations.value.find((o) => o.id === savedOrgId)
      if (found) {
        currentOrg.value = found
      } else if (organizations.value.length > 0) {
        selectOrganization(organizations.value[0])
      }
    } else if (organizations.value.length > 0) {
      selectOrganization(organizations.value[0])
    }
  }

  return {
    currentOrg,
    organizations,
    fetchOrganizations,
    selectOrganization,
    initialize,
  }
})
