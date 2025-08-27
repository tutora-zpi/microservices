import { UserDTO } from '../dto/user.dto';
import { User } from '../models/user.model';

export const USER_REPOSITORY = 'IUserRepository';

export interface IUserRepository {
  // Create or update method
  saveUsers(members: User[]): Promise<UserDTO[]>;
}
