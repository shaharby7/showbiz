import type { ShowbizClient } from "../client.js";
import type {
  CreateOrganizationInput,
  ListOptions,
  ListResult,
  Organization,
  UpdateOrganizationInput,
  User,
} from "../types.js";

export class OrganizationsResource {
  constructor(private client: ShowbizClient) {}

  async create(input: CreateOrganizationInput): Promise<Organization> {
    return this.client.request<Organization>(
      "POST",
      "/v1/organizations",
      input
    );
  }

  async get(id: string): Promise<Organization> {
    return this.client.request<Organization>(
      "GET",
      `/v1/organizations/${encodeURIComponent(id)}`
    );
  }

  async list(opts?: ListOptions): Promise<ListResult<Organization>> {
    const qs = this.client.buildQueryString(opts);
    return this.client.request<ListResult<Organization>>(
      "GET",
      `/v1/organizations${qs}`
    );
  }

  async update(
    id: string,
    input: UpdateOrganizationInput
  ): Promise<Organization> {
    return this.client.request<Organization>(
      "PUT",
      `/v1/organizations/${encodeURIComponent(id)}`,
      input
    );
  }

  async deactivate(id: string): Promise<void> {
    return this.client.request<void>(
      "POST",
      `/v1/organizations/${encodeURIComponent(id)}/deactivate`
    );
  }

  async activate(id: string): Promise<void> {
    return this.client.request<void>(
      "POST",
      `/v1/organizations/${encodeURIComponent(id)}/activate`
    );
  }

  async listMembers(orgId: string): Promise<User[]> {
    return this.client.request<User[]>(
      "GET",
      `/v1/organizations/${encodeURIComponent(orgId)}/members`
    );
  }

  async addMember(orgId: string, email: string): Promise<void> {
    return this.client.request<void>(
      "POST",
      `/v1/organizations/${encodeURIComponent(orgId)}/members`,
      { email }
    );
  }

  async removeMember(orgId: string, email: string): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/organizations/${encodeURIComponent(orgId)}/members/${encodeURIComponent(email)}`
    );
  }
}
