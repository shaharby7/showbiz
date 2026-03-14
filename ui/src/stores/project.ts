import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'

export const useProjectStore = defineStore('project', () => {
  const api = useApi()
  const currentProject = ref<Project | null>(null)
  const projects = ref<Project[]>([])

  async function fetchProjects(orgId: string) {
    const result = await api.projects.list(orgId)
    projects.value = result.data
  }

  function selectProject(project: Project | null) {
    currentProject.value = project
    if (project) {
      localStorage.setItem('currentProjectId', project.id)
    } else {
      localStorage.removeItem('currentProjectId')
    }
  }

  async function initialize() {
    const savedProjectId = localStorage.getItem('currentProjectId')
    if (savedProjectId) {
      const found = projects.value.find((p) => p.id === savedProjectId)
      if (found) {
        currentProject.value = found
      }
    }
  }

  return {
    currentProject,
    projects,
    fetchProjects,
    selectProject,
    initialize,
  }
})
