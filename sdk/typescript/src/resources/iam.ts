import type { ShowbizClient } from "../client.js";
import type {
  AttachPolicyInput,
  CreatePolicyInput,
  DetachPolicyInput,
  Policy,
  PolicyAttachment,
} from "../types.js";

export class IAMResource {
  constructor(private client: ShowbizClient) {}

  async listGlobalPolicies(): Promise<Policy[]> {
    return this.client.request<Policy[]>("GET", "/v1/iam/policies");
  }

  async getPolicy(policyId: string): Promise<Policy> {
    return this.client.request<Policy>(
      "GET",
      `/v1/iam/policies/${encodeURIComponent(policyId)}`
    );
  }

  async listOrgPolicies(orgId: string): Promise<Policy[]> {
    return this.client.request<Policy[]>(
      "GET",
      `/v1/organizations/${encodeURIComponent(orgId)}/policies`
    );
  }

  async createOrgPolicy(
    orgId: string,
    input: CreatePolicyInput
  ): Promise<Policy> {
    return this.client.request<Policy>(
      "POST",
      `/v1/organizations/${encodeURIComponent(orgId)}/policies`,
      input
    );
  }

  async deleteOrgPolicy(orgId: string, policyId: string): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/organizations/${encodeURIComponent(orgId)}/policies/${encodeURIComponent(policyId)}`
    );
  }

  async listAttachments(
    orgId: string,
    projectId: string
  ): Promise<PolicyAttachment[]> {
    return this.client.request<PolicyAttachment[]>(
      "GET",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}/attachments`
    );
  }

  async attachPolicy(
    orgId: string,
    projectId: string,
    input: AttachPolicyInput
  ): Promise<PolicyAttachment> {
    return this.client.request<PolicyAttachment>(
      "POST",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}/attachments`,
      input
    );
  }

  async detachPolicy(
    orgId: string,
    projectId: string,
    input: DetachPolicyInput
  ): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/organizations/${encodeURIComponent(orgId)}/projects/${encodeURIComponent(projectId)}/attachments`,
      input
    );
  }
}
