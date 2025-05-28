import { Document } from "mongoose";
import { Model } from "../models/model";
import { DTO } from "../dto/dto";

export interface IMapper<T extends DTO, M extends Model> {
    toDoc(dto: T): Partial<M>;
    toDto(doc: M & Document): T;
}
