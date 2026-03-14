import type { ShowbizClient } from "../client.js";
import type {
  Connection,
  CreateConnectionInput,
  ListOptions,
  ListResult,
  UpdateConnectionInput,
} from "../types.js";

export class ConnectionsResource {
  constructor(private client: ShowbizClient) {}

  async create(
    projectId: string,
    input: CreateConnectionInput
  ): Promise<Connection> {
    return this.client.request<Connection>(
      "POST",
      `/v1/projects/${encodeURIComponent(projectId)}/connections`,
      input
    );
  }

  async list(
    projectId: string,
    opts?: ListOptions
  ): Promise<ListResult<Connection>> {
    const qs = this.client.buildQueryString(opts);
    return this.client.request<ListResult<Connection>>(
      "GET",
      `/v1/projects/${encodeURIComponent(projectId)}/connections${qs}`
    );
  }

  async get(projectId: string, connectionId: string): Promise<Connection> {
    return this.client.request<Connection>(
      "GET",
      `/v1/projects/${encodeURIComponent(projectId)}/connections/${encodeURIComponent(connectionId)}`
    );
  }

  async update(
    projectId: string,
    connectionId: string,
    input: UpdateConnectionInput
  ): Promise<Connection> {
    return this.client.request<Connection>(
      "PUT",
      `/v1/projects/${encodeURIComponent(projectId)}/connections/${encodeURIComponent(connectionId)}`,
      input
    );
  }

  async delete(projectId: string, connectionId: string): Promise<void> {
    return this.client.request<void>(
      "DELETE",
      `/v1/projects/${encodeURIComponent(projectId)}/connections/${encodeURIComponent(connectionId)}`
    );
  }
}
