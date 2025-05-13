import { Inject, Injectable, Logger } from "@nestjs/common";
import { Model } from "mongoose";
import { UserDTO } from "src/domain/dto/user.dto";
import { UserMapper } from "src/domain/mappers/user.mapper";
import { User, USER_MODEL } from "src/domain/models/user.model";
import { IUserRepository } from "src/domain/repository/user.repository";



@Injectable()
export class UserRepositoryImpl implements IUserRepository {
    private readonly logger: Logger = new Logger(UserRepositoryImpl.name);
    private readonly mapper: UserMapper = new UserMapper();

    constructor(
        @Inject(USER_MODEL) private readonly userModel: Model<User>,
    ) {

    }

    async saveUsers(members: User[]): Promise<UserDTO[] | null> {
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
            this.logger.error('Error while saving users:', error);
            return null;
        }
    }

}