import type { ShowbizClient } from "../client.js";
import type { ProviderInfo } from "../types.js";

export class ProvidersResource {
  constructor(private client: ShowbizClient) {}

  async list(): Promise<ProviderInfo[]> {
    return this.client.request<ProviderInfo[]>("GET", "/v1/providers");
  }

  async get(id: string): Promise<ProviderInfo> {
    return this.client.request<ProviderInfo>(
      "GET",
      `/v1/providers/${encodeURIComponent(id)}`
    );
  }
}
