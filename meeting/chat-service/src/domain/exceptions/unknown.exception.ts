export class UnknownException extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'UnknownException';
    }
}