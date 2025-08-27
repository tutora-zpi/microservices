import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import mongoose, { Document, Schema as MongooseSchema, Model as MongooseModel } from 'mongoose';
import { Message } from './message.model';
import { User } from './user.model';
import { Model } from './model';

export type ChatDocument = Chat & Document;

@Schema({ timestamps: true })
export class Chat extends Model {
  @Prop({ required: true })
  id: string;

  @Prop({ type: [{ type: String, ref: 'User' }], default: [] })
  members: User[] | string[];

  @Prop({
    type: [{ type: MongooseSchema.Types.ObjectId, ref: 'Message' }],
    default: [],
  })
  messages: Message[] | MongooseSchema.Types.ObjectId[];
}

export const ChatSchema = SchemaFactory.createForClass(Chat);

export const CHAT_MODEL = 'CHAT_MODEL';

ChatSchema.pre('findOneAndDelete', async function () {
  const filter = this.getFilter();
  const id = filter.id;

  // removing message related to chat
  const messageModel: MongooseModel<Message> = mongoose.model<Message>('Message');
  await messageModel.deleteMany({ chatID: id });
});