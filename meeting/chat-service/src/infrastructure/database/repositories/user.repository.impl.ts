import { Inject, Injectable, Logger } from "@nestjs/common";
import { Model } from "mongoose";
import { UserDTO } from "src/domain/dto/user.dto";
import { UnknownException } from "src/domain/exceptions/unknown.exception";
import { UserMapper } from "src/domain/mappers/user/user.mapper";
import { User, USER_MODEL } from "src/domain/models/user.model";
import { IUserRepository } from "src/domain/repository/user.repository";

/**
 * Saves a list of user entities to the database. If a user already exists, it updates the existing record.
 * If the user does not exist, it inserts a new record. This operation is performed in bulk for efficiency.
 *
 * @param members - An array of `User` entities to be saved.
 * @returns A promise that resolves to an array of `UserDTO` objects representing the saved users.
 * @throws {UnknownException} If an error occurs during the save operation.
 *
 * The method performs the following steps:
 * 1. Logs the input users for debugging purposes.
 * 2. Constructs bulk write operations for upserting users.
 * 3. Executes the bulk write operation using the `userModel`.
 * 4. Retrieves the saved users from the database.
 * 5. Maps the saved users to `UserDTO` objects using the `UserMapper`.
 * 6. Logs the success or failure of the operation.
 */
@Injectable()
export class UserRepositoryImpl implements IUserRepository {
    private readonly logger: Logger = new Logger(UserRepositoryImpl.name);
    private readonly mapper: UserMapper = new UserMapper();

    constructor(
        @Inject(USER_MODEL) private readonly userModel: Model<User>,
    ) {

    }

    async saveUsers(members: User[]): Promise<UserDTO[]> {
        try {
            this.logger.debug('Saving users:', members);

            const operations = members.map(user => ({
                updateOne: {
                    filter: { _id: user._id },
                    update: {
                        $set: {
                            ...user,
                            updatedAt: new Date()
                        },
                        $setOnInsert: {
                            createdAt: new Date()
                        }
                    },
                    upsert: true
                }
            }));

            const res = await this.userModel.bulkWrite(operations);
            this.logger.debug('Bulk write result:', res);

            const savedUsers = await this.userModel.find({
                _id: { $in: members.map(u => u._id) }
            });

            this.logger.log(`Successfully saved ${savedUsers.length} users`);
            const users = savedUsers.map(u => this.mapper.toDto(u));

            return users;
        } catch (error) {
            this.logger.error('Error while saving users:', error.message);
            throw new UnknownException(`Failed to save users: ${error.message}`);
        }
    }

}