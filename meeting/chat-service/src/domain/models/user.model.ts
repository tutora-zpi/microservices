import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { Document } from "mongoose";
import { Model } from "./model";

export type UserDocument = User & Document;

@Schema({ timestamps: true })
export class User extends Model {
    @Prop({ required: true, type: String })
    _id: string;

    @Prop({ required: false })
    avatarURL?: string;

    @Prop({ required: true })
    firstName: string;

    @Prop({ required: true })
    lastName: string;
}

export const UserSchema = SchemaFactory.createForClass(User);

export const USER_MODEL = 'USER_MODEL';