import type { ShowbizClient } from "../client.js";
import type { AuthResponse, LoginInput, RegisterInput, User } from "../types.js";

export class AuthResource {
  constructor(private client: ShowbizClient) {}

  async register(input: RegisterInput): Promise<AuthResponse> {
    return this.client.request<AuthResponse>("POST", "/v1/auth/register", input);
  }

  async login(input: LoginInput): Promise<AuthResponse> {
    const res = await this.client.request<AuthResponse>(
      "POST",
      "/v1/auth/login",
      input
    );
    this.client.setTokens(res.accessToken, res.refreshToken);
    return res;
  }

  async refresh(): Promise<AuthResponse> {
    return this.client.request<AuthResponse>("POST", "/v1/auth/refresh");
  }

  async me(): Promise<User> {
    return this.client.request<User>("GET", "/v1/auth/me");
  }
}
