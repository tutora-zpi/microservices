import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { Document, Schema as MongooseSchema } from "mongoose";
import { User } from "./user.model";
import { Model } from "./model";

export type ReactionDocument = Document & Reaction;

@Schema({ timestamps: true })
export class Reaction extends Model {
    @Prop({ type: String, ref: 'User', field: 'sourceID', required: true })
    user: User | string;

    @Prop({ required: true })
    emoji: string;
}

export const ReactionSchema = SchemaFactory.createForClass(Reaction);
