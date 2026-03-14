import type { ShowbizClient } from "../client.js";
import type {
  CreateProjectInput,
  ListOptions,
  ListResult,
  Project,
  UpdateProjectInput,
} from "../types.js";

export class ProjectsResource {
  constructor(private client: ShowbizClient) {}

  async create(orgId: string, input: CreateProjectInput): Promise<Project> {
    return this.client.request<Project>(
      "POST",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects`,
      input
    );
  }

  async list(
    orgId: string,
    opts?: ListOptions
  ): Promise<ListResult<Project>> {
    const qs = this.client.buildQueryString(opts);
    return this.client.request<ListResult<Project>>(
      "GET",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects${qs}`
    );
  }

  async get(orgId: string, projectId: string): Promise<Project> {
    return this.client.request<Project>(
      "GET",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}`
    );
  }

  async update(
    orgId: string,
    projectId: string,
    input: UpdateProjectInput
  ): Promise<Project> {
    return this.client.request<Project>(
      "PUT",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}`,
      input
    );
  }

  async delete(orgId: string, projectId: string): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}`
    );
  }
}
