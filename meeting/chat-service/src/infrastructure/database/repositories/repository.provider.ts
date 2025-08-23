import { Connection } from 'mongoose';
import { CHAT_MODEL, ChatSchema } from 'src/domain/models/chat.model';
import { MESSAGE_MODEL, MessageSchema } from 'src/domain/models/message.model';
import { USER_MODEL, UserSchema } from 'src/domain/models/user.model';
import { ChatRepositoryImpl } from './chat.repository.impl';
import { MessageRepositoryImpl } from './message.repository.impl';
import { UserRepositoryImpl } from './user.repository.impl';
import { DATABASE_CONNECTION } from '../database.provider';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { CHAT_REPOSITORY } from 'src/domain/repository/chat.repository';
import { USER_REPOSITORY } from 'src/domain/repository/user.repository';
import {
  REACTION_MODEL,
  ReactionSchema,
} from 'src/domain/models/reaction.model';

export const repoProviders = [
  {
    provide: USER_MODEL,
    useFactory: (connection: Connection) =>
      connection.model('User', UserSchema),
    inject: [DATABASE_CONNECTION],
  },
  {
    provide: CHAT_MODEL,
    useFactory: (connection: Connection) =>
      connection.model('Chat', ChatSchema),
    inject: [DATABASE_CONNECTION],
  },
  {
    provide: MESSAGE_MODEL,
    useFactory: (connection: Connection) =>
      connection.model('Message', MessageSchema),
    inject: [DATABASE_CONNECTION],
  },
  {
    provide: REACTION_MODEL,
    useFactory: (connection: Connection) =>
      connection.model('Reaction', ReactionSchema),
    inject: [DATABASE_CONNECTION],
  },
  {
    provide: MESSAGE_REPOSITORY,
    useClass: MessageRepositoryImpl,
  },
  {
    provide: CHAT_REPOSITORY,
    useClass: ChatRepositoryImpl,
  },
  {
    provide: USER_REPOSITORY,
    useClass: UserRepositoryImpl,
  },
];
