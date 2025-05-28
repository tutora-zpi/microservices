export class FailedToValidate extends Error {
    constructor(message: string) {
        super("Could not validate or get from query or command: " + message);
        this.name = "FailedToValidate";
    }
}