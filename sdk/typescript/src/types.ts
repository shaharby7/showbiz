export interface Organization {
  id: string;
  name: string;
  displayName: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  email: string;
  organizationId: string;
  displayName: string;
  emailVerified: boolean;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Project {
  id: string;
  name: string;
  organizationId: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Connection {
  id: string;
  name: string;
  projectId: string;
  provider: string;
  credentials?: Record<string, unknown>;
  config?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface Resource {
  id: string;
  name: string;
  projectId: string;
  connectionId: string | null;
  resourceType: string;
  values: Record<string, unknown>;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface Policy {
  id: string;
  name: string;
  scope: string;
  organizationId?: string;
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface PolicyAttachment {
  id: string;
  projectId: string;
  userEmail: string;
  policyId: string;
  createdAt: string;
}

export interface ProviderInfo {
  name: string;
  resourceTypes: string[];
}

export interface ResourceTypeInfo {
  name: string;
  requiresConnection: boolean;
  inputSchema: FieldSchema[];
  outputSchema: FieldSchema[];
}

export interface FieldSchema {
  name: string;
  type: string;
  required: boolean;
  description: string;
}

export interface Pagination {
  nextCursor?: string;
  hasMore: boolean;
}

export interface ListResult<T> {
  data: T[];
  nextCursor?: string;
  hasMore: boolean;
}

export interface ListOptions {
  cursor?: string;
  limit?: number;
}

export interface RegisterInput {
  email: string;
  password: string;
  displayName?: string;
  organizationId?: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export interface CreateOrganizationInput {
  name: string;
  displayName?: string;
}

export interface UpdateOrganizationInput {
  displayName: string;
}

export interface CreateProjectInput {
  name: string;
  description?: string;
}

export interface UpdateProjectInput {
  description: string;
}

export interface CreateConnectionInput {
  name: string;
  provider: string;
  credentials?: Record<string, unknown>;
  config?: Record<string, unknown>;
}

export interface UpdateConnectionInput {
  config: Record<string, unknown>;
}

export interface CreateResourceInput {
  name: string;
  connectionId?: string;
  resourceType: string;
  values?: Record<string, unknown>;
}

export interface UpdateResourceInput {
  values: Record<string, unknown>;
}

export interface CreatePolicyInput {
  name: string;
  permissions: string[];
}

export interface AttachPolicyInput {
  userEmail: string;
  policyId: string;
}

export interface DetachPolicyInput {
  userEmail: string;
  policyId: string;
}
