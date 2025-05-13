import { IsBoolean, IsOptional, IsString, ValidateNested } from "class-validator";
import { DTO } from "../dto/dto";
import { Type } from "class-transformer";

export class ServiceResponse<T extends DTO> {

    @IsBoolean()
    success: boolean;

    @IsOptional()
    @ValidateNested()
    @Type(() => DTO)
    data?: T;

    @IsOptional()
    @IsString()
    error?: string;
}