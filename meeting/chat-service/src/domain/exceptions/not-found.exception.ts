class RecordNotFound extends Error {
    constructor(message: string) {
        super("Record not found: " + message);
        this.name = "RecordNotFound";
    }
}