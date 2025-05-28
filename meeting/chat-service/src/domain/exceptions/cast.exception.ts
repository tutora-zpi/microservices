class CastingError extends Error {
    constructor(message: string) {
        super("Could not cast something: " + message);
        this.name = "CastingError";
    }
}