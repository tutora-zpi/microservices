import { Document } from "mongoose";
import { Model } from "../models/model";
import { DTO } from "../dto/dto";

export interface IMapper<T extends DTO, K extends Model> {
    toDoc(dto: T): Partial<K>;
    toDto(doc: K & Document): T;
}
