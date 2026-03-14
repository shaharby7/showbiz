import type { ShowbizClient } from "../client.js";
import type {
  CreateResourceInput,
  ListOptions,
  ListResult,
  Resource,
  UpdateResourceInput,
} from "../types.js";

export class ResourcesResource {
  constructor(private client: ShowbizClient) {}

  async create(
    projectId: string,
    input: CreateResourceInput
  ): Promise<Resource> {
    return this.client.request<Resource>(
      "POST",
      `/v1/projects/${encodeURIComponent(projectId)}/resources`,
      input
    );
  }

  async list(
    projectId: string,
    opts?: ListOptions
  ): Promise<ListResult<Resource>> {
    const qs = this.client.buildQueryString(opts);
    return this.client.request<ListResult<Resource>>(
      "GET",
      `/v1/projects/${encodeURIComponent(projectId)}/resources${qs}`
    );
  }

  async get(projectId: string, resourceId: string): Promise<Resource> {
    return this.client.request<Resource>(
      "GET",
      `/v1/projects/${encodeURIComponent(projectId)}/resources/${encodeURIComponent(resourceId)}`
    );
  }

  async update(
    projectId: string,
    resourceId: string,
    input: UpdateResourceInput
  ): Promise<Resource> {
    return this.client.request<Resource>(
      "PUT",
      `/v1/projects/${encodeURIComponent(projectId)}/resources/${encodeURIComponent(resourceId)}`,
      input
    );
  }

  async delete(projectId: string, resourceId: string): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/projects/${encodeURIComponent(projectId)}/resources/${encodeURIComponent(resourceId)}`
    );
  }
}
