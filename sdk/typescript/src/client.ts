import { ShowbizError } from "./errors.js";
import { AuthResource } from "./resources/auth.js";
import { OrganizationsResource } from "./resources/organizations.js";
import { ProjectsResource } from "./resources/projects.js";
import { ConnectionsResource } from "./resources/connections.js";
import { ResourcesResource } from "./resources/resources.js";
import { IAMResource } from "./resources/iam.js";
import { ProvidersResource } from "./resources/providers.js";
import { ResourceTypesResource } from "./resources/resourceTypes.js";
import type { ListOptions } from "./types.js";

export interface ClientOptions {
  baseURL: string;
  token?: string;
  refreshToken?: string;
  onTokenRefresh?: (tokens: {
    accessToken: string;
    refreshToken: string;
  }) => void;
}

export class ShowbizClient {
  private _baseURL: string;
  private _token?: string;
  private _refreshToken?: string;
  private _onTokenRefresh?: (tokens: {
    accessToken: string;
    refreshToken: string;
  }) => void;
  private _refreshing: Promise<void> | null = null;

  readonly auth: AuthResource;
  readonly organizations: OrganizationsResource;
  readonly projects: ProjectsResource;
  readonly connections: ConnectionsResource;
  readonly resources: ResourcesResource;
  readonly iam: IAMResource;
  readonly providers: ProvidersResource;
  readonly resourceTypes: ResourceTypesResource;

  constructor(options: ClientOptions) {
    this._baseURL = options.baseURL.replace(/\/+$/, "");
    this._token = options.token;
    this._refreshToken = options.refreshToken;
    this._onTokenRefresh = options.onTokenRefresh;

    this.auth = new AuthResource(this);
    this.organizations = new OrganizationsResource(this);
    this.projects = new ProjectsResource(this);
    this.connections = new ConnectionsResource(this);
    this.resources = new ResourcesResource(this);
    this.iam = new IAMResource(this);
    this.providers = new ProvidersResource(this);
    this.resourceTypes = new ResourceTypesResource(this);
  }

  setTokens(accessToken: string, refreshToken: string): void {
    this._token = accessToken;
    this._refreshToken = refreshToken;
  }

  buildQueryString(opts?: ListOptions): string {
    if (!opts) return "";
    const params = new URLSearchParams();
    if (opts.cursor) params.set("cursor", opts.cursor);
    if (opts.limit !== undefined) params.set("limit", String(opts.limit));
    const qs = params.toString();
    return qs ? `?${qs}` : "";
  }

  async request<T>(
    method: string,
    path: string,
    body?: unknown
  ): Promise<T> {
    const url = `${this._baseURL}${path}`;
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this._token) {
      headers["Authorization"] = `Bearer ${this._token}`;
    }

    const res = await fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (res.status === 401 && this._refreshToken && !path.includes("/auth/")) {
      await this._doRefresh();
      return this._retry<T>(method, url, body);
    }

    return this._handleResponse<T>(res);
  }

  private async _doRefresh(): Promise<void> {
    if (this._refreshing) {
      await this._refreshing;
      return;
    }
    this._refreshing = (async () => {
      const res = await fetch(`${this._baseURL}/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${this._refreshToken}`,
        },
      });
      if (!res.ok) {
        const err = await this._parseError(res);
        throw err;
      }
      const data = (await res.json()) as {
        accessToken: string;
        refreshToken: string;
      };
      this._token = data.accessToken;
      this._refreshToken = data.refreshToken;
      this._onTokenRefresh?.({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
      });
    })();
    try {
      await this._refreshing;
    } finally {
      this._refreshing = null;
    }
  }

  private async _retry<T>(
    method: string,
    url: string,
    body?: unknown
  ): Promise<T> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this._token) {
      headers["Authorization"] = `Bearer ${this._token}`;
    }
    const res = await fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
    return this._handleResponse<T>(res);
  }

  private async _handleResponse<T>(res: Response): Promise<T> {
    if (!res.ok) {
      throw await this._parseError(res);
    }
    if (res.status === 204) {
      return undefined as T;
    }
    return (await res.json()) as T;
  }

  private async _parseError(res: Response): Promise<ShowbizError> {
    let code = "unknown";
    let message = res.statusText;
    try {
      const body = (await res.json()) as {
        code?: string;
        message?: string;
        error?: string;
      };
      if (body.code) code = body.code;
      if (body.message) message = body.message;
      else if (body.error) message = body.error;
    } catch {
      // body not JSON
    }
    return new ShowbizError(message, code, res.status);
  }
}
