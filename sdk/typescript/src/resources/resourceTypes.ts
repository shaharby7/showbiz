import type { ShowbizClient } from "../client.js";
import type { ResourceTypeInfo } from "../types.js";

export class ResourceTypesResource {
  constructor(private client: ShowbizClient) {}

  async list(): Promise<ResourceTypeInfo[]> {
    return this.client.request<ResourceTypeInfo[]>("GET", "/v1/resource-types");
  }

  async get(name: string): Promise<ResourceTypeInfo> {
    return this.client.request<ResourceTypeInfo>(
      "GET",
      `/v1/resource-types/${encodeURIComponent(name)}`
    );
  }
}
