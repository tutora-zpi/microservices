export class UserTyping {
    constructor(public readonly chatID: string, public readonly userID: string, public readonly isTyping: boolean) { }
}

// { chatID: string; userID: string; isTyping: boolean }