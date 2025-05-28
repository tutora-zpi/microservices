import { DTO } from "src/domain/dto/dto";

export interface IMeetingService {
    start<T extends DTO>(dto: T): Promise<void>;
    end<T extends DTO>(dto: T): Promise<void>;
}