import { DTO } from "../dto/dto";

export class ServiceResponse<T extends DTO> {
  readonly success: boolean;

  readonly data?: T;

  readonly error?: string;

  constructor(data?: T, error?: string) {
    this.success = !!data && !error;
    this.data = data;
    this.error = error;
  }
}
