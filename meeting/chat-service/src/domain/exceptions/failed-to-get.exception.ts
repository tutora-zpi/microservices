class FailedToGet extends Error {
    constructor(message: string) {
        super("Failed to get from query or command: " + message);
        this.name = "FailedToGet";
    }
}