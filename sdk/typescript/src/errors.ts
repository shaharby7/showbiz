export class ShowbizError extends Error {
  readonly code: string;
  readonly statusCode: number;

  constructor(message: string, code: string, statusCode: number) {
    super(message);
    this.name = "ShowbizError";
    this.code = code;
    this.statusCode = statusCode;
  }
}

export function isNotFound(err: unknown): err is ShowbizError {
  return err instanceof ShowbizError && err.statusCode === 404;
}

export function isConflict(err: unknown): err is ShowbizError {
  return err instanceof ShowbizError && err.statusCode === 409;
}

export function isUnauthorized(err: unknown): err is ShowbizError {
  return err instanceof ShowbizError && err.statusCode === 401;
}
