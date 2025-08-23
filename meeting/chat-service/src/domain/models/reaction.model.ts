import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document } from 'mongoose';
import { User } from './user.model';
import { Model } from './model';

export type ReactionDocument = Document & Reaction;

@Schema({ timestamps: true })
export class Reaction extends Model {
  @Prop({ type: String, ref: 'User', required: true })
  user: User | string;

  @Prop({ required: true })
  emoji: string;

  @Prop({ type: String, ref: 'Message', required: true })
  messageID: string;
}

export const ReactionSchema = SchemaFactory.createForClass(Reaction);
export const REACTION_MODEL = 'REACTION_MODEL';
