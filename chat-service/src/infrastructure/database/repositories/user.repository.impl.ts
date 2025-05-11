import { Inject, Injectable } from "@nestjs/common";
import { Model } from "mongoose";
import { User, USER_MODEL } from "src/domain/models/user.model";
import { IUserRepository } from "src/domain/repository/user.repository";



@Injectable()
export class UserRepositoryImpl implements IUserRepository {
    constructor(
        @Inject(USER_MODEL) private readonly collection: Model<User>,
    ) {

    }

}