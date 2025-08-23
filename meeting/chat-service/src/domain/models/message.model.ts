import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document, Types } from 'mongoose';
import { User } from './user.model';
import { Model } from './model';
import { Reaction } from './reaction.model';

export type MessageDocument = Message & Document;

@Schema({ timestamps: { createdAt: 'sentAt', updatedAt: 'updatedAt' } })
export class Message extends Model {
  @Prop({ required: true })
  content: string;

  @Prop({ type: String, ref: 'User', required: true })
  sender: User | string;

  @Prop({ type: String, required: true })
  chatID: string;

  @Prop({ default: false })
  isRead: boolean;

  @Prop({ type: [{ type: String, ref: 'Reaction' }], default: [] })
  reactions: Reaction[] | string[];

  @Prop({ type: [{ type: String, ref: 'Message' }], default: [] })
  answers: Message[] | string[];

  sentAt: Date;
  updatedAt: Date;
}

export const MessageSchema = SchemaFactory.createForClass(Message);
export const MESSAGE_MODEL = 'MESSAGE_MODEL';
