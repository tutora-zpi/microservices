import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { Document, Schema as MongooseSchema } from "mongoose";
import { User } from "./user.model";
import { Model } from "./model";
import { Reaction } from "./reaction.model";

export type MessageDocument = Message & Document;

@Schema({ timestamps: { createdAt: 'sentAt' } })
export class Message extends Model {
    @Prop({ required: true })
    content: string;

    @Prop({ type: String, ref: 'User', required: true })
    sender: User | string;

    @Prop({ type: String, ref: 'User', required: true })
    receiver: User | string;

    @Prop({ type: MongooseSchema.Types.ObjectId, ref: 'Chat' })
    chatID: MongooseSchema.Types.ObjectId;

    @Prop({ default: false })
    isRead: boolean;

    // message has a lot of reactions 
    @Prop({ type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Reaction' }], default: [] })
    reactions: Reaction[] | MongooseSchema.Types.ObjectId[];

    @Prop({ type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Message' }], default: [] })
    answers: Message[] | MongooseSchema.Types.ObjectId[];
}

export const MessageSchema = SchemaFactory.createForClass(Message);
export const MESSAGE_MODEL = 'MESSAGE_MODEL';